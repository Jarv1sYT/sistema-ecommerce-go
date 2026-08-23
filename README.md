# 🛒 Sistema de Gestión de E-Commerce

> Proyecto integrador de la asignatura **Programación Orientada a Objetos**

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![Tests](https://img.shields.io/badge/Tests-44%20passed-brightgreen?style=flat)](./servicios/)
[![Endpoints](https://img.shields.io/badge/API%20REST-8%20endpoints-blue?style=flat)](./api/)
[![License](https://img.shields.io/badge/License-Academic-yellow?style=flat)](#-licencia)

---

## 📋 Descripción

Sistema integral de comercio electrónico en Go desarrollado con POO estricta, API REST nativa, interfaz web Bootstrap, persistencia JSON y procesamiento concurrente de compras mediante Worker Pool (Goroutines y Channels).

El software opera en **dos modos simultáneos**:
- 🖥️ **Consola interactiva** con menús y navegación por terminal.
- 🌐 **Servidor API REST** en `http://localhost:8080` con 8 endpoints JSON.

---

## 👤 Autor

| Nombre |
|---|
| **Marlon Rivera** |

---

## 🏗️ Arquitectura del Software

```text
sistema-ecommerce-go/
│
├── main.go                         # Punto de entrada (orquestador)
├── go.mod                          # Módulo Go
│
├── modelos/                        # Capa de Dominio (POO + Encapsulación)
│   ├── producto.go                 # Struct Producto (campos privados + getters/setters)
│   ├── cliente.go                  # Struct Cliente (interfaz Mostrable)
│   ├── carrito.go                  # Struct Carrito + ElementoCarrito
│   └── pedido.go                   # Struct Pedido (inmutable, solo getters)
│
├── servicios/                      # Capa de Lógica de Negocio
│   ├── productos.go                # CRUD de productos
│   ├── clientes.go                 # CRUD de clientes + interfaz Mostrable
│   ├── inventario.go               # Consultas de inventario
│   ├── carrito.go                  # Gestión del carrito + Checkout transaccional
│   ├── pedidos.go                  # Historial y búsqueda de pedidos
│   ├── productos_test.go           # 15 pruebas unitarias
│   └── pedidos_test.go             # 14 pruebas unitarias
│
├── utilidades/                     # Capa de Soporte Transversal
│   ├── entradas.go                 # Lectura segura por consola (bufio.Reader)
│   ├── archivos.go                 # Persistencia en archivos JSON
│   └── menu.go                     # Interfaces de menús en consola
│
├── api/                            # Capa Web (Unidad 4)
│   ├── servidor.go                 # Configuración del servidor HTTP y rutas
│   ├── handlers.go                 # 8 Endpoints REST con JSON
│   ├── worker.go                   # Worker Pool (Goroutines + Channels + Mutex)
│   └── handlers_test.go            # 16 pruebas de integración HTTP
│
├── public/                         # Interfaz Gráfica Web
│   └── index.html                  # Panel interactivo en Bootstrap 5 (Tema Claro)
│
└── datos/                          # Persistencia en disco
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
git clone [https://github.com/Jarv1sYT/sistema-ecommerce-go.git](https://github.com/Jarv1sYT/sistema-ecommerce-go.git)
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

### Tabla de Endpoints

| # | Método | Ruta | Descripción | Payload JSON (Entrada) | Código HTTP |
|---|---|---|---|---|---|
| 1 | `GET` | `/api/productos` | Listar catálogo completo | — | `200 OK` |
| 2 | `GET` | `/api/productos/{id}` | Consultar detalle y stock de un producto | — | `200 OK` / `404 Not Found` |
| 3 | `POST` | `/api/productos` | Registrar nuevo producto | `{"codigo":"P001", "nombre":"Laptop HP", "precio":899.99, "cantidad":10}` | `201 Created` / `409 Conflict` |
| 4 | `GET` | `/api/clientes` | Listar clientes registrados | — | `200 OK` |
| 5 | `POST` | `/api/clientes` | Registrar nuevo cliente | `{"identificacion":"1234567890", "nombre":"Juan Pérez", "correo":"juan@email.com"}` | `201 Created` / `409 Conflict` |
| 6 | `POST` | `/api/carrito/agregar` | Añadir producto al carrito | `{"codigo_producto":"P001", "cantidad":2}` | `200 OK` / `404 Not Found` |
| 7 | `POST` | `/api/pedidos/checkout` | Confirmar compra (Checkout) | `{"identificacion_cliente":"1234567890"}` | `201 Created` / `400 Bad Request` |
| 8 | `GET` | `/api/pedidos` | Listar historial de pedidos | — (opcional: `?cliente=1234567890`) | `200 OK` |

### Ejemplos con curl

```bash
# Registrar producto
curl -X POST http://localhost:8080/api/productos \
  -H "Content-Type: application/json" \
  -d '{"codigo":"P001","nombre":"Laptop HP","precio":899.99,"cantidad":10}'

# Registrar cliente
curl -X POST http://localhost:8080/api/clientes \
  -H "Content-Type: application/json" \
  -d '{"identificacion":"1234567890","nombre":"Juan Perez","correo":"juan@email.com"}'

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

```text
Petición HTTP ──► CheckoutHandler ──► [Canal de Solicitudes] ──► Worker Pool (3 Goroutines)
                                                                       │
                                                                 sync.RWMutex
                                                              (Protege inventario)
                                                                       │
                                                              [Canal de Respuestas]
                                                                       │
                                                              Respuesta HTTP JSON
```

### Conceptos Implementados

| Concepto | Archivo | Descripción |
|---|---|---|
| **Goroutines** | `api/worker.go` | 3 funciones `worker()` ejecutándose concurrentemente con `go p.worker(i)` |
| **Channels** | `api/worker.go` | `chan SolicitudPedido` para enviar pedidos al pool y `chan ResultadoPedido` para recibir respuestas |
| **sync.RWMutex** | `api/servidor.go`, `api/worker.go` | `RLock()/RUnlock()` para lecturas concurrentes, `Lock()/Unlock()` para escrituras exclusivas al descontar stock |
| **Simulación asíncrona** | `api/worker.go` | `time.Sleep(500ms)` simula la latencia de una pasarela de pagos externa |

---

## ✅ Testing Automatizado (44 Pruebas)

| Paquete | Archivo | Tests | Tipo |
|---|---|---|---|
| `servicios` | `productos_test.go` | 15 | Unitarias (validaciones, búsqueda, stock) |
| `servicios` | `pedidos_test.go` | 14 | Unitarias (pedidos, carrito, checkout, clientes) |
| `api` | `handlers_test.go` | 16 | Integración HTTP con `httptest.NewRecorder()` |
| **Total** | | **44** | **100% PASS** |

```text
ok   sistema-ecommerce-go/api         0.948s    ← 16 tests PASS
ok   sistema-ecommerce-go/servicios   0.258s    ← 28 tests PASS
```

---

## 📁 Persistencia de Datos

Los datos se almacenan en archivos JSON dentro de la carpeta `datos/`:

| Archivo | Contenido |
|---|---|
| `datos/productos.json` | Inventario de productos con stock actualizado |
| `datos/clientes.json` | Clientes registrados en el sistema |
| `datos/pedidos.json` | Historial completo de compras realizadas |

La carpeta `datos/` se crea automáticamente al ejecutar la aplicación por primera vez.

---

## 📄 Licencia

Proyecto académico desarrollado como parte del proyecto integrador de la asignatura Programación Orientada a Objetos, desarrollado por Marlon Rivera.
