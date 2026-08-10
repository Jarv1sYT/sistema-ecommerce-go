package servicios

import (
	"fmt"
	"strings"

	"sistema-ecommerce-go/modelos"
	"sistema-ecommerce-go/utilidades"
)

// Ruta donde se almacena el archivo de persistencia de clientes
const rutaClientesJSON = "datos/clientes.json"

// Mostrable representa la interfaz para cualquier entidad que pueda imprimir su ficha (Unidad 3: Interfaces)
type Mostrable interface {
	MostrarFicha()
}

// CargarClientes lee la lista de clientes registrados desde el archivo JSON local.
func CargarClientes() ([]modelos.Cliente, error) {
	var clientes []modelos.Cliente
	err := utilidades.LeerJSON(rutaClientesJSON, &clientes)
	if err != nil {
		return nil, err
	}
	return clientes, nil
}

// GuardarClientes serializa y almacena el listado de clientes en el archivo JSON.
func GuardarClientes(clientes []modelos.Cliente) error {
	return utilidades.GuardarJSON(rutaClientesJSON, clientes)
}

// BuscarClientePorIdentificacion busca un cliente dentro del slice utilizando el getter GetIdentificacion().
// Retorna la estructura del cliente, su posición (índice) y un booleano que indica si existe.
func BuscarClientePorIdentificacion(clientes []modelos.Cliente, id string) (modelos.Cliente, int, bool) {
	idLimpia := strings.TrimSpace(id)
	for i, c := range clientes {
		if strings.EqualFold(c.GetIdentificacion(), idLimpia) {
			return c, i, true
		}
	}
	return modelos.Cliente{}, -1, false
}

// RegistrarCliente solicita los datos de un nuevo cliente, los valida mediante el constructor modelos.NuevoCliente,
// y lo agrega al slice persistiendo los datos de inmediato en JSON.
func RegistrarCliente(clientes *[]modelos.Cliente) {
	fmt.Println("\n--- REGISTRAR NUEVO CLIENTE ---")

	// 1. Validar identificación única y no vacía
	var identificacion string
	for {
		identificacion = utilidades.LeerTexto("Ingrese la cédula o identificación del cliente: ")
		_, _, existe := BuscarClientePorIdentificacion(*clientes, identificacion)
		if existe {
			fmt.Println("Error: Ya existe un cliente registrado con esa identificación. Intente nuevamente.")
			continue
		}
		break
	}

	// 2. Solicitar nombre completo
	nombre := utilidades.LeerTexto("Ingrese el nombre completo del cliente: ")

	// 3. Solicitar correo electrónico
	correo := utilidades.LeerTexto("Ingrese el correo electrónico del cliente: ")

	// Crear el objeto cliente utilizando su CONSTRUCTOR (valida internamente formato y datos)
	nuevoCliente, err := modelos.NuevoCliente(identificacion, nombre, correo)
	if err != nil {
		fmt.Printf("\nError de validación al registrar el cliente: %v\n", err)
		return
	}

	// Agregar cliente al slice por puntero
	*clientes = append(*clientes, *nuevoCliente)

	// Persistir de inmediato en archivo JSON
	err = GuardarClientes(*clientes)
	if err != nil {
		fmt.Printf("Error al guardar el cliente en JSON: %v\n", err)
		return
	}

	fmt.Println("\n¡Cliente registrado y guardado exitosamente en datos/clientes.json!")
}

// ListarClientes muestra todos los clientes registrados utilizando la interfaz Mostrable.
func ListarClientes(clientes []modelos.Cliente) {
	fmt.Println("\n--- LISTADO DE CLIENTES REGISTRADOS ---")

	if len(clientes) == 0 {
		fmt.Println("No existen clientes registrados en el sistema.")
		return
	}

	fmt.Println(strings.Repeat("-", 70))
	for _, c := range clientes {
		// Aplicación práctica de Interfaz Mostrable (Unidad 3)
		var m Mostrable = c
		m.MostrarFicha()
	}
	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("Total de clientes registrados: %d\n", len(clientes))
}

// BuscarClienteMenu solicita una identificación y muestra la información del cliente encontrado.
func BuscarClienteMenu(clientes []modelos.Cliente) {
	fmt.Println("\n--- BUSCAR CLIENTE POR IDENTIFICACIÓN ---")

	if len(clientes) == 0 {
		fmt.Println("No existen clientes registrados en el sistema.")
		return
	}

	id := utilidades.LeerTexto("Ingrese la identificación del cliente a buscar: ")
	cliente, _, encontrado := BuscarClientePorIdentificacion(clientes, id)

	if !encontrado {
		fmt.Printf("No se encontró ningún cliente con la identificación '%s'.\n", id)
		return
	}

	fmt.Println("\n--- DETALLES DEL CLIENTE ENCONTRADO ---")
	var m Mostrable = cliente
	m.MostrarFicha()
}

// EjecutarModuloClientes controla la navegación y opciones del submenú de clientes.
func EjecutarModuloClientes(clientes *[]modelos.Cliente) {
	for {
		utilidades.MostrarSubmenuClientes()
		opcion := utilidades.LeerOpcion("Seleccione una opción: ")

		switch opcion {
		case 1:
			RegistrarCliente(clientes)
			utilidades.Pausar()
		case 2:
			ListarClientes(*clientes)
			utilidades.Pausar()
		case 3:
			BuscarClienteMenu(*clientes)
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
