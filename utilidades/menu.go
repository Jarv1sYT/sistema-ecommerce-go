package utilidades

import "fmt"

// MostrarMenuPrincipal imprime en consola la interfaz gráfica del menú principal del sistema.
func MostrarMenuPrincipal() {
	fmt.Println("\n========================================")
	fmt.Println("    SISTEMA DE GESTIÓN DE E-COMMERCE    ")
	fmt.Println("========================================")
	fmt.Println("1. Gestión de productos")
	fmt.Println("2. Gestión de clientes")
	fmt.Println("3. Gestión de inventario")
	fmt.Println("4. Carrito de compras")
	fmt.Println("5. Gestión de pedidos")
	fmt.Println("6. Reportes")
	fmt.Println("0. Salir")
	fmt.Println("========================================")
}

// MostrarSubmenuProductos muestra el menú de opciones para la administración de productos.
func MostrarSubmenuProductos() {
	fmt.Println("\n========================================")
	fmt.Println("          GESTIÓN DE PRODUCTOS          ")
	fmt.Println("========================================")
	fmt.Println("1. Registrar producto")
	fmt.Println("2. Listar productos")
	fmt.Println("3. Buscar producto por código")
	fmt.Println("4. Modificar producto")
	fmt.Println("0. Volver al menú principal")
	fmt.Println("========================================")
}

// MostrarSubmenuClientes muestra el menú de opciones para la administración de clientes.
func MostrarSubmenuClientes() {
	fmt.Println("\n========================================")
	fmt.Println("          GESTIÓN DE CLIENTES           ")
	fmt.Println("========================================")
	fmt.Println("1. Registrar cliente")
	fmt.Println("2. Listar clientes")
	fmt.Println("3. Buscar cliente por identificación")
	fmt.Println("0. Volver al menú principal")
	fmt.Println("========================================")
}

// MostrarSubmenuInventario muestra el menú de opciones para el control de inventario.
func MostrarSubmenuInventario() {
	fmt.Println("\n========================================")
	fmt.Println("         GESTIÓN DE INVENTARIO          ")
	fmt.Println("========================================")
	fmt.Println("1. Mostrar inventario general")
	fmt.Println("2. Buscar producto en inventario")
	fmt.Println("3. Mostrar productos sin existencias")
	fmt.Println("4. Mostrar productos con existencias bajas (<= 5)")
	fmt.Println("0. Volver al menú principal")
	fmt.Println("========================================")
}

// MostrarSubmenuCarrito muestra las opciones disponibles dentro del carrito de compras.
func MostrarSubmenuCarrito() {
	fmt.Println("\n========================================")
	fmt.Println("           CARRITO DE COMPRAS           ")
	fmt.Println("========================================")
	fmt.Println("1. Seleccionar cliente")
	fmt.Println("2. Agregar producto al carrito")
	fmt.Println("3. Mostrar contenido del carrito")
	fmt.Println("4. Eliminar producto del carrito")
	fmt.Println("5. Vaciar carrito")
	fmt.Println("6. Confirmar compra")
	fmt.Println("0. Cancelar y volver")
	fmt.Println("========================================")
}

// MostrarSubmenuPedidos muestra las opciones para consultar órdenes registradas.
func MostrarSubmenuPedidos() {
	fmt.Println("\n========================================")
	fmt.Println("           GESTIÓN DE PEDIDOS           ")
	fmt.Println("========================================")
	fmt.Println("1. Listar todos los pedidos")
	fmt.Println("2. Buscar pedido por código")
	fmt.Println("3. Buscar pedidos por identificación del cliente")
	fmt.Println("0. Volver al menú principal")
	fmt.Println("========================================")
}

// MostrarSubmenuReportes muestra las opciones para consultar métricas del negocio.
func MostrarSubmenuReportes() {
	fmt.Println("\n========================================")
	fmt.Println("          REPORTES DEL SISTEMA          ")
	fmt.Println("========================================")
	fmt.Println("1. Resumen general")
	fmt.Println("2. Reporte de inventario")
	fmt.Println("3. Reporte de ventas acumuladas")
	fmt.Println("0. Volver al menú principal")
	fmt.Println("========================================")
}
