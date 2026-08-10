package servicios

import (
	"fmt"
	"strings"

	"sistema-ecommerce-go/modelos"
	"sistema-ecommerce-go/utilidades"
)

// Ruta donde se almacena el archivo de persistencia de productos
const rutaProductosJSON = "datos/productos.json"

// CargarProductos lee los productos guardados en el archivo JSON local.
func CargarProductos() ([]modelos.Producto, error) {
	var productos []modelos.Producto
	err := utilidades.LeerJSON(rutaProductosJSON, &productos)
	if err != nil {
		return nil, err
	}
	return productos, nil
}

// GuardarProductos serializa y almacena el listado de productos en el archivo JSON.
func GuardarProductos(productos []modelos.Producto) error {
	return utilidades.GuardarJSON(rutaProductosJSON, productos)
}

// BuscarProductoPorCodigo busca un producto dentro del slice utilizando el getter GetCodigo().
// Retorna la estructura encontrada, su posición (índice) en el slice y un booleano que indica si existe.
func BuscarProductoPorCodigo(productos []modelos.Producto, codigo string) (modelos.Producto, int, bool) {
	codigoLimpio := strings.ToLower(strings.TrimSpace(codigo))
	for i, p := range productos {
		if strings.ToLower(p.GetCodigo()) == codigoLimpio {
			return p, i, true
		}
	}
	return modelos.Producto{}, -1, false
}

// RegistrarProducto solicita los datos de un nuevo producto, los valida mediante el constructor modelos.NuevoProducto,
// y lo agrega al listado si cumple con todas las reglas.
func RegistrarProducto(productos *[]modelos.Producto) {
	fmt.Println("\n--- REGISTRAR NUEVO PRODUCTO ---")

	// 1. Validar que el código no esté vacío ni repetido
	var codigo string
	for {
		codigo = utilidades.LeerTexto("Ingrese el código del producto (ej: P001): ")
		_, _, existe := BuscarProductoPorCodigo(*productos, codigo)
		if existe {
			fmt.Println("Error: Ya existe un producto registrado con ese código. Intente con otro.")
			continue
		}
		break
	}

	// 2. Solicitar nombre
	nombre := utilidades.LeerTexto("Ingrese el nombre del producto: ")

	// 3. Solicitar precio
	precio := utilidades.LeerDecimal("Ingrese el precio del producto: $")

	// 4. Solicitar cantidad inicial disponible
	cantidad := utilidades.LeerEntero("Ingrese la cantidad disponible en inventario: ")

	// Crear el producto utilizando el CONSTRUCTOR del modelo (aplica validaciones internas)
	nuevoProducto, err := modelos.NuevoProducto(codigo, nombre, precio, cantidad)
	if err != nil {
		fmt.Printf("\nError de validación al crear el producto: %v\n", err)
		return
	}

	// Agregar el producto al slice dereferenciado por puntero
	*productos = append(*productos, *nuevoProducto)

	// Persistir inmediatamente en JSON
	err = GuardarProductos(*productos)
	if err != nil {
		fmt.Printf("Error al guardar el producto en archivo JSON: %v\n", err)
		return
	}

	fmt.Println("\n¡Producto registrado y guardado exitosamente en datos/productos.json!")
}

// ListarProductos muestra en consola todos los productos utilizando los Getters de la estructura.
func ListarProductos(productos []modelos.Producto) {
	fmt.Println("\n--- LISTADO DE PRODUCTOS ---")

	if len(productos) == 0 {
		fmt.Println("No existen productos registrados en el sistema.")
		return
	}

	fmt.Printf("%-10s | %-25s | %-12s | %-10s\n", "CÓDIGO", "NOMBRE", "PRECIO", "STOCK")
	fmt.Println(strings.Repeat("-", 65))

	for _, p := range productos {
		fmt.Printf("%-10s | %-25s | $%11.2f | %-10d\n", p.GetCodigo(), p.GetNombre(), p.GetPrecio(), p.GetCantidadDisponible())
	}
	fmt.Println(strings.Repeat("-", 65))
	fmt.Printf("Total de productos registrados: %d\n", len(productos))
}

