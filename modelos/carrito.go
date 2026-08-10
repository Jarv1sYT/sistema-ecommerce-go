package modelos

import (
	"encoding/json"
	"errors"
	"strings"
)

// ElementoCarrito representa un ítem individual en el carrito de compras.
// Sus atributos son PRIVADOS en minúscula para garantizar la ENCAPSULACIÓN (Unidad 3).
type ElementoCarrito struct {
	producto Producto // Información del producto seleccionado
	cantidad int      // Unidades solicitadas (debe ser > 0)
	subtotal float64  // Calculado automáticamente como precio * cantidad
}

// ----------------------------------------------------------------------------
// ESTRUCTURA DTO Y JSON PARA ElementoCarrito
// ----------------------------------------------------------------------------

type elementoCarritoDTO struct {
	Producto Producto `json:"producto"`
	Cantidad int      `json:"cantidad"`
	Subtotal float64  `json:"subtotal"`
}

func (e ElementoCarrito) MarshalJSON() ([]byte, error) {
	dto := elementoCarritoDTO{
		Producto: e.producto,
		Cantidad: e.cantidad,
		Subtotal: e.subtotal,
	}
	return json.Marshal(dto)
}

func (e *ElementoCarrito) UnmarshalJSON(b []byte) error {
	var dto elementoCarritoDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	e.producto = dto.Producto
	e.cantidad = dto.Cantidad
	e.subtotal = dto.Subtotal
	return nil
}

// ----------------------------------------------------------------------------
// CONSTRUCTOR Y MÉTODOS DE ElementoCarrito
// ----------------------------------------------------------------------------

// NuevoElementoCarrito construye un ítem del carrito calculando automáticamente su subtotal.
func NuevoElementoCarrito(producto Producto, cantidad int) (*ElementoCarrito, error) {
	if cantidad <= 0 {
		return nil, errors.New("la cantidad elegida debe ser mayor que cero")
	}

	subtotal := producto.GetPrecio() * float64(cantidad)

	return &ElementoCarrito{
		producto: producto,
		cantidad: cantidad,
		subtotal: subtotal,
	}, nil
}

func (e ElementoCarrito) GetProducto() Producto {
	return e.producto
}

func (e ElementoCarrito) GetCantidad() int {
	return e.cantidad
}

func (e ElementoCarrito) GetSubtotal() float64 {
	return e.subtotal
}

func (e *ElementoCarrito) SetCantidad(nuevaCantidad int) error {
	if nuevaCantidad <= 0 {
		return errors.New("la cantidad debe ser mayor que cero")
	}
	e.cantidad = nuevaCantidad
	e.subtotal = e.producto.GetPrecio() * float64(nuevaCantidad)
	return nil
}

// ============================================================================
// ESTRUCTURA Carrito
// ============================================================================

// Carrito representa el contenedor de compras temporales.
// Mantiene sus atributos PRIVADOS para controlar la adición, eliminación y recálculo de totales.
type Carrito struct {
	cliente   *Cliente          // Cliente asociado al carrito (opcional hasta confirmar)
	elementos []ElementoCarrito // Lista interna de ítems
	total     float64           // Total acumulado
}

// ----------------------------------------------------------------------------
// ESTRUCTURA DTO Y JSON PARA Carrito
// ----------------------------------------------------------------------------

type carritoDTO struct {
	Cliente   *Cliente          `json:"cliente,omitempty"`
	Elementos []ElementoCarrito `json:"elementos"`
	Total     float64           `json:"total"`
}

func (c Carrito) MarshalJSON() ([]byte, error) {
	dto := carritoDTO{
		Cliente:   c.cliente,
		Elementos: c.elementos,
		Total:     c.total,
	}
	return json.Marshal(dto)
}

func (c *Carrito) UnmarshalJSON(b []byte) error {
	var dto carritoDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	c.cliente = dto.Cliente
	c.elementos = dto.Elementos
	c.total = dto.Total
	return nil
}

// ----------------------------------------------------------------------------
// CONSTRUCTOR Y MÉTODOS DE Carrito
// ----------------------------------------------------------------------------

// NuevoCarrito crea una nueva instancia limpia de Carrito.
func NuevoCarrito() *Carrito {
	return &Carrito{
		cliente:   nil,
		elementos: []ElementoCarrito{},
		total:     0.0,
	}
}

// SetCliente asocia el cliente comprador al carrito.
func (c *Carrito) SetCliente(cliente *Cliente) {
	c.cliente = cliente
}

// GetCliente retorna el cliente asociado (o nil si no se ha asignado).
func (c Carrito) GetCliente() *Cliente {
	return c.cliente
}

// GetElementos retorna una COPIA SEGURA de la lista de elementos para preservar la encapsulación.
func (c Carrito) GetElementos() []ElementoCarrito {
	copia := make([]ElementoCarrito, len(c.elementos))
	copy(copia, c.elementos)
	return copia
}

// GetTotal retorna el valor total del carrito.
func (c Carrito) GetTotal() float64 {
	return c.total
}

// EsVacio comprueba si el carrito no tiene productos agregados.
func (c Carrito) EsVacio() bool {
	return len(c.elementos) == 0
}

// CalcularTotal recalcula la suma acumulada de los subtotales.
func (c *Carrito) CalcularTotal() float64 {
	suma := 0.0
	for _, elem := range c.elementos {
		suma += elem.GetSubtotal()
	}
	c.total = suma
	return c.total
}

// AgregarElemento incluye un producto en el carrito validando stock y recalculando el total.
func (c *Carrito) AgregarElemento(producto Producto, cantidad int) error {
	if cantidad <= 0 {
		return errors.New("la cantidad a agregar debe ser mayor que cero")
	}

	if !producto.TieneStockSuficiente(cantidad) {
		return errors.New("no hay stock suficiente en inventario para agregar esa cantidad")
	}

	// Verificar si el producto ya está en el carrito para sumar la cantidad
	for i, elem := range c.elementos {
		if strings.EqualFold(elem.GetProducto().GetCodigo(), producto.GetCodigo()) {
			nuevaCantidadTotal := elem.GetCantidad() + cantidad
			if !producto.TieneStockSuficiente(nuevaCantidadTotal) {
				return errors.New("la cantidad total solicitada supera el stock disponible en inventario")
			}
			err := c.elementos[i].SetCantidad(nuevaCantidadTotal)
			if err != nil {
				return err
			}
			c.CalcularTotal()
			return nil
		}
	}

	// Si es un ítem nuevo, crearlo mediante su constructor y anexarlo
	nuevoItem, err := NuevoElementoCarrito(producto, cantidad)
	if err != nil {
		return err
	}

	c.elementos = append(c.elementos, *nuevoItem)
	c.CalcularTotal()
	return nil
}

// EliminarElemento remueve un producto del carrito por su código y actualiza el total.
func (c *Carrito) EliminarElemento(codigoProducto string) bool {
	codigoLimpio := strings.ToLower(strings.TrimSpace(codigoProducto))
	for i, elem := range c.elementos {
		if strings.ToLower(elem.GetProducto().GetCodigo()) == codigoLimpio {
			c.elementos = append(c.elementos[:i], c.elementos[i+1:]...)
			c.CalcularTotal()
			return true
		}
	}
	return false
}

// Vaciar reinicia la lista de elementos y resetea el total a cero.
func (c *Carrito) Vaciar() {
	c.elementos = []ElementoCarrito{}
	c.cliente = nil
	c.total = 0.0
}
