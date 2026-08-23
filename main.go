package main

import (
	"fmt"

	"sistema-ecommerce-go/api"
	"sistema-ecommerce-go/modelos"
	"sistema-ecommerce-go/servicios"
	"sistema-ecommerce-go/utilidades"
)

func main() {
	// Carga inicial de datos persistidos en archivos JSON
	productos, err := servicios.CargarProductos()
	if err != nil {
		fmt.Printf("Aviso: No se pudieron cargar los productos existentes: %v\n", err)
		productos = []modelos.Producto{}
	}

	clientes, err := servicios.CargarClientes()
	if err != nil {
		fmt.Printf("Aviso: No se pudieron cargar los clientes existentes: %v\n", err)
		clientes = []modelos.Cliente{}
	}

	pedidos, err := servicios.CargarPedidos()
	if err != nil {
		fmt.Printf("Aviso: No se pudieron cargar los pedidos existentes: %v\n", err)
		pedidos = []modelos.Pedido{}
	}

	// Crear el carrito de compras en memoria (no se persiste entre sesiones)
	carrito := modelos.NuevoCarrito()

	fmt.Println("¡Datos cargados correctamente desde la carpeta datos/!")

	// =========================================================================
	// INICIAR SERVIDOR API REST EN SEGUNDO PLANO (Goroutine)
	// =========================================================================
	// El servidor web comparte el mismo estado que la consola a través de EstadoAPI.
	// Se lanza como goroutine para que la consola interactiva siga funcionando.
	estado := api.NuevoEstadoAPI(productos, clientes, pedidos)
	go api.IniciarServidor(estado, "8080")

	// =========================================================================
	// BUCLE INTERACTIVO PRINCIPAL (CONSOLA)
	// =========================================================================
	for {
		utilidades.MostrarMenuPrincipal()
		opcion := utilidades.LeerOpcion("Seleccione una opción: ")

		switch opcion {
		case 1:
			servicios.EjecutarModuloProductos(&productos)
			// Sincronizar cambios con el estado de la API
			estado.Mu.Lock()
			estado.Productos = productos
			estado.Mu.Unlock()
		case 2:
			servicios.EjecutarModuloClientes(&clientes)
			// Sincronizar cambios con el estado de la API
			estado.Mu.Lock()
			estado.Clientes = clientes
			estado.Mu.Unlock()
		case 3:
			servicios.EjecutarModuloInventario(productos)
		case 4:
			servicios.EjecutarModuloCarrito(carrito, &productos, clientes)
			// Recargar pedidos después de posibles compras
			pedidos, err = servicios.CargarPedidos()
			if err != nil {
				fmt.Printf("Aviso: No se pudieron recargar los pedidos: %v\n", err)
			}
			// Sincronizar con el estado de la API
			estado.Mu.Lock()
			estado.Productos = productos
			estado.Pedidos = pedidos
			estado.Mu.Unlock()
		case 5:
			servicios.EjecutarModuloPedidos(&pedidos)
		case 6:
			fmt.Println("\nMódulo de reportes en desarrollo...")
			utilidades.Pausar()
		case 0:
			fmt.Println("\n==================================================")
			fmt.Println(" Gracias por utilizar el Sistema de E-Commerce")
			fmt.Println("==================================================")
			return
		default:
			fmt.Println("Opción no válida. Intente nuevamente.")
			utilidades.Pausar()
		}
	}
}