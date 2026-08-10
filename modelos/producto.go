package modelos

import (
	"encoding/json"
	"errors"
	"strings"
)

// Producto representa un artículo disponible en la tienda virtual.
// Se aplican campos PRIVADOS en minúscula para garantizar la ENCAPSULACIÓN (Unidad 3).
// Ningún paquete externo puede modificar directamente su estado sin pasar por las reglas del modelo.
type Producto struct {
	codigo             string  // Identificador único del producto (ej: "P001")
	nombre             string  // Nombre descriptivo del producto
	precio             float64 // Precio unitario en valor monetario (debe ser > 0)
	cantidadDisponible int     // Existencias o stock disponible (debe ser >= 0)
}

// ----------------------------------------------------------------------------
// ESTRUCTURA DTO AUXILIAR Y PERSISTENCIA JSON
// ----------------------------------------------------------------------------

// productoDTO es una estructura auxiliar privada utilizada EXCLUSIVAMENTE para la
// serialización y deserialización con encoding/json, evitando exponer los campos
// del modelo original Producto hacia otros paquetes.
type productoDTO struct {
	Codigo             string  `json:"codigo"`
	Nombre             string  `json:"nombre"`
	Precio             float64 `json:"precio"`
	CantidadDisponible int     `json:"cantidad_disponible"`
}

// MarshalJSON implementa json.Marshaler. Convierte el Producto privado en un DTO
// para que encoding/json genere la estructura JSON esperada.
func (p Producto) MarshalJSON() ([]byte, error) {
	dto := productoDTO{
		Codigo:             p.codigo,
		Nombre:             p.nombre,
		Precio:             p.precio,
		CantidadDisponible: p.cantidadDisponible,
	}
	return json.Marshal(dto)
}

// UnmarshalJSON implementa json.Unmarshaler. Recibe los bytes desde la lectura JSON,
// los mapea al DTO y reconstruye el objeto Producto privado.
func (p *Producto) UnmarshalJSON(b []byte) error {
	var dto productoDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	p.codigo = dto.Codigo
	p.nombre = dto.Nombre
	p.precio = dto.Precio
	p.cantidadDisponible = dto.CantidadDisponible
	return nil
}

// ----------------------------------------------------------------------------
// CONSTRUCTOR
// ----------------------------------------------------------------------------

// NuevoProducto es la función constructora que valida el estado inicial de un Producto
// antes de crearlo. Retorna un puntero al Producto o un error si alguna regla no se cumple.
func NuevoProducto(codigo, nombre string, precio float64, cantidadDisponible int) (*Producto, error) {
	codigoLimpio := strings.ToUpper(strings.TrimSpace(codigo))
	if codigoLimpio == "" {
		return nil, errors.New("el código del producto no puede estar vacío")
	}

	nombreLimpio := strings.TrimSpace(nombre)
	if nombreLimpio == "" {
		return nil, errors.New("el nombre del producto no puede estar vacío")
	}

	if precio <= 0 {
		return nil, errors.New("el precio del producto debe ser mayor que cero")
	}

	if cantidadDisponible < 0 {
		return nil, errors.New("la cantidad disponible no puede ser negativa")
	}

	return &Producto{
		codigo:             codigoLimpio,
		nombre:             nombreLimpio,
		precio:             precio,
		cantidadDisponible: cantidadDisponible,
	}, nil
}

// ----------------------------------------------------------------------------
// GETTERS (Permiten consultar la información privada)
// ----------------------------------------------------------------------------

// GetCodigo permite consultar el código único del producto.
func (p Producto) GetCodigo() string {
	return p.codigo
}

// GetNombre permite consultar el nombre descriptivo del producto.
func (p Producto) GetNombre() string {
	return p.nombre
}

// GetPrecio permite consultar el precio unitario del producto.
func (p Producto) GetPrecio() float64 {
	return p.precio
}

// GetCantidadDisponible permite consultar el stock actual del producto.
func (p Producto) GetCantidadDisponible() int {
	return p.cantidadDisponible
}

// ----------------------------------------------------------------------------
// SETTERS (Permiten modificar la información protegiendo el estado interno)
// ----------------------------------------------------------------------------

// SetNombre modifica el nombre del producto previa validación.
func (p *Producto) SetNombre(nuevoNombre string) error {
	nombreLimpio := strings.TrimSpace(nuevoNombre)
	if nombreLimpio == "" {
		return errors.New("el nuevo nombre no puede estar vacío")
	}
	p.nombre = nombreLimpio
	return nil
}

// SetPrecio modifica el precio del producto previa validación.
func (p *Producto) SetPrecio(nuevoPrecio float64) error {
	if nuevoPrecio <= 0 {
		return errors.New("el precio debe ser mayor que cero")
	}
	p.precio = nuevoPrecio
	return nil
}

// SetCantidadDisponible modifica el stock del producto previa validación.
func (p *Producto) SetCantidadDisponible(nuevaCantidad int) error {
	if nuevaCantidad < 0 {
		return errors.New("la cantidad disponible no puede ser negativa")
	}
	p.cantidadDisponible = nuevaCantidad
	return nil
}

// ----------------------------------------------------------------------------
// MÉTODOS DE COMPORTAMIENTO DEL NEGOCIO
// ----------------------------------------------------------------------------

// TieneStockSuficiente comprueba si el producto cuenta con existencias para cubrir una cantidad solicitada.
func (p Producto) TieneStockSuficiente(cantidad int) bool {
	return cantidad > 0 && p.cantidadDisponible >= cantidad
}

// DescontarStock reduce las existencias disponibles según las unidades vendidas.
func (p *Producto) DescontarStock(cantidad int) error {
	if !p.TieneStockSuficiente(cantidad) {
		return errors.New("inventario insuficiente para realizar el descuento")
	}
	p.cantidadDisponible -= cantidad
	return nil
}
