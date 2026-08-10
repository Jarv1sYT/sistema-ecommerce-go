package modelos

import (
	"encoding/json"
	"errors"
	"strings"
)

// Pedido representa una transacción de compra finalizada e inmutable en el sistema.
// Sus atributos son PRIVADOS para evitar modificaciones accidentales posteriores a la venta (Unidad 3).
type Pedido struct {
	codigo    string            // Código único de la orden (ej: "PED-1001")
	cliente   Cliente           // Datos del cliente que efectuó la compra
	elementos []ElementoCarrito // Copia inmutable de los productos adquiridos
	total     float64           // Valor total cobrado
	fecha     string            // Fecha y hora de confirmación (AAAA-MM-DD HH:MM:SS)
}

// ----------------------------------------------------------------------------
// ESTRUCTURA DTO AUXILIAR Y PERSISTENCIA JSON
// ----------------------------------------------------------------------------

type pedidoDTO struct {
	Codigo    string            `json:"codigo"`
	Cliente   Cliente           `json:"cliente"`
	Elementos []ElementoCarrito `json:"elementos"`
	Total     float64           `json:"total"`
	Fecha     string            `json:"fecha"`
}

// MarshalJSON serializa el pedido a JSON utilizando el DTO auxiliar.
func (p Pedido) MarshalJSON() ([]byte, error) {
	dto := pedidoDTO{
		Codigo:    p.codigo,
		Cliente:   p.cliente,
		Elementos: p.elementos,
		Total:     p.total,
		Fecha:     p.fecha,
	}
	return json.Marshal(dto)
}

// UnmarshalJSON deserializa los bytes JSON y reconstruye el Pedido encapsulado.
func (p *Pedido) UnmarshalJSON(b []byte) error {
	var dto pedidoDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	p.codigo = dto.Codigo
	p.cliente = dto.Cliente
	p.elementos = dto.Elementos
	p.total = dto.Total
	p.fecha = dto.Fecha
	return nil
}

// ----------------------------------------------------------------------------
// CONSTRUCTOR (Única vía controlada para crear un Pedido)
// ----------------------------------------------------------------------------

// NuevoPedido construye una orden inmutable. Valida la existencia de productos y calcula el total final.
func NuevoPedido(codigo string, cliente Cliente, elementos []ElementoCarrito, fecha string) (*Pedido, error) {
	codigoLimpio := strings.ToUpper(strings.TrimSpace(codigo))
	if codigoLimpio == "" {
		return nil, errors.New("el código del pedido no puede estar vacío")
	}

	if len(elementos) == 0 {
		return nil, errors.New("no se puede crear un pedido sin productos en la compra")
	}

	fechaLimpia := strings.TrimSpace(fecha)
	if fechaLimpia == "" {
		return nil, errors.New("la fecha del pedido es obligatoria")
	}

	// Copia defensiva de los elementos para asegurar la inmutabilidad
	elementosCopia := make([]ElementoCarrito, len(elementos))
	copy(elementosCopia, elementos)

	// Calcular el total del pedido sumando subtotales
	totalCalculado := 0.0
	for _, elem := range elementosCopia {
		totalCalculado += elem.GetSubtotal()
	}

	return &Pedido{
		codigo:    codigoLimpio,
		cliente:   cliente,
		elementos: elementosCopia,
		total:     totalCalculado,
		fecha:     fechaLimpia,
	}, nil
}

// ----------------------------------------------------------------------------
// GETTERS (Sin Setters para garantizar Inmutabilidad)
// ----------------------------------------------------------------------------

// GetCodigo retorna el identificador único del pedido.
func (p Pedido) GetCodigo() string {
	return p.codigo
}

// GetCliente retorna una copia de los datos del cliente comprador.
func (p Pedido) GetCliente() Cliente {
	return p.cliente
}

// GetElementos retorna una COPIA SEGURA del slice de elementos para evitar mutaciones externas.
func (p Pedido) GetElementos() []ElementoCarrito {
	copia := make([]ElementoCarrito, len(p.elementos))
	copy(copia, p.elementos)
	return copia
}

// GetTotal retorna el monto total cobrado en la transacción.
func (p Pedido) GetTotal() float64 {
	return p.total
}

// GetFecha retorna la fecha y hora congelada del pedido.
func (p Pedido) GetFecha() string {
	return p.fecha
}
