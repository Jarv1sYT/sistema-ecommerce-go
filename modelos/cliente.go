package modelos

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Cliente representa a un comprador registrado en la tienda virtual.
// Se aplican campos PRIVADOS en minúscula para garantizar la ENCAPSULACIÓN (Unidad 3).
// Ningún paquete externo puede modificar directamente sus datos sin pasar por sus métodos.
type Cliente struct {
	identificacion string // Documento de identidad o cédula del cliente (único)
	nombre         string // Nombre completo del cliente
	correo         string // Correo electrónico de contacto
}

// ----------------------------------------------------------------------------
// ESTRUCTURA DTO AUXILIAR Y PERSISTENCIA JSON
// ----------------------------------------------------------------------------

// clienteDTO es una estructura auxiliar privada empleada EXCLUSIVAMENTE para la
// serialización y deserialización con encoding/json, conservando los nombres de campos en JSON.
type clienteDTO struct {
	Identificacion string `json:"identificacion"`
	Nombre         string `json:"nombre"`
	Correo         string `json:"correo"`
}

// MarshalJSON implementa json.Marshaler. Convierte el objeto Cliente privado en un DTO
// para que encoding/json genere la estructura JSON sin perder los atributos encapsulados.
func (c Cliente) MarshalJSON() ([]byte, error) {
	dto := clienteDTO{
		Identificacion: c.identificacion,
		Nombre:         c.nombre,
		Correo:         c.correo,
	}
	return json.Marshal(dto)
}

// UnmarshalJSON implementa json.Unmarshaler. Reconstruye el objeto Cliente privado
// mapeando los bytes leídos desde el archivo JSON hacia el DTO.
func (c *Cliente) UnmarshalJSON(b []byte) error {
	var dto clienteDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	c.identificacion = dto.Identificacion
	c.nombre = dto.Nombre
	c.correo = dto.Correo
	return nil
}

// ----------------------------------------------------------------------------
// CONSTRUCTOR
// ----------------------------------------------------------------------------

// ValidarCorreoFormato es una función auxiliar del modelo que verifica el formato básico de un correo.
func ValidarCorreoFormato(correo string) bool {
	correoLimpio := strings.TrimSpace(correo)
	tieneArroba := strings.Contains(correoLimpio, "@")
	tienePunto := strings.Contains(correoLimpio, ".")
	return tieneArroba && tienePunto && len(correoLimpio) >= 5
}

// NuevoCliente es la función constructora que valida el estado inicial de un Cliente
// antes de crearlo. Retorna un puntero al Cliente o un error si alguna regla no se cumple.
func NuevoCliente(identificacion, nombre, correo string) (*Cliente, error) {
	idLimpia := strings.TrimSpace(identificacion)
	if idLimpia == "" {
		return nil, errors.New("la identificación del cliente no puede estar vacía")
	}

	nombreLimpio := strings.TrimSpace(nombre)
	if nombreLimpio == "" {
		return nil, errors.New("el nombre del cliente no puede estar vacío")
	}

	correoLimpio := strings.ToLower(strings.TrimSpace(correo))
	if !ValidarCorreoFormato(correoLimpio) {
		return nil, errors.New("el correo electrónico no tiene un formato válido (debe contener '@' y '.')")
	}

	return &Cliente{
		identificacion: idLimpia,
		nombre:         nombreLimpio,
		correo:         correoLimpio,
	}, nil
}

// ----------------------------------------------------------------------------
// GETTERS (Permiten consultar la información privada)
// ----------------------------------------------------------------------------

// GetIdentificacion permite consultar el documento o cédula del cliente.
func (c Cliente) GetIdentificacion() string {
	return c.identificacion
}

// GetNombre permite consultar el nombre completo del cliente.
func (c Cliente) GetNombre() string {
	return c.nombre
}

// GetCorreo permite consultar el correo electrónico del cliente.
func (c Cliente) GetCorreo() string {
	return c.correo
}

// ----------------------------------------------------------------------------
// SETTERS (Permiten modificar los datos protegiendo la integridad del objeto)
// ----------------------------------------------------------------------------

// SetNombre modifica el nombre del cliente previa validación.
func (c *Cliente) SetNombre(nuevoNombre string) error {
	nombreLimpio := strings.TrimSpace(nuevoNombre)
	if nombreLimpio == "" {
		return errors.New("el nuevo nombre del cliente no puede estar vacío")
	}
	c.nombre = nombreLimpio
	return nil
}

// SetCorreo modifica el correo electrónico del cliente previa validación.
func (c *Cliente) SetCorreo(nuevoCorreo string) error {
	correoLimpio := strings.ToLower(strings.TrimSpace(nuevoCorreo))
	if !ValidarCorreoFormato(correoLimpio) {
		return errors.New("el nuevo correo electrónico no tiene un formato válido (debe contener '@' y '.')")
	}
	c.correo = correoLimpio
	return nil
}

// ----------------------------------------------------------------------------
// MÉTODOS E INTERFACES (Unidad 3)
// ----------------------------------------------------------------------------

// MostrarFicha imprime los datos del cliente utilizando sus métodos de consulta.
// Este método permite a Cliente cumplir con la interfaz Mostrable.
func (c Cliente) MostrarFicha() {
	fmt.Printf("Identificación: %-12s | Nombre: %-25s | Correo: %s\n", c.GetIdentificacion(), c.GetNombre(), c.GetCorreo())
}
