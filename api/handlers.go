package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"sistema-ecommerce-go/modelos"
	"sistema-ecommerce-go/servicios"
)

// ----------------------------------------------------------------------------
// ESTRUCTURAS DTO PARA RECIBIR DATOS JSON DESDE EL CLIENTE HTTP
// ----------------------------------------------------------------------------

// ProductoRequest representa el cuerpo JSON esperado al crear un producto vía POST.
type ProductoRequest struct {
	Codigo   string  `json:"codigo"`
	Nombre   string  `json:"nombre"`
	Precio   float64 `json:"precio"`
	Cantidad int     `json:"cantidad"`
}

// ClienteRequest representa el cuerpo JSON esperado al registrar un cliente vía POST.
type ClienteRequest struct {
	Identificacion string `json:"identificacion"`
	Nombre         string `json:"nombre"`
	Correo         string `json:"correo"`
}

// CarritoRequest representa el cuerpo JSON esperado al agregar un producto al carrito vía POST.
type CarritoRequest struct {
	CodigoProducto string `json:"codigo_producto"`
	Cantidad       int    `json:"cantidad"`
}

// CheckoutRequest representa el cuerpo JSON esperado al confirmar una compra vía POST.
type CheckoutRequest struct {
	IdentificacionCliente string `json:"identificacion_cliente"`
}

// RespuestaJSON es una estructura genérica para respuestas de la API.
type RespuestaJSON struct {
	Exito   bool        `json:"exito"`
	Mensaje string      `json:"mensaje"`
	Datos   interface{} `json:"datos,omitempty"`
}

// ----------------------------------------------------------------------------
// FUNCIONES AUXILIARES PARA RESPUESTAS HTTP
// ----------------------------------------------------------------------------

// responderJSON serializa y escribe una respuesta JSON con el código de estado indicado.
func responderJSON(w http.ResponseWriter, codigo int, respuesta RespuestaJSON) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(respuesta)
}

// responderError envía una respuesta de error JSON con el código HTTP correspondiente.
func responderError(w http.ResponseWriter, codigo int, mensaje string) {
	responderJSON(w, codigo, RespuestaJSON{
		Exito:   false,
		Mensaje: mensaje,
	})
}

// responderExito envía una respuesta exitosa JSON con datos opcionales.
func responderExito(w http.ResponseWriter, codigo int, mensaje string, datos interface{}) {
	responderJSON(w, codigo, RespuestaJSON{
		Exito:   true,
		Mensaje: mensaje,
		Datos:   datos,
	})
}

// ============================================================================
// ENDPOINT 1: GET /api/productos — Listar catálogo de productos
// ============================================================================

func (e *EstadoAPI) ListarProductosHandler(w http.ResponseWriter, r *http.Request) {
	e.Mu.RLock()
	defer e.Mu.RUnlock()

	responderExito(w, http.StatusOK, "Listado de productos obtenido correctamente", e.Productos)
}

// ============================================================================
// ENDPOINT 2: GET /api/productos/{id} — Consultar detalle/stock de un producto
// ============================================================================

func (e *EstadoAPI) ObtenerProductoHandler(w http.ResponseWriter, r *http.Request) {
	e.Mu.RLock()
	defer e.Mu.RUnlock()

	// Extraer el parámetro {id} de la URL usando el enrutador nativo de Go 1.22+
	id := r.PathValue("id")
	if id == "" {
		responderError(w, http.StatusBadRequest, "Debe proporcionar un código de producto en la URL")
		return
	}

	producto, _, encontrado := servicios.BuscarProductoPorCodigo(e.Productos, id)
	if !encontrado {
		responderError(w, http.StatusNotFound, fmt.Sprintf("No se encontró un producto con el código '%s'", id))
		return
	}

	responderExito(w, http.StatusOK, "Producto encontrado", producto)
}

// ============================================================================
// ENDPOINT 3: POST /api/productos — Registrar nuevo producto
// ============================================================================

