package servicios

import (
	"testing"

	"sistema-ecommerce-go/modelos"
)

// ============================================================================
// PRUEBAS UNITARIAS PARA EL SERVICIO DE PEDIDOS Y CHECKOUT (Unidad 4 - Testing)
// ============================================================================
// Ejecutar: go test -v ./servicios/ -run TestPedido
// ============================================================================

// --- PRUEBAS DEL CONSTRUCTOR NuevoPedido (Inmutabilidad) ---

// TestPedido_CrearPedidoValido verifica que un pedido se crea correctamente con datos válidos.
func TestPedido_CrearPedidoValido(t *testing.T) {
	// Preparar datos de prueba
	producto, _ := modelos.NuevoProducto("P001", "Laptop HP", 899.99, 10)
	cliente, _ := modelos.NuevoCliente("1234567890", "Juan Pérez", "juan@email.com")
	elemento, _ := modelos.NuevoElementoCarrito(*producto, 2)

	elementos := []modelos.ElementoCarrito{*elemento}

	pedido, err := modelos.NuevoPedido("PED-0001", *cliente, elementos, "2026-08-21 12:00:00")

	if err != nil {
		t.Fatalf("Error inesperado al crear pedido válido: %v", err)
	}
	if pedido == nil {
		t.Fatal("El pedido creado no debe ser nil")
	}
	if pedido.GetCodigo() != "PED-0001" {
		t.Errorf("Código esperado 'PED-0001', obtenido '%s'", pedido.GetCodigo())
	}
	if pedido.GetCliente().GetNombre() != "Juan Pérez" {
		t.Errorf("Cliente esperado 'Juan Pérez', obtenido '%s'", pedido.GetCliente().GetNombre())
	}

	// Verificar cálculo automático del total: 899.99 * 2 = 1799.98
	totalEsperado := 1799.98
	if pedido.GetTotal() != totalEsperado {
		t.Errorf("Total esperado %.2f, obtenido %.2f", totalEsperado, pedido.GetTotal())
	}
}

// TestPedido_SinElementos verifica que NO se permita crear un pedido sin productos.
func TestPedido_SinElementos(t *testing.T) {
	cliente, _ := modelos.NuevoCliente("1234567890", "Juan Pérez", "juan@email.com")
	elementosVacios := []modelos.ElementoCarrito{}

	_, err := modelos.NuevoPedido("PED-0002", *cliente, elementosVacios, "2026-08-21 12:00:00")

	if err == nil {
		t.Fatal("Se esperaba un error al crear un pedido sin productos")
	}
}

// TestPedido_CodigoVacio verifica que NO se permita crear un pedido sin código.
func TestPedido_CodigoVacio(t *testing.T) {
	producto, _ := modelos.NuevoProducto("P001", "Laptop", 899.99, 10)
	cliente, _ := modelos.NuevoCliente("1234567890", "Juan", "juan@email.com")
	elemento, _ := modelos.NuevoElementoCarrito(*producto, 1)

	_, err := modelos.NuevoPedido("", *cliente, []modelos.ElementoCarrito{*elemento}, "2026-08-21 12:00:00")

	if err == nil {
		t.Fatal("Se esperaba un error al crear un pedido con código vacío")
	}
}

// TestPedido_FechaVacia verifica que NO se permita crear un pedido sin fecha.
func TestPedido_FechaVacia(t *testing.T) {
	producto, _ := modelos.NuevoProducto("P001", "Laptop", 899.99, 10)
	cliente, _ := modelos.NuevoCliente("1234567890", "Juan", "juan@email.com")
	elemento, _ := modelos.NuevoElementoCarrito(*producto, 1)

	_, err := modelos.NuevoPedido("PED-0003", *cliente, []modelos.ElementoCarrito{*elemento}, "")

	if err == nil {
		t.Fatal("Se esperaba un error al crear un pedido con fecha vacía")
	}
}

