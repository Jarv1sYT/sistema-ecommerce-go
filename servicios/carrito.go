package servicios

import (
	"fmt"
	"strings"
	"time"

	"sistema-ecommerce-go/modelos"
	"sistema-ecommerce-go/utilidades"
)

// Rutas de persistencia JSON
const rutaPedidosJSON = "datos/pedidos.json"

// ----------------------------------------------------------------------------
// FUNCIONES AUXILIARES DE CARGA Y GUARDADO DE PEDIDOS
// ----------------------------------------------------------------------------

// CargarPedidos lee los pedidos guardados en el archivo JSON local.
func CargarPedidos() ([]modelos.Pedido, error) {
	var pedidos []modelos.Pedido
	err := utilidades.LeerJSON(rutaPedidosJSON, &pedidos)
	if err != nil {
		return nil, err
	}
	return pedidos, nil
}

// GuardarPedidos serializa y almacena el listado de pedidos en el archivo JSON.
func GuardarPedidos(pedidos []modelos.Pedido) error {
	return utilidades.GuardarJSON(rutaPedidosJSON, pedidos)
}

// generarCodigoPedido crea un código único de pedido basado en la cantidad actual de pedidos.
// Formato: PED-0001, PED-0002, etc.
func generarCodigoPedido(pedidosExistentes int) string {
	return fmt.Sprintf("PED-%04d", pedidosExistentes+1)
}

// ----------------------------------------------------------------------------
// FUNCIONES DE INTERACCIÓN CON EL CARRITO (CONSOLA)
// ----------------------------------------------------------------------------

// SeleccionarCliente permite al usuario asociar un cliente registrado al carrito de compras.
// Busca al cliente por identificación usando el servicio de clientes existente.
func SeleccionarCliente(carrito *modelos.Carrito, clientes []modelos.Cliente) {
	fmt.Println("\n--- SELECCIONAR CLIENTE ---")

	if len(clientes) == 0 {
		fmt.Println("No hay clientes registrados. Registre un cliente primero desde el menú principal.")
		return
	}

	id := utilidades.LeerTexto("Ingrese la identificación del cliente: ")
	cliente, _, encontrado := BuscarClientePorIdentificacion(clientes, id)

	if !encontrado {
		fmt.Printf("No se encontró ningún cliente con la identificación '%s'.\n", id)
		return
	}

	// Asociar el cliente al carrito usando el setter del modelo encapsulado
	carrito.SetCliente(&cliente)
	fmt.Printf("Cliente seleccionado: %s (%s)\n", cliente.GetNombre(), cliente.GetIdentificacion())
}

// AgregarProductoAlCarrito solicita un código de producto y la cantidad deseada,
// valida que exista en el inventario y que haya stock suficiente, y lo agrega al carrito.
func AgregarProductoAlCarrito(carrito *modelos.Carrito, productos []modelos.Producto) {
	fmt.Println("\n--- AGREGAR PRODUCTO AL CARRITO ---")

	if len(productos) == 0 {
		fmt.Println("No hay productos registrados en el inventario.")
		return
	}

	// Mostrar catálogo resumido para facilitar la selección
	fmt.Println("\nProductos disponibles:")
	fmt.Printf("%-10s | %-25s | %-12s | %-10s\n", "CÓDIGO", "NOMBRE", "PRECIO", "STOCK")
	fmt.Println(strings.Repeat("-", 65))
	for _, p := range productos {
		fmt.Printf("%-10s | %-25s | $%11.2f | %-10d\n",
			p.GetCodigo(), p.GetNombre(), p.GetPrecio(), p.GetCantidadDisponible())
	}
	fmt.Println(strings.Repeat("-", 65))

	// Solicitar código del producto
	codigo := utilidades.LeerTexto("Ingrese el código del producto a agregar: ")
	producto, _, encontrado := BuscarProductoPorCodigo(productos, codigo)

	if !encontrado {
		fmt.Printf("No se encontró ningún producto con el código '%s'.\n", codigo)
		return
	}

	// Verificar que el producto tiene stock antes de pedir cantidad
	if producto.GetCantidadDisponible() <= 0 {
		fmt.Printf("El producto '%s' no tiene stock disponible en este momento.\n", producto.GetNombre())
		return
	}

	// Solicitar la cantidad deseada
	fmt.Printf("Stock disponible de '%s': %d unidades\n", producto.GetNombre(), producto.GetCantidadDisponible())
	cantidad := utilidades.LeerEntero("Ingrese la cantidad que desea agregar: ")

	if cantidad <= 0 {
		fmt.Println("La cantidad debe ser mayor que cero.")
		return
	}

	// Delegar la lógica de agregar al método encapsulado del Carrito
	// (AgregarElemento ya valida stock y acumula si el producto ya estaba en el carrito)
	err := carrito.AgregarElemento(producto, cantidad)
	if err != nil {
		fmt.Printf("No se pudo agregar al carrito: %v\n", err)
		return
	}

	fmt.Printf("\n¡Producto '%s' (x%d) agregado al carrito exitosamente!\n", producto.GetNombre(), cantidad)
	fmt.Printf("Total parcial del carrito: $%.2f\n", carrito.GetTotal())
}