// BuscarProductoMenu solicita un código al usuario y muestra los detalles mediante Getters.
func BuscarProductoMenu(productos []modelos.Producto) {
	fmt.Println("\n--- BUSCAR PRODUCTO POR CÓDIGO ---")

	if len(productos) == 0 {
		fmt.Println("No existen productos registrados en el sistema.")
		return
	}

	codigo := utilidades.LeerTexto("Ingrese el código del producto a buscar: ")
	producto, _, encontrado := BuscarProductoPorCodigo(productos, codigo)

	if !encontrado {
		fmt.Printf("No se encontró ningún producto registrado con el código '%s'.\n", codigo)
		return
	}

	fmt.Println("\n--- DETALLES DEL PRODUCTO ENCONTRADO ---")
	fmt.Printf("Código:              %s\n", producto.GetCodigo())
	fmt.Printf("Nombre:              %s\n", producto.GetNombre())
	fmt.Printf("Precio:              $%.2f\n", producto.GetPrecio())
	fmt.Printf("Cantidad Disponible: %d unidades\n", producto.GetCantidadDisponible())
}

// ModificarProducto permite cambiar atributos de un producto utilizando exclusivamente sus Setters controlados.
func ModificarProducto(productos *[]modelos.Producto) {
	fmt.Println("\n--- MODIFICAR PRODUCTO ---")

	if len(*productos) == 0 {
		fmt.Println("No existen productos registrados para modificar.")
		return
	}

	codigo := utilidades.LeerTexto("Ingrese el código del producto que desea modificar: ")
	producto, pos, encontrado := BuscarProductoPorCodigo(*productos, codigo)

	if !encontrado {
		fmt.Printf("No se encontró ningún producto con el código '%s'.\n", codigo)
		return
	}

	fmt.Printf("\nProducto seleccionado: %s - %s ($%.2f, Stock: %d)\n",
		producto.GetCodigo(), producto.GetNombre(), producto.GetPrecio(), producto.GetCantidadDisponible())
	fmt.Println("1. Modificar nombre")
	fmt.Println("2. Modificar precio")
	fmt.Println("3. Modificar cantidad disponible")
	fmt.Println("0. Cancelar")

	opcion := utilidades.LeerOpcion("Seleccione qué campo desea modificar: ")

	switch opcion {
	case 1:
		nuevoNombre := utilidades.LeerTexto("Ingrese el nuevo nombre: ")
		err := (*productos)[pos].SetNombre(nuevoNombre)
		if err != nil {
			fmt.Printf("Error al modificar el nombre: %v\n", err)
			return
		}
		fmt.Println("Nombre actualizado correctamente.")
	case 2:
		nuevoPrecio := utilidades.LeerDecimal("Ingrese el nuevo precio: $")
		err := (*productos)[pos].SetPrecio(nuevoPrecio)
		if err != nil {
			fmt.Printf("Error al modificar el precio: %v\n", err)
			return
		}
		fmt.Println("Precio actualizado correctamente.")
	case 3:
		nuevaCantidad := utilidades.LeerEntero("Ingrese la nueva cantidad disponible: ")
		err := (*productos)[pos].SetCantidadDisponible(nuevaCantidad)
		if err != nil {
			fmt.Printf("Error al modificar la cantidad disponible: %v\n", err)
			return
		}
		fmt.Println("Cantidad disponible actualizada correctamente.")
	case 0:
		fmt.Println("Operación cancelada. No se realizaron cambios.")
		return
	default:
		fmt.Println("Opción inválida. No se realizaron cambios.")
		return
	}

	// Guardar los cambios en el archivo JSON
	err := GuardarProductos(*productos)
	if err != nil {
		fmt.Printf("Error al guardar los cambios en JSON: %v\n", err)
		return
	}

	fmt.Println("¡Producto modificado y guardado con éxito en datos/productos.json!")
}

// EjecutarModuloProductos controla la navegación y ejecución del submenú de productos.
func EjecutarModuloProductos(productos *[]modelos.Producto) {
	for {
		utilidades.MostrarSubmenuProductos()
		opcion := utilidades.LeerOpcion("Seleccione una opción: ")

		switch opcion {
		case 1:
			RegistrarProducto(productos)
			utilidades.Pausar()
		case 2:
			ListarProductos(*productos)
			utilidades.Pausar()
		case 3:
			BuscarProductoMenu(*productos)
			utilidades.Pausar()
		case 4:
			ModificarProducto(productos)
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
