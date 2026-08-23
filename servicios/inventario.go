package servicios

import (
	"fmt"
	"strings"

	"sistema-ecommerce-go/modelos"
	"sistema-ecommerce-go/utilidades"
)

// ----------------------------------------------------------------------------
// FUNCIONES DE CONSULTA DE INVENTARIO (CONSOLA)
// ----------------------------------------------------------------------------

// MostrarInventarioGeneral muestra todos los productos con formato de inventario,
// incluyendo el valor total del inventario (precio × cantidad por producto).
func MostrarInventarioGeneral(productos []modelos.Producto) {
	fmt.Println("\n--- INVENTARIO GENERAL ---")

	if len(productos) == 0 {
		fmt.Println("No existen productos registrados en el inventario.")
		return
	}

	fmt.Printf("%-10s | %-25s | %-12s | %-8s | %-14s\n", "CÓDIGO", "NOMBRE", "PRECIO", "STOCK", "VALOR TOTAL")
	fmt.Println(strings.Repeat("-", 80))

	valorTotalInventario := 0.0
	for _, p := range productos {
		valorProducto := p.GetPrecio() * float64(p.GetCantidadDisponible())
		valorTotalInventario += valorProducto
		fmt.Printf("%-10s | %-25s | $%11.2f | %8d | $%13.2f\n",
			p.GetCodigo(), p.GetNombre(), p.GetPrecio(), p.GetCantidadDisponible(), valorProducto)
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("Total de productos: %d | Valor total del inventario: $%.2f\n", len(productos), valorTotalInventario)
}

// MostrarProductosSinStock filtra y muestra los productos con existencias en cero.
func MostrarProductosSinStock(productos []modelos.Producto) {
	fmt.Println("\n--- PRODUCTOS SIN EXISTENCIAS (STOCK = 0) ---")

	sinStock := []modelos.Producto{}
	for _, p := range productos {
		if p.GetCantidadDisponible() == 0 {
			sinStock = append(sinStock, p)
		}
	}

	if len(sinStock) == 0 {
		fmt.Println("Todos los productos tienen stock disponible. ¡No hay productos agotados!")
		return
	}

	fmt.Printf("%-10s | %-25s | %-12s\n", "CÓDIGO", "NOMBRE", "PRECIO")
	fmt.Println(strings.Repeat("-", 55))
	for _, p := range sinStock {
		fmt.Printf("%-10s | %-25s | $%11.2f\n", p.GetCodigo(), p.GetNombre(), p.GetPrecio())
	}
	fmt.Println(strings.Repeat("-", 55))
	fmt.Printf("Productos sin stock: %d\n", len(sinStock))
}

// MostrarProductosBajoStock filtra y muestra los productos con existencias menores o iguales a 5 unidades.
func MostrarProductosBajoStock(productos []modelos.Producto) {
	fmt.Println("\n--- PRODUCTOS CON EXISTENCIAS BAJAS (<= 5 UNIDADES) ---")

	bajoStock := []modelos.Producto{}
	for _, p := range productos {
		if p.GetCantidadDisponible() > 0 && p.GetCantidadDisponible() <= 5 {
			bajoStock = append(bajoStock, p)
		}
	}

	if len(bajoStock) == 0 {
		fmt.Println("No hay productos con existencias bajas en este momento.")
		return
	}

	fmt.Printf("%-10s | %-25s | %-12s | %-8s\n", "CÓDIGO", "NOMBRE", "PRECIO", "STOCK")
	fmt.Println(strings.Repeat("-", 65))
	for _, p := range bajoStock {
		fmt.Printf("%-10s | %-25s | $%11.2f | %8d\n",
			p.GetCodigo(), p.GetNombre(), p.GetPrecio(), p.GetCantidadDisponible())
	}
	fmt.Println(strings.Repeat("-", 65))
	fmt.Printf("Productos con stock bajo: %d\n", len(bajoStock))
}

// ----------------------------------------------------------------------------
// CONTROLADOR DEL SUBMENÚ DE INVENTARIO
// ----------------------------------------------------------------------------

// EjecutarModuloInventario controla la navegación del submenú de inventario.
func EjecutarModuloInventario(productos []modelos.Producto) {
	for {
		utilidades.MostrarSubmenuInventario()
		opcion := utilidades.LeerOpcion("Seleccione una opción: ")

		switch opcion {
		case 1:
			MostrarInventarioGeneral(productos)
			utilidades.Pausar()
		case 2:
			BuscarProductoMenu(productos)
			utilidades.Pausar()
		case 3:
			MostrarProductosSinStock(productos)
			utilidades.Pausar()
		case 4:
			MostrarProductosBajoStock(productos)
			utilidades.Pausar()
		case 0:
			fmt.Println("Regresando al menú principal...")
			return
		default:
			fmt.Println("Opción no válida. Intente nuevamente.")
			utilidades.Pausar()
		}
	}
}