// MostrarCarrito imprime el contenido actual del carrito de compras con formato tabular.
func MostrarCarrito(carrito *modelos.Carrito) {
	fmt.Println("\n--- CONTENIDO DEL CARRITO ---")

	if carrito.EsVacio() {
		fmt.Println("El carrito está vacío. Agregue productos primero.")
		return
	}

	// Mostrar cliente asociado (si existe)
	cliente := carrito.GetCliente()
	if cliente != nil {
		fmt.Printf("Cliente: %s (%s)\n", cliente.GetNombre(), cliente.GetIdentificacion())
	} else {
		fmt.Println("Cliente: (No seleccionado)")
	}

	fmt.Println()
	fmt.Printf("%-10s | %-25s | %-10s | %-8s | %-12s\n", "CÓDIGO", "PRODUCTO", "PRECIO", "CANT.", "SUBTOTAL")
	fmt.Println(strings.Repeat("-", 75))

	// Obtener la COPIA SEGURA de los elementos (respetando la encapsulación)
	elementos := carrito.GetElementos()
	for _, elem := range elementos {
		prod := elem.GetProducto()
		fmt.Printf("%-10s | %-25s | $%9.2f | %8d | $%11.2f\n",
			prod.GetCodigo(),
			prod.GetNombre(),
			prod.GetPrecio(),
			elem.GetCantidad(),
			elem.GetSubtotal())
	}

	fmt.Println(strings.Repeat("-", 75))
	fmt.Printf("TOTAL A PAGAR: $%.2f\n", carrito.GetTotal())
}

// EliminarProductoDelCarrito permite al usuario eliminar un ítem del carrito por su código.
func EliminarProductoDelCarrito(carrito *modelos.Carrito) {
	fmt.Println("\n--- ELIMINAR PRODUCTO DEL CARRITO ---")

	if carrito.EsVacio() {
		fmt.Println("El carrito está vacío. No hay productos que eliminar.")
		return
	}

	// Mostrar contenido actual para referencia
	MostrarCarrito(carrito)

	codigo := utilidades.LeerTexto("Ingrese el código del producto a eliminar del carrito: ")

	eliminado := carrito.EliminarElemento(codigo)
	if !eliminado {
		fmt.Printf("No se encontró el producto con código '%s' en el carrito.\n", codigo)
		return
	}

	fmt.Println("Producto eliminado del carrito exitosamente.")
	fmt.Printf("Total actualizado del carrito: $%.2f\n", carrito.GetTotal())
}

// VaciarCarrito permite al usuario vaciar completamente el carrito previa confirmación.
func VaciarCarrito(carrito *modelos.Carrito) {
	fmt.Println("\n--- VACIAR CARRITO ---")

	if carrito.EsVacio() {
		fmt.Println("El carrito ya está vacío.")
		return
	}

	confirmacion := utilidades.Confirmar("¿Está seguro de que desea vaciar todo el carrito?")
	if !confirmacion {
		fmt.Println("Operación cancelada. El carrito no fue modificado.")
		return
	}

	carrito.Vaciar()
	fmt.Println("El carrito ha sido vaciado completamente.")
}

// ----------------------------------------------------------------------------
// PROCESO DE CHECKOUT (CONFIRMAR COMPRA)
// ----------------------------------------------------------------------------

