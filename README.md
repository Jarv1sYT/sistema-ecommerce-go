# 🛒 Sistema de Gestión de E-Commerce

> Proyecto integrador de la asignatura **Programación Orientada a Objetos**

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![Tests](https://img.shields.io/badge/Tests-44%20passed-brightgreen?style=flat)](./servicios/)
[![Endpoints](https://img.shields.io/badge/API%20REST-8%20endpoints-blue?style=flat)](./api/)
[![License](https://img.shields.io/badge/License-Academic-yellow?style=flat)]()

---

## 📋 Descripción

Sistema completo de gestión de comercio electrónico desarrollado en **Go (Golang)**, que implementa los 4 pilares fundamentales de la asignatura:

El software opera en **dos modos simultáneos**:
- 🖥️ **Consola interactiva** con menús y navegación por terminal.
- 🌐 **Servidor API REST** en `http://localhost:8080` con 8 endpoints JSON.

---

## 👤 Autor

| Nombre |
|--------|
| **Marlon Rivera**

---

## 🏗️ Arquitectura del Software

```
sistema-ecommerce-go/
│
├── main.go                        # Punto de entrada (orquestador)
├── go.mod                         # Módulo Go
│
├── modelos/                       # Capa de Dominio (POO + Encapsulación)
│   ├── producto.go                # Struct Producto (campos privados + getters/setters)
│   ├── cliente.go                 # Struct Cliente (interfaz Mostrable)
│   ├── carrito.go                 # Struct Carrito + ElementoCarrito
│   └── pedido.go                  # Struct Pedido (inmutable, solo getters)
│
├── servicios/                     # Capa de Lógica de Negocio
│   ├── productos.go               # CRUD de productos
│   ├── clientes.go                # CRUD de clientes + interfaz Mostrable
│   ├── inventario.go              # Consultas de inventario
│   ├── carrito.go                 # Gestión del carrito + Checkout transaccional
│   ├── pedidos.go                 # Historial y búsqueda de pedidos
│   ├── productos_test.go          # 15 pruebas unitarias
│   └── pedidos_test.go            # 14 pruebas unitarias
│
├── utilidades/                    # Capa de Soporte Transversal
│   ├── entradas.go                # Lectura segura por consola (bufio.Reader)
│   ├── archivos.go                # Persistencia en archivos JSON
│   └── menu.go                    # Interfaces de menús en consola
│
├── api/                           # Capa Web (Unidad 4)
│   ├── servidor.go                # Configuración del servidor HTTP y rutas
│   ├── handlers.go                # 8 Endpoints REST con JSON
│   ├── worker.go                  # Worker Pool (Goroutines + Channels + Mutex)
│   └── handlers_test.go           # 16 pruebas de integración HTTP
│
└── datos/                         # Persistencia en disco
    ├── productos.json
    ├── clientes.json
    └── pedidos.json
```
### Principios de Diseño Aplicados

- **Encapsulación estricta**: Todos los campos de las estructuras son privados (minúsculas). Solo se accede mediante getters y setters controlados con validaciones.
- **Constructores seguros**: Cada entidad se instancia exclusivamente a través de funciones constructoras (`NuevoProducto`, `NuevoCliente`, `NuevoPedido`) que validan el estado inicial.
- **Inmutabilidad**: La estructura `Pedido` no tiene setters; una vez creada, no puede modificarse.
- **Copias defensivas**: `GetElementos()` en `Carrito` y `Pedido` retorna copias del slice interno para proteger el estado encapsulado.
- **Serialización encapsulada**: Se utilizan DTOs privados + `MarshalJSON`/`UnmarshalJSON` para serializar campos privados sin exponerlos.

---

## 🚀 Guía de Instalación y Ejecución

### Requisitos Previos

- [Go 1.22 o superior](https://go.dev/dl/) instalado.
- Terminal de comandos (PowerShell, CMD o Bash).

### Clonar el Repositorio

```bash
git clone https://github.com/[usuario]/sistema-ecommerce-go.git
cd sistema-ecommerce-go
```
### Compilar y Ejecutar

```bash
# Ejecutar directamente
go run main.go

# O compilar y luego ejecutar el binario
go build -o ecommerce.exe
./ecommerce.exe
```

Al iniciar, el sistema:
1. Carga los datos existentes desde `datos/*.json`.
2. Inicia el **servidor API REST** en `http://localhost:8080`.
3. Muestra el **menú interactivo** en la consola.

### Ejecutar las Pruebas Automatizadas

```bash
# Ejecutar todos los 44 tests
go test ./...

# Ejecutar con detalle completo
go test -v ./...

# Solo tests de lógica de negocio
go test -v ./servicios/

# Solo tests de integración HTTP
go test -v ./api/
```
---

## 🌐 Documentación de Endpoints REST (8 Servicios Web)

El servidor escucha en `http://localhost:8080`. Todos los endpoints responden en formato JSON.

### Estructura de Respuesta Estándar

```json
{
  "exito": true,
  "mensaje": "Descripción del resultado",
  "datos": { }
}
```
### Ejemplos con curl

```bash
# Registrar producto
curl -X POST http://localhost:8080/api/productos \
  -H "Content-Type: application/json" \
  -d '{"codigo":"P001","nombre":"Laptop HP","precio":899.99,"cantidad":10}'

# Agregar al carrito
curl -X POST http://localhost:8080/api/carrito/agregar \
  -H "Content-Type: application/json" \
  -d '{"codigo_producto":"P001","cantidad":2}'

# Confirmar compra
curl -X POST http://localhost:8080/api/pedidos/checkout \
  -H "Content-Type: application/json" \
  -d '{"identificacion_cliente":"1234567890"}'
```

---

## ⚡ Concurrencia: Worker Pool con Goroutines y Channels

El sistema implementa un **Worker Pool** para procesar pedidos de forma concurrente sin bloquear las respuestas HTTP.

## 📁 Persistencia de Datos

Los datos se almacenan en archivos JSON dentro de la carpeta `datos/`:

| Archivo | Contenido |
|---------|-----------|
| `datos/productos.json` | Inventario de productos con stock actualizado |
| `datos/clientes.json`  | Clientes registrados en el sistema |
| `datos/pedidos.json`   | Historial completo de compras realizadas |

La carpeta `datos/` se crea automáticamente al ejecutar la aplicación por primera vez.

---

## 📄 Licencia

Proyecto académico desarrollado como parte del proyecto integrador de la asignatura Programación Orientada a Objetos, desarollado por Marlon Rivera.
