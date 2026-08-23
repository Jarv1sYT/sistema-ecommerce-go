package servicios

import (
	"testing"

	"sistema-ecommerce-go/modelos"
)

// ============================================================================
// PRUEBAS UNITARIAS PARA EL SERVICIO DE PRODUCTOS (Unidad 4 - Testing)
// ============================================================================
// Ejecutar: go test -v ./servicios/ -run TestProductos
// ============================================================================

// --- PRUEBAS DEL CONSTRUCTOR NuevoProducto (Validaciones de Encapsulación) ---

// TestProducto_CrearProductoValido verifica que se pueda crear un producto con datos válidos.
func TestProducto_CrearProductoValido(t *testing.T) {
	producto, err := modelos.NuevoProducto("P001", "Laptop HP", 899.99, 10)

	if err != nil {
		t.Fatalf("Error inesperado al crear producto válido: %v", err)
	}
	if producto == nil {
		t.Fatal("El producto creado no debe ser nil")
	}
	if producto.GetCodigo() != "P001" {
		t.Errorf("Código esperado 'P001', obtenido '%s'", producto.GetCodigo())
	}
	if producto.GetNombre() != "Laptop HP" {
		t.Errorf("Nombre esperado 'Laptop HP', obtenido '%s'", producto.GetNombre())
	}
	if producto.GetPrecio() != 899.99 {
		t.Errorf("Precio esperado 899.99, obtenido %.2f", producto.GetPrecio())
	}
	if producto.GetCantidadDisponible() != 10 {
		t.Errorf("Cantidad esperada 10, obtenida %d", producto.GetCantidadDisponible())
	}
}

// TestProducto_PrecioNegativo verifica que NO se permita crear un producto con precio negativo.
func TestProducto_PrecioNegativo(t *testing.T) {
	producto, err := modelos.NuevoProducto("P002", "Producto inválido", -5.00, 10)

	if err == nil {
		t.Fatal("Se esperaba un error al crear un producto con precio negativo, pero no se recibió ninguno")
	}
	if producto != nil {
		t.Fatal("El producto debe ser nil cuando hay error de validación")
	}
}

// TestProducto_PrecioCero verifica que NO se permita precio igual a cero.
func TestProducto_PrecioCero(t *testing.T) {
	_, err := modelos.NuevoProducto("P003", "Producto gratis", 0, 5)

	if err == nil {
		t.Fatal("Se esperaba un error al crear un producto con precio cero")
	}
}

// TestProducto_StockNegativo verifica que NO se permita stock negativo.
func TestProducto_StockNegativo(t *testing.T) {
	_, err := modelos.NuevoProducto("P004", "Producto stock negativo", 10.00, -3)

	if err == nil {
		t.Fatal("Se esperaba un error al crear un producto con stock negativo")
	}
}

// TestProducto_CodigoVacio verifica que NO se permita un código vacío.
func TestProducto_CodigoVacio(t *testing.T) {
	_, err := modelos.NuevoProducto("", "Producto sin código", 10.00, 5)

	if err == nil {
		t.Fatal("Se esperaba un error al crear un producto con código vacío")
	}
}

// TestProducto_NombreVacio verifica que NO se permita un nombre vacío.
func TestProducto_NombreVacio(t *testing.T) {
	_, err := modelos.NuevoProducto("P005", "", 10.00, 5)

	if err == nil {
		t.Fatal("Se esperaba un error al crear un producto con nombre vacío")
	}
}

// --- PRUEBAS DE DESCUENTO DE STOCK ---

// TestProducto_DescontarStockValido verifica el descuento correcto de stock.
func TestProducto_DescontarStockValido(t *testing.T) {
	producto, _ := modelos.NuevoProducto("P010", "Mouse Logitech", 25.50, 20)

	err := producto.DescontarStock(5)
	if err != nil {
		t.Fatalf("Error inesperado al descontar stock válido: %v", err)
	}
	if producto.GetCantidadDisponible() != 15 {
		t.Errorf("Stock esperado 15 después de descontar 5 de 20, obtenido %d", producto.GetCantidadDisponible())
	}
}