// ConfirmarCompra ejecuta la transacción completa de compra:
// 1. Valida que el carrito tenga productos y un cliente asociado.
// 2. Verifica que cada producto del carrito aún tenga stock suficiente en el inventario actual.
// 3. Descuenta el stock de cada producto en el slice de productos.
// 4. Crea un Pedido inmutable usando el constructor NuevoPedido.
// 5. Persiste los pedidos actualizados y los productos con stock reducido en JSON.
// 6. Vacía el carrito.
func ConfirmarCompra(carrito *modelos.Carrito, productos *[]modelos.Producto) {
	fmt.Println("\n--- CONFIRMAR COMPRA (CHECKOUT) ---")

	// Validación 1: El carrito no puede estar vacío
	if carrito.EsVacio() {
		fmt.Println("Error: El carrito está vacío. Agregue productos antes de confirmar la compra.")
		return
	}

	// Validación 2: Debe existir un cliente asociado
	cliente := carrito.GetCliente()
	if cliente == nil {
		fmt.Println("Error: No se ha seleccionado un cliente. Seleccione un cliente antes de confirmar la compra.")
		return
	}

	// Mostrar resumen de la compra antes de confirmar
	fmt.Println("\n========================================")
	fmt.Println("      RESUMEN DE LA COMPRA")
	fmt.Println("========================================")
	fmt.Printf("Cliente: %s (%s)\n", cliente.GetNombre(), cliente.GetIdentificacion())
	fmt.Println()

	elementos := carrito.GetElementos()
	fmt.Printf("%-25s | %-8s | %-12s\n", "PRODUCTO", "CANT.", "SUBTOTAL")
	fmt.Println(strings.Repeat("-", 50))
	for _, elem := range elementos {
		fmt.Printf("%-25s | %8d | $%11.2f\n",
			elem.GetProducto().GetNombre(),
			elem.GetCantidad(),
			elem.GetSubtotal())
	}
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("TOTAL: $%.2f\n", carrito.GetTotal())
	fmt.Println("========================================")

	// Solicitar confirmación final del usuario
	confirmacion := utilidades.Confirmar("¿Desea confirmar esta compra?")
	if !confirmacion {
		fmt.Println("Compra cancelada. El carrito no fue modificado.")
		return
	}

	// Validación 3: Verificar stock suficiente para CADA producto del carrito
	// (Podría haber cambiado desde que el usuario agregó los ítems)
	for _, elem := range elementos {
		codigoProducto := elem.GetProducto().GetCodigo()
		cantidadSolicitada := elem.GetCantidad()

		_, posInventario, existe := BuscarProductoPorCodigo(*productos, codigoProducto)
		if !existe {
			fmt.Printf("Error: El producto '%s' ya no existe en el inventario.\n", codigoProducto)
			return
		}

		productoReal := (*productos)[posInventario]
		if !productoReal.TieneStockSuficiente(cantidadSolicitada) {
			fmt.Printf("Error: Stock insuficiente para '%s'. Disponible: %d, Solicitado: %d\n",
				productoReal.GetNombre(), productoReal.GetCantidadDisponible(), cantidadSolicitada)
			return
		}
	}

	// PASO TRANSACCIONAL: Descontar stock de cada producto en el inventario
	for _, elem := range elementos {
		codigoProducto := elem.GetProducto().GetCodigo()
		cantidadSolicitada := elem.GetCantidad()

		_, posInventario, _ := BuscarProductoPorCodigo(*productos, codigoProducto)

		err := (*productos)[posInventario].DescontarStock(cantidadSolicitada)
		if err != nil {
			fmt.Printf("Error crítico al descontar stock de '%s': %v\n", codigoProducto, err)
			return
		}
	}

	// Guardar los productos con stock actualizado
	err := GuardarProductos(*productos)
	if err != nil {
		fmt.Printf("Error al guardar los productos actualizados: %v\n", err)
		return
	}

	// Cargar pedidos existentes para generar el siguiente código
	pedidos, err := CargarPedidos()
	if err != nil {
		fmt.Printf("Aviso: No se pudieron cargar los pedidos anteriores: %v\n", err)
		pedidos = []modelos.Pedido{}
	}

	// Generar código único y fecha actual
	codigoPedido := generarCodigoPedido(len(pedidos))
	fechaActual := time.Now().Format("2006-01-02 15:04:05")

	// Crear el Pedido INMUTABLE usando el constructor del modelo
	nuevoPedido, err := modelos.NuevoPedido(codigoPedido, *cliente, elementos, fechaActual)
	if err != nil {
		fmt.Printf("Error al crear el pedido: %v\n", err)
		return
	}

	// Agregar el nuevo pedido al historial y persistir
	pedidos = append(pedidos, *nuevoPedido)
	err = GuardarPedidos(pedidos)
	if err != nil {
		fmt.Printf("Error al guardar el pedido en JSON: %v\n", err)
		return
	}

	// Vaciar el carrito después de la compra exitosa
	carrito.Vaciar()

	// Mostrar comprobante final
	fmt.Println("\n========================================")
	fmt.Println("     ¡COMPRA REALIZADA CON ÉXITO!")
	fmt.Println("========================================")
	fmt.Printf("Código de pedido: %s\n", nuevoPedido.GetCodigo())
	fmt.Printf("Cliente:          %s\n", nuevoPedido.GetCliente().GetNombre())
	fmt.Printf("Total cobrado:    $%.2f\n", nuevoPedido.GetTotal())
	fmt.Printf("Fecha:            %s\n", nuevoPedido.GetFecha())
	fmt.Println("========================================")
	fmt.Println("El inventario ha sido actualizado automáticamente.")
	fmt.Println("El pedido ha sido guardado en datos/pedidos.json.")
}

// ----------------------------------------------------------------------------
// CONTROLADOR DEL SUBMENÚ DE CARRITO
// ----------------------------------------------------------------------------

// EjecutarModuloCarrito controla la navegación del submenú del carrito de compras.
// Recibe el carrito (en memoria), el listado de productos y el listado de clientes.
func EjecutarModuloCarrito(carrito *modelos.Carrito, productos *[]modelos.Producto, clientes []modelos.Cliente) {
	for {
		utilidades.MostrarSubmenuCarrito()
		opcion := utilidades.LeerOpcion("Seleccione una opción: ")

		switch opcion {
		case 1:
			SeleccionarCliente(carrito, clientes)
			utilidades.Pausar()
		case 2:
			AgregarProductoAlCarrito(carrito, *productos)
			utilidades.Pausar()
		case 3:
			MostrarCarrito(carrito)
			utilidades.Pausar()
		case 4:
			EliminarProductoDelCarrito(carrito)
			utilidades.Pausar()
		case 5:
			VaciarCarrito(carrito)
			utilidades.Pausar()
		case 6:
			ConfirmarCompra(carrito, productos)
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
