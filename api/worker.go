package api

import (
	"fmt"
	"log"
	"time"

	"sistema-ecommerce-go/modelos"
	"sistema-ecommerce-go/servicios"
)

// ============================================================================
// CONCURRENCIA: WORKER POOL CON GOROUTINES Y CHANNELS (Unidad 4)
// ============================================================================
//
// Este archivo implementa el patrón Worker Pool para procesar pedidos de forma
// ASÍNCRONA sin bloquear las respuestas HTTP del servidor.
//
// FLUJO DE CONCURRENCIA:
//
//   Endpoint HTTP (CheckoutHandler)
//        |
//        v
//   [Canal de Solicitudes] ──► Goroutine Worker 1 ──► [Canal de Resultados]
//                          ──► Goroutine Worker 2 ──►
//                          ──► Goroutine Worker 3 ──►
//                                                          |
//                                                          v
//                                                   Respuesta al cliente
//
// CONCEPTOS DEMOSTRADOS:
// 1. Goroutines: Funciones que se ejecutan concurrentemente (workers).
// 2. Channels (chan): Comunicación segura entre goroutines sin memoria compartida.
// 3. sync.RWMutex: Protección de datos compartidos (inventario) ante accesos simultáneos.
// ============================================================================

// SolicitudPedido encapsula toda la información necesaria para que un Worker
// procese un pedido en segundo plano a través de un canal.
type SolicitudPedido struct {
	Cliente       modelos.Cliente           // Cliente que realiza la compra
	Elementos     []modelos.ElementoCarrito // Productos y cantidades del carrito
	CanalRespuesta chan ResultadoPedido      // Canal para devolver el resultado al handler
}

// ResultadoPedido contiene el resultado del procesamiento de un pedido por el Worker.
// Se envía de vuelta al handler HTTP a través del canal de respuesta.
type ResultadoPedido struct {
	Exito     bool           // Indica si el pedido se procesó correctamente
	Pedido    *modelos.Pedido // El pedido creado (nil si hubo error)
	Mensaje   string          // Mensaje descriptivo del resultado
}

// ProcesadorPedidos es el componente central de concurrencia.
// Administra un pool de goroutines (workers) que procesan pedidos desde un canal compartido.
type ProcesadorPedidos struct {
	canalSolicitudes chan SolicitudPedido // Canal por donde llegan las solicitudes de compra
	estado           *EstadoAPI           // Referencia al estado compartido (protegido por Mutex)
	cantidadWorkers  int                  // Número de goroutines trabajando en paralelo
}

// NuevoProcesadorPedidos crea e inicializa el procesador con el número de workers indicado.
// El tamaño del buffer del canal determina cuántas solicitudes pueden encolarse sin bloquear.
func NuevoProcesadorPedidos(estado *EstadoAPI, cantidadWorkers int) *ProcesadorPedidos {
	return &ProcesadorPedidos{
		canalSolicitudes: make(chan SolicitudPedido, 100), // Buffer de 100 solicitudes
		estado:           estado,
		cantidadWorkers:  cantidadWorkers,
	}
}

// Iniciar lanza las goroutines del Worker Pool.
// Cada worker es una goroutine independiente que escucha del mismo canal de solicitudes.
// Go distribuye automáticamente las solicitudes entre los workers disponibles.
func (p *ProcesadorPedidos) Iniciar() {
	for i := 1; i <= p.cantidadWorkers; i++ {
		go p.worker(i) // Cada llamada a 'go' lanza una nueva goroutine concurrente
	}
	log.Printf("[CONCURRENCIA] Worker Pool iniciado con %d workers procesando pedidos\n", p.cantidadWorkers)
}

// GetCanalSolicitudes expone el canal de solicitudes para que el handler pueda enviar pedidos.
func (p *ProcesadorPedidos) GetCanalSolicitudes() chan SolicitudPedido {
	return p.canalSolicitudes
}

// worker es la función que ejecuta cada goroutine del pool.
// Escucha continuamente del canal de solicitudes usando 'range', que bloquea hasta recibir datos.
// Cuando recibe una solicitud, la procesa y envía el resultado por el canal de respuesta.
func (p *ProcesadorPedidos) worker(idWorker int) {
	// 'range' sobre un canal: la goroutine se bloquea esperando datos.
	// Cuando llega una solicitud, la procesa. Si el canal se cierra, el for termina.
	for solicitud := range p.canalSolicitudes {
		log.Printf("[Worker %d] Procesando pedido para cliente: %s\n",
			idWorker, solicitud.Cliente.GetNombre())

		resultado := p.procesarPedido(solicitud, idWorker)

		// Enviar el resultado de vuelta al handler HTTP a través del canal de respuesta
		solicitud.CanalRespuesta <- resultado
	}
}

