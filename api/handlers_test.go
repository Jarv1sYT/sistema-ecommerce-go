package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sistema-ecommerce-go/modelos"
)

// ============================================================================
// PRUEBAS DE INTEGRACIÓN HTTP PARA LOS ENDPOINTS REST (Unidad 4 - Testing)
// ============================================================================
// Ejecutar: go test -v ./api/ -run TestHandler
// ============================================================================
//
// Estas pruebas utilizan httptest.NewRecorder() del paquete estándar de Go
// para simular peticiones HTTP reales sin necesidad de levantar un servidor.

// crearEstadoPrueba genera un EstadoAPI limpio con datos de prueba precargados.
// Se usa en cada test para garantizar independencia entre pruebas.
func crearEstadoPrueba() *EstadoAPI {
	p1, _ := modelos.NuevoProducto("P001", "Laptop HP", 899.99, 10)
	p2, _ := modelos.NuevoProducto("P002", "Mouse Logitech", 25.50, 50)

	c1, _ := modelos.NuevoCliente("1234567890", "Juan Pérez", "juan@email.com")

	estado := &EstadoAPI{
		Productos: []modelos.Producto{*p1, *p2},
		Clientes:  []modelos.Cliente{*c1},
		Pedidos:   []modelos.Pedido{},
		Carrito:   modelos.NuevoCarrito(),
	}

	// Inicializar el procesador de pedidos para tests que usen checkout
	estado.Procesador = NuevoProcesadorPedidos(estado, 1)
	estado.Procesador.Iniciar()

	return estado
}

// ejecutarRequest es una función auxiliar que ejecuta una petición HTTP simulada
// y retorna el ResponseRecorder para inspeccionar la respuesta.
func ejecutarRequest(mux *http.ServeMux, metodo, ruta string, cuerpo interface{}) *httptest.ResponseRecorder {
	var req *http.Request

	if cuerpo != nil {
		datosJSON, _ := json.Marshal(cuerpo)
		req = httptest.NewRequest(metodo, ruta, bytes.NewBuffer(datosJSON))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(metodo, ruta, nil)
	}

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)
	return recorder
}

// decodificarRespuesta extrae y parsea el JSON de respuesta en una RespuestaJSON.
func decodificarRespuesta(t *testing.T, recorder *httptest.ResponseRecorder) RespuestaJSON {
	var respuesta RespuestaJSON
	err := json.NewDecoder(recorder.Body).Decode(&respuesta)
	if err != nil {
		t.Fatalf("Error al decodificar respuesta JSON: %v", err)
	}
	return respuesta
}

// ============================================================================
// ENDPOINT 1: GET /api/productos — Listar productos
// ============================================================================

func TestHandler_ListarProductos_OK(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	recorder := ejecutarRequest(mux, "GET", "/api/productos", nil)

	if recorder.Code != http.StatusOK {
		t.Errorf("Código HTTP esperado %d, obtenido %d", http.StatusOK, recorder.Code)
	}

	respuesta := decodificarRespuesta(t, recorder)
	if !respuesta.Exito {
		t.Error("Se esperaba respuesta exitosa al listar productos")
	}
}

// ============================================================================
// ENDPOINT 2: GET /api/productos/{id} — Detalle de producto
// ============================================================================

func TestHandler_ObtenerProducto_Existente(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	recorder := ejecutarRequest(mux, "GET", "/api/productos/P001", nil)

	if recorder.Code != http.StatusOK {
		t.Errorf("Código HTTP esperado %d, obtenido %d", http.StatusOK, recorder.Code)
	}

	respuesta := decodificarRespuesta(t, recorder)
	if !respuesta.Exito {
		t.Error("Se esperaba respuesta exitosa al buscar producto existente")
	}
}

func TestHandler_ObtenerProducto_Inexistente(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	recorder := ejecutarRequest(mux, "GET", "/api/productos/P999", nil)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Código HTTP esperado %d, obtenido %d", http.StatusNotFound, recorder.Code)
	}

	respuesta := decodificarRespuesta(t, recorder)
	if respuesta.Exito {
		t.Error("Se esperaba respuesta de error al buscar producto inexistente")
	}
}

// ============================================================================
// ENDPOINT 3: POST /api/productos — Crear producto
// ============================================================================

