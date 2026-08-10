package utilidades

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// lector es una instancia global de bufio.Reader configurada para leer desde la entrada estándar (consola).
// Reutilizar el mismo lector evita problemas de buffers no sincronizados.
var lector = bufio.NewReader(os.Stdin)

// LeerTexto solicita una entrada de texto al usuario, elimina espacios en blanco al inicio y final,
// y valida que la entrada no esté vacía. Se repite en bucle hasta recibir un texto válido.
func LeerTexto(mensaje string) string {
	for {
		fmt.Print(mensaje)
		entrada, err := lector.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada. Intente nuevamente.")
			continue
		}

		// Limpiar saltos de línea (\n y \r) y espacios extras
		texto := strings.TrimSpace(entrada)

		if texto == "" {
			fmt.Println("El campo no puede estar vacío. Intente nuevamente.")
			continue
		}

		return texto
	}
}

// LeerEntero solicita un número entero al usuario. Si el usuario ingresa letras o símbolos,
// captura el error con strconv.Atoi y solicita el número nuevamente sin cerrar el programa.
func LeerEntero(mensaje string) int {
	for {
		fmt.Print(mensaje)
		entrada, err := lector.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada. Intente nuevamente.")
			continue
		}

		texto := strings.TrimSpace(entrada)
		numero, err := strconv.Atoi(texto)
		if err != nil {
			fmt.Println("Error: Debe ingresar un número entero válido (sin letras ni decimales).")
			continue
		}

		return numero
	}
}

// LeerDecimal solicita un número decimal (float64) al usuario (ej: 15.50).
// Utiliza strconv.ParseFloat para la conversión y valida que el formato sea correcto.
func LeerDecimal(mensaje string) float64 {
	for {
		fmt.Print(mensaje)
		entrada, err := lector.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada. Intente nuevamente.")
			continue
		}

		texto := strings.TrimSpace(entrada)
		numero, err := strconv.ParseFloat(texto, 64)
		if err != nil {
			fmt.Println("Error: Debe ingresar un número decimal válido (ejemplo: 12.50).")
			continue
		}

		return numero
	}
}

// LeerOpcion lee la opción seleccionada en los menús de forma segura.
func LeerOpcion(mensaje string) int {
	return LeerEntero(mensaje)
}

// Confirmar realiza una pregunta de confirmación (sí/no) al usuario.
// Retorna true si la respuesta empieza con 's' o 'S', y false en cualquier otro caso.
func Confirmar(mensaje string) bool {
	fmt.Print(mensaje + " (s/n): ")
	entrada, err := lector.ReadString('\n')
	if err != nil {
		return false
	}
	respuesta := strings.ToLower(strings.TrimSpace(entrada))
	return respuesta == "s" || respuesta == "si"
}

// Pausar detiene temporalmente la ejecución del programa hasta que el usuario presione la tecla ENTER,
// lo que permite leer los mensajes de la consola antes de volver al menú.
func Pausar() {
	fmt.Println("\nPresione ENTER para continuar...")
	lector.ReadString('\n')
}