func (e *EstadoAPI) CrearProductoHandler(w http.ResponseWriter, r *http.Request) {
	var req ProductoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responderError(w, http.StatusBadRequest, "Error al leer el cuerpo JSON: "+err.Error())
		return
	}

	e.Mu.Lock()
	defer e.Mu.Unlock()

	// Verificar que no exista un producto con el mismo código
	_, _, existe := servicios.BuscarProductoPorCodigo(e.Productos, req.Codigo)
	if existe {
		responderError(w, http.StatusConflict, "Ya existe un producto registrado con ese código")
		return
	}

	// Crear el producto usando el constructor encapsulado del modelo
	nuevoProducto, err := modelos.NuevoProducto(req.Codigo, req.Nombre, req.Precio, req.Cantidad)
	if err != nil {
		responderError(w, http.StatusBadRequest, "Error de validación: "+err.Error())
		return
	}

	e.Productos = append(e.Productos, *nuevoProducto)

	// Persistir en JSON
	err = servicios.GuardarProductos(e.Productos)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al guardar en archivo JSON: "+err.Error())
		return
	}

	responderExito(w, http.StatusCreated, "Producto registrado exitosamente", nuevoProducto)
}

// ============================================================================
// ENDPOINT 4: GET /api/clientes — Listar clientes registrados
// ============================================================================

func (e *EstadoAPI) ListarClientesHandler(w http.ResponseWriter, r *http.Request) {
	e.Mu.RLock()
	defer e.Mu.RUnlock()

	responderExito(w, http.StatusOK, "Listado de clientes obtenido correctamente", e.Clientes)
}

// ============================================================================
// ENDPOINT 5: POST /api/clientes — Registrar nuevo cliente
// ============================================================================

func (e *EstadoAPI) CrearClienteHandler(w http.ResponseWriter, r *http.Request) {
	var req ClienteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responderError(w, http.StatusBadRequest, "Error al leer el cuerpo JSON: "+err.Error())
		return
	}

	e.Mu.Lock()
	defer e.Mu.Unlock()

	// Verificar que no exista un cliente con la misma identificación
	_, _, existe := servicios.BuscarClientePorIdentificacion(e.Clientes, req.Identificacion)
	if existe {
		responderError(w, http.StatusConflict, "Ya existe un cliente registrado con esa identificación")
		return
	}

	// Crear el cliente usando el constructor encapsulado del modelo
	nuevoCliente, err := modelos.NuevoCliente(req.Identificacion, req.Nombre, req.Correo)
	if err != nil {
		responderError(w, http.StatusBadRequest, "Error de validación: "+err.Error())
		return
	}

	e.Clientes = append(e.Clientes, *nuevoCliente)

	// Persistir en JSON
	err = servicios.GuardarClientes(e.Clientes)
	if err != nil {
		responderError(w, http.StatusInternalServerError, "Error al guardar en archivo JSON: "+err.Error())
		return
	}

	responderExito(w, http.StatusCreated, "Cliente registrado exitosamente", nuevoCliente)
}

// ============================================================================
// ENDPOINT 6: POST /api/carrito/agregar — Añadir producto al carrito
// ============================================================================

func (e *EstadoAPI) AgregarAlCarritoHandler(w http.ResponseWriter, r *http.Request) {
	var req CarritoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responderError(w, http.StatusBadRequest, "Error al leer el cuerpo JSON: "+err.Error())
		return
	}

	e.Mu.Lock()
	defer e.Mu.Unlock()

	// Buscar el producto en el inventario
	producto, _, encontrado := servicios.BuscarProductoPorCodigo(e.Productos, req.CodigoProducto)
	if !encontrado {
		responderError(w, http.StatusNotFound,
			fmt.Sprintf("No se encontró un producto con el código '%s'", req.CodigoProducto))
		return
	}

	// Delegar la lógica de agregar al método encapsulado del Carrito
	err = e.Carrito.AgregarElemento(producto, req.Cantidad)
	if err != nil {
		responderError(w, http.StatusBadRequest, "No se pudo agregar al carrito: "+err.Error())
		return
	}

	// Estructura de respuesta con el estado actual del carrito
	carritoResumen := map[string]interface{}{
		"elementos":       e.Carrito.GetElementos(),
		"total":           e.Carrito.GetTotal(),
		"cantidad_items":  len(e.Carrito.GetElementos()),
	}

	responderExito(w, http.StatusOK,
		fmt.Sprintf("Producto '%s' (x%d) agregado al carrito", producto.GetNombre(), req.Cantidad),
		carritoResumen)
}