// TestProducto_DescontarStockInsuficiente verifica que falle si no hay stock suficiente.
func TestProducto_DescontarStockInsuficiente(t *testing.T) {
	producto, _ := modelos.NuevoProducto("P011", "Teclado mecánico", 75.00, 3)

	err := producto.DescontarStock(10)
	if err == nil {
		t.Fatal("Se esperaba un error al descontar más stock del disponible")
	}
	// El stock NO debe haber cambiado
	if producto.GetCantidadDisponible() != 3 {
		t.Errorf("El stock no debe cambiar cuando el descuento falla. Esperado 3, obtenido %d", producto.GetCantidadDisponible())
	}
}

// TestProducto_DescontarStockExacto verifica que se pueda descontar exactamente todo el stock.
func TestProducto_DescontarStockExacto(t *testing.T) {
	producto, _ := modelos.NuevoProducto("P012", "Cable HDMI", 12.00, 5)

	err := producto.DescontarStock(5)
	if err != nil {
		t.Fatalf("Debería permitir descontar la cantidad exacta del stock: %v", err)
	}
	if producto.GetCantidadDisponible() != 0 {
		t.Errorf("Stock esperado 0 después de descontar todo, obtenido %d", producto.GetCantidadDisponible())
	}
}

// --- PRUEBAS DE BÚSQUEDA ---

// TestBuscarProductoPorCodigo_Existente verifica la búsqueda exitosa de un producto.
func TestBuscarProductoPorCodigo_Existente(t *testing.T) {
	p1, _ := modelos.NuevoProducto("P001", "Laptop", 900.00, 5)
	p2, _ := modelos.NuevoProducto("P002", "Mouse", 25.00, 20)
	p3, _ := modelos.NuevoProducto("P003", "Teclado", 50.00, 15)

	productos := []modelos.Producto{*p1, *p2, *p3}

	producto, indice, encontrado := BuscarProductoPorCodigo(productos, "P002")

	if !encontrado {
		t.Fatal("Se esperaba encontrar el producto P002")
	}
	if indice != 1 {
		t.Errorf("Índice esperado 1, obtenido %d", indice)
	}
	if producto.GetNombre() != "Mouse" {
		t.Errorf("Nombre esperado 'Mouse', obtenido '%s'", producto.GetNombre())
	}
}

// TestBuscarProductoPorCodigo_Inexistente verifica que no encuentre un producto inexistente.
func TestBuscarProductoPorCodigo_Inexistente(t *testing.T) {
	p1, _ := modelos.NuevoProducto("P001", "Laptop", 900.00, 5)
	productos := []modelos.Producto{*p1}

	_, indice, encontrado := BuscarProductoPorCodigo(productos, "P999")

	if encontrado {
		t.Fatal("No se esperaba encontrar un producto inexistente")
	}
	if indice != -1 {
		t.Errorf("El índice debe ser -1 cuando no se encuentra, obtenido %d", indice)
	}
}

// TestBuscarProductoPorCodigo_InsensibleMayusculas verifica búsqueda sin distinción de mayúsculas.
func TestBuscarProductoPorCodigo_InsensibleMayusculas(t *testing.T) {
	p1, _ := modelos.NuevoProducto("P001", "Laptop", 900.00, 5)
	productos := []modelos.Producto{*p1}

	_, _, encontrado := BuscarProductoPorCodigo(productos, "p001")

	if !encontrado {
		t.Fatal("La búsqueda debería ser insensible a mayúsculas/minúsculas")
	}
}

// --- PRUEBAS DE SETTERS CONTROLADOS ---

// TestProducto_SetPrecioValido verifica que el setter de precio acepte valores válidos.
func TestProducto_SetPrecioValido(t *testing.T) {
	producto, _ := modelos.NuevoProducto("P020", "Monitor", 300.00, 8)

	err := producto.SetPrecio(350.00)
	if err != nil {
		t.Fatalf("Error inesperado al cambiar precio: %v", err)
	}
	if producto.GetPrecio() != 350.00 {
		t.Errorf("Precio esperado 350.00, obtenido %.2f", producto.GetPrecio())
	}
}

// TestProducto_SetPrecioInvalido verifica que el setter rechace precios inválidos.
func TestProducto_SetPrecioInvalido(t *testing.T) {
	producto, _ := modelos.NuevoProducto("P021", "Auriculares", 80.00, 12)

	err := producto.SetPrecio(-10.00)
	if err == nil {
		t.Fatal("Se esperaba un error al establecer precio negativo")
	}
	// El precio NO debe haber cambiado
	if producto.GetPrecio() != 80.00 {
		t.Errorf("El precio no debe cambiar cuando la validación falla. Esperado 80.00, obtenido %.2f", producto.GetPrecio())
	}
}