// TestPedido_InmutabilidadElementos verifica que modificar el slice externo NO afecte el pedido.
func TestPedido_InmutabilidadElementos(t *testing.T) {
	producto, _ := modelos.NuevoProducto("P001", "Laptop", 899.99, 10)
	cliente, _ := modelos.NuevoCliente("1234567890", "Juan", "juan@email.com")
	elemento, _ := modelos.NuevoElementoCarrito(*producto, 2)

	elementosOriginales := []modelos.ElementoCarrito{*elemento}
	pedido, _ := modelos.NuevoPedido("PED-0004", *cliente, elementosOriginales, "2026-08-21 12:00:00")

	// Vaciar el slice original externo
	elementosOriginales = elementosOriginales[:0]

	// El pedido debe seguir teniendo sus elementos intactos (copia defensiva)
	elementosDelPedido := pedido.GetElementos()
	if len(elementosDelPedido) != 1 {
		t.Errorf("El pedido debe mantener sus elementos a pesar de modificaciones externas. Esperado 1, obtenido %d", len(elementosDelPedido))
	}
}

// --- PRUEBAS DE CHECKOUT (Descuento de Stock Transaccional) ---

// TestCheckout_DescontarStockCorrectamente verifica que el checkout descuenta el stock exacto.
func TestCheckout_DescontarStockCorrectamente(t *testing.T) {
	// Crear producto con stock de 10
	producto, _ := modelos.NuevoProducto("P001", "Laptop HP", 899.99, 10)

	// Simular compra de 3 unidades
	cantidadComprada := 3
	err := producto.DescontarStock(cantidadComprada)

	if err != nil {
		t.Fatalf("Error inesperado al descontar stock: %v", err)
	}

	stockEsperado := 7
	if producto.GetCantidadDisponible() != stockEsperado {
		t.Errorf("Stock después del checkout esperado %d, obtenido %d", stockEsperado, producto.GetCantidadDisponible())
	}
}

// TestCheckout_FallaStockInsuficiente verifica que el checkout falle si no hay suficiente stock.
func TestCheckout_FallaStockInsuficiente(t *testing.T) {
	producto, _ := modelos.NuevoProducto("P002", "Mouse", 25.00, 2)

	// Intentar comprar más de lo que hay
	err := producto.DescontarStock(5)

	if err == nil {
		t.Fatal("Se esperaba un error al intentar descontar más stock del disponible")
	}

	// El stock debe permanecer intacto
	if producto.GetCantidadDisponible() != 2 {
		t.Errorf("El stock no debe cambiar cuando el checkout falla. Esperado 2, obtenido %d", producto.GetCantidadDisponible())
	}
}

// TestCheckout_MultiplesProductos verifica el descuento correcto de múltiples productos.
func TestCheckout_MultiplesProductos(t *testing.T) {
	p1, _ := modelos.NuevoProducto("P001", "Laptop", 899.99, 10)
	p2, _ := modelos.NuevoProducto("P002", "Mouse", 25.50, 50)
	p3, _ := modelos.NuevoProducto("P003", "Teclado", 75.00, 30)

	// Simular compra: 2 laptops, 3 mouses, 1 teclado
	compras := map[*modelos.Producto]int{
		p1: 2,
		p2: 3,
		p3: 1,
	}

	for producto, cantidad := range compras {
		err := producto.DescontarStock(cantidad)
		if err != nil {
			t.Fatalf("Error al descontar stock de '%s': %v", producto.GetNombre(), err)
		}
	}

	// Verificar stocks finales
	if p1.GetCantidadDisponible() != 8 {
		t.Errorf("Stock de Laptop esperado 8, obtenido %d", p1.GetCantidadDisponible())
	}
	if p2.GetCantidadDisponible() != 47 {
		t.Errorf("Stock de Mouse esperado 47, obtenido %d", p2.GetCantidadDisponible())
	}
	if p3.GetCantidadDisponible() != 29 {
		t.Errorf("Stock de Teclado esperado 29, obtenido %d", p3.GetCantidadDisponible())
	}
}

// --- PRUEBAS DEL CARRITO ---