// ============================================================================
// ENDPOINT 7: POST /api/pedidos/checkout — Procesar compra (Checkout)
// ============================================================================
// Este endpoint demuestra CONCURRENCIA (Unidad 4):
// 1. El handler valida las precondiciones básicas (carrito no vacío, cliente válido).
// 2. Envía la solicitud al Worker Pool a través de un CHANNEL (canal de Go).
// 3. Una GOROUTINE worker recoge la solicitud del canal y la procesa en segundo plano.
// 4. El worker protege el inventario con sync.RWMutex para evitar RACE CONDITIONS.
// 5. El resultado se devuelve al handler por otro channel de respuesta.

func (e *EstadoAPI) CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	var req CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		responderError(w, http.StatusBadRequest, "Error al leer el cuerpo JSON: "+err.Error())
		return
	}

	// Validación previa (fuera del worker para responder rápido si hay error básico)
	e.Mu.RLock()
	carritoVacio := e.Carrito.EsVacio()
	e.Mu.RUnlock()

	if carritoVacio {
		responderError(w, http.StatusBadRequest, "El carrito está vacío. Agregue productos antes de confirmar la compra")
		return
	}

	// Buscar el cliente
	e.Mu.RLock()
	cliente, _, encontrado := servicios.BuscarClientePorIdentificacion(e.Clientes, req.IdentificacionCliente)
	elementos := e.Carrito.GetElementos()
	e.Mu.RUnlock()

	if !encontrado {
		responderError(w, http.StatusNotFound,
			fmt.Sprintf("No se encontró un cliente con la identificación '%s'", req.IdentificacionCliente))
		return
	}

	// =========================================================================
	// ENVIAR SOLICITUD AL WORKER POOL A TRAVÉS DEL CANAL (CHANNEL)
	// =========================================================================
	// Se crea un canal de respuesta exclusivo para esta solicitud.
	// El handler se bloquea esperando el resultado, pero el procesamiento
	// lo realiza una goroutine del pool en segundo plano.
	canalRespuesta := make(chan ResultadoPedido, 1)

	solicitud := SolicitudPedido{
		Cliente:        cliente,
		Elementos:      elementos,
		CanalRespuesta: canalRespuesta,
	}

	// Enviar la solicitud al canal del Worker Pool
	// Un worker goroutine la recogerá y procesará concurrentemente
	e.Procesador.GetCanalSolicitudes() <- solicitud

	// Esperar la respuesta del worker (el canal bloquea hasta recibir el resultado)
	resultado := <-canalRespuesta
	close(canalRespuesta)

	// Vaciar el carrito si la compra fue exitosa
	if resultado.Exito {
		e.Mu.Lock()
		e.Carrito.Vaciar()
		e.Mu.Unlock()

		comprobante := map[string]interface{}{
			"codigo_pedido": resultado.Pedido.GetCodigo(),
			"cliente":       resultado.Pedido.GetCliente().GetNombre(),
			"elementos":     resultado.Pedido.GetElementos(),
			"total":         resultado.Pedido.GetTotal(),
			"fecha":         resultado.Pedido.GetFecha(),
		}

		responderExito(w, http.StatusCreated, resultado.Mensaje, comprobante)
	} else {
		responderError(w, http.StatusConflict, resultado.Mensaje)
	}
}

// ============================================================================
// ENDPOINT 8: GET /api/pedidos — Listar historial de pedidos
// ============================================================================

func (e *EstadoAPI) ListarPedidosHandler(w http.ResponseWriter, r *http.Request) {
	e.Mu.RLock()
	defer e.Mu.RUnlock()

	// Permitir filtrar pedidos por cliente usando query string: ?cliente=1234567890
	clienteID := strings.TrimSpace(r.URL.Query().Get("cliente"))

	if clienteID != "" {
		pedidosFiltrados := []modelos.Pedido{}
		for _, p := range e.Pedidos {
			if strings.EqualFold(p.GetCliente().GetIdentificacion(), clienteID) {
				pedidosFiltrados = append(pedidosFiltrados, p)
			}
		}
		responderExito(w, http.StatusOK,
			fmt.Sprintf("Pedidos encontrados para el cliente '%s'", clienteID),
			pedidosFiltrados)
		return
	}

	responderExito(w, http.StatusOK, "Historial de pedidos obtenido correctamente", e.Pedidos)
}