func TestHandler_CrearProducto_Valido(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	nuevoProducto := ProductoRequest{
		Codigo:   "P003",
		Nombre:   "Teclado Mecánico",
		Precio:   75.00,
		Cantidad: 20,
	}

	recorder := ejecutarRequest(mux, "POST", "/api/productos", nuevoProducto)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Código HTTP esperado %d, obtenido %d", http.StatusCreated, recorder.Code)
	}

	respuesta := decodificarRespuesta(t, recorder)
	if !respuesta.Exito {
		t.Errorf("Se esperaba respuesta exitosa, obtenido: %s", respuesta.Mensaje)
	}

	// Verificar que el producto fue agregado al estado
	if len(estado.Productos) != 3 {
		t.Errorf("Se esperaban 3 productos después de crear uno nuevo, obtenido %d", len(estado.Productos))
	}
}

func TestHandler_CrearProducto_Duplicado(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	productoDuplicado := ProductoRequest{
		Codigo:   "P001", // Ya existe en los datos de prueba
		Nombre:   "Duplicado",
		Precio:   10.00,
		Cantidad: 5,
	}

	recorder := ejecutarRequest(mux, "POST", "/api/productos", productoDuplicado)

	if recorder.Code != http.StatusConflict {
		t.Errorf("Código HTTP esperado %d (conflicto), obtenido %d", http.StatusConflict, recorder.Code)
	}
}

func TestHandler_CrearProducto_PrecioNegativo(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	productoInvalido := ProductoRequest{
		Codigo:   "P099",
		Nombre:   "Inválido",
		Precio:   -5.00,
		Cantidad: 1,
	}

	recorder := ejecutarRequest(mux, "POST", "/api/productos", productoInvalido)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Código HTTP esperado %d (bad request), obtenido %d", http.StatusBadRequest, recorder.Code)
	}
}

// ============================================================================
// ENDPOINT 4: GET /api/clientes — Listar clientes
// ============================================================================

func TestHandler_ListarClientes_OK(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	recorder := ejecutarRequest(mux, "GET", "/api/clientes", nil)

	if recorder.Code != http.StatusOK {
		t.Errorf("Código HTTP esperado %d, obtenido %d", http.StatusOK, recorder.Code)
	}

	respuesta := decodificarRespuesta(t, recorder)
	if !respuesta.Exito {
		t.Error("Se esperaba respuesta exitosa al listar clientes")
	}
}

// ============================================================================
// ENDPOINT 5: POST /api/clientes — Crear cliente
// ============================================================================

func TestHandler_CrearCliente_Valido(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	nuevoCliente := ClienteRequest{
		Identificacion: "0987654321",
		Nombre:         "María López",
		Correo:         "maria@email.com",
	}

	recorder := ejecutarRequest(mux, "POST", "/api/clientes", nuevoCliente)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Código HTTP esperado %d, obtenido %d", http.StatusCreated, recorder.Code)
	}

	if len(estado.Clientes) != 2 {
		t.Errorf("Se esperaban 2 clientes después de crear uno nuevo, obtenido %d", len(estado.Clientes))
	}
}

func TestHandler_CrearCliente_Duplicado(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	clienteDuplicado := ClienteRequest{
		Identificacion: "1234567890", // Ya existe
		Nombre:         "Duplicado",
		Correo:         "dup@email.com",
	}

	recorder := ejecutarRequest(mux, "POST", "/api/clientes", clienteDuplicado)

	if recorder.Code != http.StatusConflict {
		t.Errorf("Código HTTP esperado %d (conflicto), obtenido %d", http.StatusConflict, recorder.Code)
	}
}

// ============================================================================
// ENDPOINT 6: POST /api/carrito/agregar — Agregar al carrito
// ============================================================================

func TestHandler_AgregarAlCarrito_OK(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	reqCarrito := CarritoRequest{
		CodigoProducto: "P001",
		Cantidad:       2,
	}

	recorder := ejecutarRequest(mux, "POST", "/api/carrito/agregar", reqCarrito)

	if recorder.Code != http.StatusOK {
		t.Errorf("Código HTTP esperado %d, obtenido %d", http.StatusOK, recorder.Code)
	}

	respuesta := decodificarRespuesta(t, recorder)
	if !respuesta.Exito {
		t.Errorf("Se esperaba respuesta exitosa, obtenido: %s", respuesta.Mensaje)
	}

	// Verificar que el carrito ya no está vacío
	if estado.Carrito.EsVacio() {
		t.Error("El carrito no debe estar vacío después de agregar un producto")
	}
}