// procesarPedido ejecuta la lógica transaccional de checkout dentro del worker.
// Protege el acceso al inventario compartido usando sync.RWMutex para evitar Race Conditions.
func (p *ProcesadorPedidos) procesarPedido(solicitud SolicitudPedido, idWorker int) ResultadoPedido {
	// =========================================================================
	// PASO 1: Simulación de validación de pago (operación asíncrona)
	// =========================================================================
	// En un sistema real, aquí se contactaría a un gateway de pagos externo.
	// La simulación con time.Sleep demuestra que el worker puede tardar sin bloquear otros workers.
	log.Printf("[Worker %d] Validando pago simulado...\n", idWorker)
	time.Sleep(500 * time.Millisecond) // Simula latencia de pasarela de pagos
	log.Printf("[Worker %d] Pago validado correctamente\n", idWorker)

	// =========================================================================
	// PASO 2: Bloquear acceso exclusivo al inventario (sync.RWMutex)
	// =========================================================================
	// Usamos Lock() (escritura exclusiva) porque vamos a MODIFICAR el stock.
	// Mientras este worker tiene el lock, ningún otro worker ni handler puede
	// leer ni escribir el inventario, previniendo Race Conditions.
	p.estado.Mu.Lock()
	defer p.estado.Mu.Unlock()

	// =========================================================================
	// PASO 3: Verificar stock (dentro de la zona protegida por Mutex)
	// =========================================================================
	for _, elem := range solicitud.Elementos {
		codigoProducto := elem.GetProducto().GetCodigo()
		cantidadSolicitada := elem.GetCantidad()

		_, posInventario, existe := servicios.BuscarProductoPorCodigo(p.estado.Productos, codigoProducto)
		if !existe {
			return ResultadoPedido{
				Exito:   false,
				Mensaje: fmt.Sprintf("El producto '%s' ya no existe en el inventario", codigoProducto),
			}
		}

		productoReal := p.estado.Productos[posInventario]
		if !productoReal.TieneStockSuficiente(cantidadSolicitada) {
			return ResultadoPedido{
				Exito: false,
				Mensaje: fmt.Sprintf("Stock insuficiente para '%s'. Disponible: %d, Solicitado: %d",
					productoReal.GetNombre(), productoReal.GetCantidadDisponible(), cantidadSolicitada),
			}
		}
	}

	// =========================================================================
	// PASO 4: Descontar stock (operación protegida por Mutex)
	// =========================================================================
	for _, elem := range solicitud.Elementos {
		codigoProducto := elem.GetProducto().GetCodigo()
		_, posInventario, _ := servicios.BuscarProductoPorCodigo(p.estado.Productos, codigoProducto)

		err := p.estado.Productos[posInventario].DescontarStock(elem.GetCantidad())
		if err != nil {
			return ResultadoPedido{
				Exito:   false,
				Mensaje: fmt.Sprintf("Error al descontar stock de '%s': %v", codigoProducto, err),
			}
		}
	}

	// Persistir productos con stock actualizado
	err := servicios.GuardarProductos(p.estado.Productos)
	if err != nil {
		return ResultadoPedido{
			Exito:   false,
			Mensaje: "Error al guardar productos actualizados: " + err.Error(),
		}
	}

	// =========================================================================
	// PASO 5: Crear el Pedido inmutable y persistir
	// =========================================================================
	codigoPedido := fmt.Sprintf("PED-%04d", len(p.estado.Pedidos)+1)
	fechaActual := time.Now().Format("2006-01-02 15:04:05")

	nuevoPedido, err := modelos.NuevoPedido(codigoPedido, solicitud.Cliente, solicitud.Elementos, fechaActual)
	if err != nil {
		return ResultadoPedido{
			Exito:   false,
			Mensaje: "Error al crear el pedido: " + err.Error(),
		}
	}

	p.estado.Pedidos = append(p.estado.Pedidos, *nuevoPedido)
	err = servicios.GuardarPedidos(p.estado.Pedidos)
	if err != nil {
		return ResultadoPedido{
			Exito:   false,
			Mensaje: "Error al guardar pedido: " + err.Error(),
		}
	}

	// =========================================================================
	// PASO 6: Emitir comprobante (simulación de notificación)
	// =========================================================================
	log.Printf("[Worker %d] ¡Pedido %s completado! Total: $%.2f | Cliente: %s\n",
		idWorker, nuevoPedido.GetCodigo(), nuevoPedido.GetTotal(), solicitud.Cliente.GetNombre())

	return ResultadoPedido{
		Exito:   true,
		Pedido:  nuevoPedido,
		Mensaje: "¡Compra procesada exitosamente por Worker " + fmt.Sprintf("%d", idWorker) + "!",
	}
}
