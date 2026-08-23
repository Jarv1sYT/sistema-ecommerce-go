package servicios

import (
	"fmt"
	"strings"

	"sistema-ecommerce-go/modelos"
	"sistema-ecommerce-go/utilidades"
)

// ----------------------------------------------------------------------------
// FUNCIONES DE CONSULTA DE PEDIDOS (CONSOLA)
// ----------------------------------------------------------------------------

// ListarPedidos muestra en consola todos los pedidos registrados en el sistema.
func ListarPedidos(pedidos []modelos.Pedido) {
	fmt.Println("\n--- HISTORIAL DE PEDIDOS ---")

	if len(pedidos) == 0 {
		fmt.Println("No existen pedidos registrados en el sistema.")
		return
	}

	fmt.Printf("%-12s | %-20s | %-12s | %-20s\n", "CÓDIGO", "CLIENTE", "TOTAL", "FECHA")
	fmt.Println(strings.Repeat("-", 70))

	for _, p := range pedidos {
		fmt.Printf("%-12s | %-20s | $%11.2f | %-20s\n",
			p.GetCodigo(),
			p.GetCliente().GetNombre(),
			p.GetTotal(),
			p.GetFecha())
	}

	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("Total de pedidos registrados: %d\n", len(pedidos))
}

// BuscarPedidoPorCodigo busca un pedido por su código dentro del slice.
// Retorna el pedido encontrado, su posición y un booleano que indica si existe.
func BuscarPedidoPorCodigo(pedidos []modelos.Pedido, codigo string) (modelos.Pedido, int, bool) {
	codigoLimpio := strings.ToUpper(strings.TrimSpace(codigo))
	for i, p := range pedidos {
		if p.GetCodigo() == codigoLimpio {
			return p, i, true
		}
	}
	return modelos.Pedido{}, -1, false
}

// BuscarPedidoMenu solicita el código de un pedido y muestra sus detalles completos.
func BuscarPedidoMenu(pedidos []modelos.Pedido) {
	fmt.Println("\n--- BUSCAR PEDIDO POR CÓDIGO ---")

	if len(pedidos) == 0 {
		fmt.Println("No existen pedidos registrados en el sistema.")
		return
	}

	codigo := utilidades.LeerTexto("Ingrese el código del pedido a buscar (ej: PED-0001): ")
	pedido, _, encontrado := BuscarPedidoPorCodigo(pedidos, codigo)

	if !encontrado {
		fmt.Printf("No se encontró ningún pedido con el código '%s'.\n", codigo)
		return
	}

	mostrarDetallePedido(pedido)
}

// BuscarPedidosPorCliente busca todos los pedidos asociados a una identificación de cliente.
func BuscarPedidosPorCliente(pedidos []modelos.Pedido) {
	fmt.Println("\n--- BUSCAR PEDIDOS POR CLIENTE ---")

	if len(pedidos) == 0 {
		fmt.Println("No existen pedidos registrados en el sistema.")
		return
	}

	id := utilidades.LeerTexto("Ingrese la identificación del cliente: ")
	idLimpia := strings.TrimSpace(id)

	encontrados := []modelos.Pedido{}
	for _, p := range pedidos {
		if strings.EqualFold(p.GetCliente().GetIdentificacion(), idLimpia) {
			encontrados = append(encontrados, p)
		}
	}

	if len(encontrados) == 0 {
		fmt.Printf("No se encontraron pedidos para el cliente con identificación '%s'.\n", idLimpia)
		return
	}

	fmt.Printf("\nPedidos encontrados para el cliente '%s':\n", encontrados[0].GetCliente().GetNombre())
	fmt.Printf("%-12s | %-12s | %-20s\n", "CÓDIGO", "TOTAL", "FECHA")
	fmt.Println(strings.Repeat("-", 50))

	totalAcumulado := 0.0
	for _, p := range encontrados {
		fmt.Printf("%-12s | $%11.2f | %-20s\n",
			p.GetCodigo(), p.GetTotal(), p.GetFecha())
		totalAcumulado += p.GetTotal()
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Pedidos encontrados: %d | Total acumulado: $%.2f\n", len(encontrados), totalAcumulado)
}

// mostrarDetallePedido imprime el detalle completo de un pedido individual.
func mostrarDetallePedido(pedido modelos.Pedido) {
	fmt.Println("\n========================================")
	fmt.Println("        DETALLE DEL PEDIDO")
	fmt.Println("========================================")
	fmt.Printf("Código:  %s\n", pedido.GetCodigo())
	fmt.Printf("Cliente: %s (%s)\n", pedido.GetCliente().GetNombre(), pedido.GetCliente().GetIdentificacion())
	fmt.Printf("Fecha:   %s\n", pedido.GetFecha())
	fmt.Println()

	elementos := pedido.GetElementos()
	fmt.Printf("%-25s | %-8s | %-12s\n", "PRODUCTO", "CANT.", "SUBTOTAL")
	fmt.Println(strings.Repeat("-", 50))

	for _, elem := range elementos {
		fmt.Printf("%-25s | %8d | $%11.2f\n",
			elem.GetProducto().GetNombre(),
			elem.GetCantidad(),
			elem.GetSubtotal())
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("TOTAL COBRADO: $%.2f\n", pedido.GetTotal())
	fmt.Println("========================================")
}

// ----------------------------------------------------------------------------
// CONTROLADOR DEL SUBMENÚ DE PEDIDOS
// ----------------------------------------------------------------------------

// EjecutarModuloPedidos controla la navegación del submenú de gestión de pedidos.
func EjecutarModuloPedidos(pedidos *[]modelos.Pedido) {
	for {
		utilidades.MostrarSubmenuPedidos()
		opcion := utilidades.LeerOpcion("Seleccione una opción: ")

		switch opcion {
		case 1:
			ListarPedidos(*pedidos)
			utilidades.Pausar()
		case 2:
			BuscarPedidoMenu(*pedidos)
			utilidades.Pausar()
		case 3:
			BuscarPedidosPorCliente(*pedidos)
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
