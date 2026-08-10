package main

import (
	"fmt"

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

	fmt.Println("¡Datos cargados correctamente desde la carpeta datos/!")

	// Bucle interactivo principal de la aplicación
	for {
		utilidades.MostrarMenuPrincipal()
		opcion := utilidades.LeerOpcion("Seleccione una opción: ")

		switch opcion {
		case 1:
			servicios.EjecutarModuloProductos(&productos)
		case 2:
			servicios.EjecutarModuloClientes(&clientes)
		case 3:
			fmt.Println("\nMódulo de inventario en desarrollo...")
			utilidades.Pausar()
		case 4:
			fmt.Println("\nMódulo de carrito de compras en desarrollo...")
			utilidades.Pausar()
		case 5:
			fmt.Println("\nMódulo de gestión de pedidos en desarrollo...")
			utilidades.Pausar()
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