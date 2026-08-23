package api

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"sistema-ecommerce-go/modelos"
)

// EstadoAPI contiene los datos compartidos entre todos los handlers de la API.
// Centraliza el acceso a los slices de datos y el carrito en un solo lugar.
// Se usa sync.RWMutex para proteger el acceso concurrente (Unidad 4 - Concurrencia).
type EstadoAPI struct {
	Mu         sync.RWMutex          // Mutex de lectura/escritura para proteger los datos
	Productos  []modelos.Producto    // Inventario de productos
	Clientes   []modelos.Cliente     // Listado de clientes registrados
	Pedidos    []modelos.Pedido      // Historial de pedidos
	Carrito    *modelos.Carrito      // Carrito de compras en memoria
	Procesador *ProcesadorPedidos    // Worker Pool para procesar pedidos concurrentemente
}

// NuevoEstadoAPI crea una instancia del estado compartido de la API.
// Recibe los datos cargados desde JSON al arrancar la aplicación.
func NuevoEstadoAPI(productos []modelos.Producto, clientes []modelos.Cliente, pedidos []modelos.Pedido) *EstadoAPI {
	estado := &EstadoAPI{
		Productos: productos,
		Clientes:  clientes,
		Pedidos:   pedidos,
		Carrito:   modelos.NuevoCarrito(),
	}

	// Inicializar el Worker Pool con 3 goroutines concurrentes
	// Cada worker es una goroutine independiente que procesa pedidos desde un canal compartido
	estado.Procesador = NuevoProcesadorPedidos(estado, 3)

	return estado
}

// ConfigurarRutas registra los 8 endpoints REST en el multiplexor HTTP proporcionado.
// Utiliza el enrutamiento avanzado de Go 1.22+ que permite especificar método y parámetros.
func ConfigurarRutas(mux *http.ServeMux, estado *EstadoAPI) {
	// --- ENDPOINTS DE PRODUCTOS (3) ---
	mux.HandleFunc("GET /api/productos", estado.ListarProductosHandler)        // 1. Listar catálogo
	mux.HandleFunc("GET /api/productos/{id}", estado.ObtenerProductoHandler)   // 2. Detalle/stock de un producto
	mux.HandleFunc("POST /api/productos", estado.CrearProductoHandler)         // 3. Registrar nuevo producto

	// --- ENDPOINTS DE CLIENTES (2) ---
	mux.HandleFunc("GET /api/clientes", estado.ListarClientesHandler)          // 4. Listar clientes
	mux.HandleFunc("POST /api/clientes", estado.CrearClienteHandler)           // 5. Registrar cliente

	// --- ENDPOINT DE CARRITO (1) ---
	mux.HandleFunc("POST /api/carrito/agregar", estado.AgregarAlCarritoHandler) // 6. Añadir elemento al carrito

	// --- ENDPOINTS DE PEDIDOS (2) ---
	mux.HandleFunc("POST /api/pedidos/checkout", estado.CheckoutHandler)       // 7. Procesar compra (usa Worker Pool)
	mux.HandleFunc("GET /api/pedidos", estado.ListarPedidosHandler)            // 8. Historial de órdenes

	// --- INTERFAZ WEB (Archivos estáticos desde carpeta public/) ---
	archivoEstatico := http.FileServer(http.Dir("public"))
	mux.Handle("/", archivoEstatico)
}

// IniciarServidor configura y arranca el servidor HTTP en el puerto especificado.
// Se ejecuta como goroutine desde main.go para no bloquear la consola interactiva.
func IniciarServidor(estado *EstadoAPI, puerto string) {
	// Iniciar el Worker Pool de procesamiento concurrente de pedidos
	estado.Procesador.Iniciar()

	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	direccion := ":" + puerto
	fmt.Printf("\n[API REST] Servidor iniciado en http://localhost%s\n", direccion)
	fmt.Println("[API REST] Interfaz Web:  http://localhost" + direccion)
	fmt.Println("[API REST] Endpoints disponibles:")
	fmt.Println("  GET  /api/productos           - Listar catálogo de productos")
	fmt.Println("  GET  /api/productos/{id}       - Consultar detalle de un producto")
	fmt.Println("  POST /api/productos            - Registrar nuevo producto (JSON)")
	fmt.Println("  GET  /api/clientes             - Listar clientes registrados")
	fmt.Println("  POST /api/clientes             - Registrar nuevo cliente (JSON)")
	fmt.Println("  POST /api/carrito/agregar       - Añadir producto al carrito (JSON)")
	fmt.Println("  POST /api/pedidos/checkout      - Confirmar compra / Checkout (JSON)")
	fmt.Println("  GET  /api/pedidos              - Listar historial de pedidos")
	fmt.Println("[CONCURRENCIA] Worker Pool activo con 3 goroutines procesando pedidos")
	fmt.Println()

	err := http.ListenAndServe(direccion, mux)
	if err != nil {
		log.Fatalf("[API REST] Error al iniciar el servidor: %v\n", err)
	}
}