// TestCarrito_AgregarElemento verifica que se pueda agregar un producto al carrito.
func TestCarrito_AgregarElemento(t *testing.T) {
	carrito := modelos.NuevoCarrito()
	producto, _ := modelos.NuevoProducto("P001", "Laptop", 899.99, 10)

	err := carrito.AgregarElemento(*producto, 2)
	if err != nil {
		t.Fatalf("Error inesperado al agregar elemento al carrito: %v", err)
	}

	if carrito.EsVacio() {
		t.Fatal("El carrito no debe estar vacío después de agregar un elemento")
	}

	elementos := carrito.GetElementos()
	if len(elementos) != 1 {
		t.Errorf("Se esperaba 1 elemento en el carrito, obtenido %d", len(elementos))
	}

	// Total = 899.99 * 2 = 1799.98
	if carrito.GetTotal() != 1799.98 {
		t.Errorf("Total esperado 1799.98, obtenido %.2f", carrito.GetTotal())
	}
}

// TestCarrito_AgregarStockInsuficiente verifica que no se agregue si no hay stock.
func TestCarrito_AgregarStockInsuficiente(t *testing.T) {
	carrito := modelos.NuevoCarrito()
	producto, _ := modelos.NuevoProducto("P001", "Laptop", 899.99, 3)

	err := carrito.AgregarElemento(*producto, 10)
	if err == nil {
		t.Fatal("Se esperaba un error al agregar más cantidad de la que hay en stock")
	}

	if !carrito.EsVacio() {
		t.Fatal("El carrito debe seguir vacío si la operación falló")
	}
}

// TestCarrito_AcumularMismoProducto verifica que agregar el mismo producto acumule la cantidad.
func TestCarrito_AcumularMismoProducto(t *testing.T) {
	carrito := modelos.NuevoCarrito()
	producto, _ := modelos.NuevoProducto("P001", "Laptop", 100.00, 20)

	carrito.AgregarElemento(*producto, 3)
	carrito.AgregarElemento(*producto, 2)

	elementos := carrito.GetElementos()
	if len(elementos) != 1 {
		t.Errorf("Se esperaba 1 elemento acumulado, obtenido %d", len(elementos))
	}
	if elementos[0].GetCantidad() != 5 {
		t.Errorf("Cantidad acumulada esperada 5, obtenida %d", elementos[0].GetCantidad())
	}
	// Total = 100.00 * 5 = 500.00
	if carrito.GetTotal() != 500.00 {
		t.Errorf("Total esperado 500.00, obtenido %.2f", carrito.GetTotal())
	}
}

// TestCarrito_Vaciar verifica que el carrito se vacíe completamente.
func TestCarrito_Vaciar(t *testing.T) {
	carrito := modelos.NuevoCarrito()
	producto, _ := modelos.NuevoProducto("P001", "Laptop", 899.99, 10)

	carrito.AgregarElemento(*producto, 2)
	carrito.Vaciar()

	if !carrito.EsVacio() {
		t.Fatal("El carrito debe estar vacío después de vaciarlo")
	}
	if carrito.GetTotal() != 0 {
		t.Errorf("Total esperado 0 después de vaciar, obtenido %.2f", carrito.GetTotal())
	}
}

// --- PRUEBAS DEL CLIENTE ---

// TestCliente_CrearClienteValido verifica la creación exitosa de un cliente.
func TestCliente_CrearClienteValido(t *testing.T) {
	cliente, err := modelos.NuevoCliente("1234567890", "María López", "maria@email.com")

	if err != nil {
		t.Fatalf("Error inesperado al crear cliente: %v", err)
	}
	if cliente.GetIdentificacion() != "1234567890" {
		t.Errorf("Identificación esperada '1234567890', obtenida '%s'", cliente.GetIdentificacion())
	}
}

// TestCliente_CorreoInvalido verifica que se rechace un correo sin formato válido.
func TestCliente_CorreoInvalido(t *testing.T) {
	_, err := modelos.NuevoCliente("9876543210", "Pedro", "correo-sin-arroba")

	if err == nil {
		t.Fatal("Se esperaba un error al crear un cliente con correo inválido")
	}
}