func TestHandler_AgregarAlCarrito_ProductoInexistente(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	reqCarrito := CarritoRequest{
		CodigoProducto: "P999",
		Cantidad:       1,
	}

	recorder := ejecutarRequest(mux, "POST", "/api/carrito/agregar", reqCarrito)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Código HTTP esperado %d, obtenido %d", http.StatusNotFound, recorder.Code)
	}
}

func TestHandler_AgregarAlCarrito_StockInsuficiente(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	reqCarrito := CarritoRequest{
		CodigoProducto: "P001",
		Cantidad:       99999, // Más que el stock disponible (10)
	}

	recorder := ejecutarRequest(mux, "POST", "/api/carrito/agregar", reqCarrito)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Código HTTP esperado %d, obtenido %d", http.StatusBadRequest, recorder.Code)
	}
}

// ============================================================================
// ENDPOINT 7: POST /api/pedidos/checkout — Checkout
// ============================================================================

func TestHandler_Checkout_CarritoVacio(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	reqCheckout := CheckoutRequest{
		IdentificacionCliente: "1234567890",
	}

	recorder := ejecutarRequest(mux, "POST", "/api/pedidos/checkout", reqCheckout)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Código HTTP esperado %d (carrito vacío), obtenido %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestHandler_Checkout_ClienteInexistente(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	// Primero agregar algo al carrito
	estado.Carrito.AgregarElemento(estado.Productos[0], 1)

	reqCheckout := CheckoutRequest{
		IdentificacionCliente: "0000000000", // No existe
	}

	recorder := ejecutarRequest(mux, "POST", "/api/pedidos/checkout", reqCheckout)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Código HTTP esperado %d (cliente no encontrado), obtenido %d", http.StatusNotFound, recorder.Code)
	}
}

// ============================================================================
// ENDPOINT 8: GET /api/pedidos — Listar pedidos
// ============================================================================

func TestHandler_ListarPedidos_Vacio(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	recorder := ejecutarRequest(mux, "GET", "/api/pedidos", nil)

	if recorder.Code != http.StatusOK {
		t.Errorf("Código HTTP esperado %d, obtenido %d", http.StatusOK, recorder.Code)
	}

	respuesta := decodificarRespuesta(t, recorder)
	if !respuesta.Exito {
		t.Error("Se esperaba respuesta exitosa al listar pedidos")
	}
}

// ============================================================================
// PRUEBA DE FLUJO COMPLETO (Agregar al carrito → Checkout → Verificar stock)
// ============================================================================

func TestHandler_FlujoCompleto_AgregarYCheckout(t *testing.T) {
	estado := crearEstadoPrueba()
	mux := http.NewServeMux()
	ConfigurarRutas(mux, estado)

	// Paso 1: Agregar producto al carrito
	reqCarrito := CarritoRequest{
		CodigoProducto: "P002", // Mouse Logitech, stock inicial: 50
		Cantidad:       5,
	}
	recorder := ejecutarRequest(mux, "POST", "/api/carrito/agregar", reqCarrito)
	if recorder.Code != http.StatusOK {
		t.Fatalf("Error al agregar al carrito: HTTP %d", recorder.Code)
	}

	// Paso 2: Hacer checkout
	reqCheckout := CheckoutRequest{
		IdentificacionCliente: "1234567890",
	}
	recorder = ejecutarRequest(mux, "POST", "/api/pedidos/checkout", reqCheckout)
	if recorder.Code != http.StatusCreated {
		respuesta := decodificarRespuesta(t, recorder)
		t.Fatalf("Error en checkout: HTTP %d - %s", recorder.Code, respuesta.Mensaje)
	}

	// Paso 3: Verificar que el stock se descontó correctamente
	// Stock inicial: 50, Comprado: 5 → Stock esperado: 45
	stockFinal := estado.Productos[1].GetCantidadDisponible()
	if stockFinal != 45 {
		t.Errorf("Stock después del checkout esperado 45, obtenido %d", stockFinal)
	}

	// Paso 4: Verificar que se creó un pedido
	if len(estado.Pedidos) != 1 {
		t.Errorf("Se esperaba 1 pedido después del checkout, obtenido %d", len(estado.Pedidos))
	}

	// Paso 5: Verificar que el carrito quedó vacío
	if !estado.Carrito.EsVacio() {
		t.Error("El carrito debe estar vacío después de un checkout exitoso")
	}
}
