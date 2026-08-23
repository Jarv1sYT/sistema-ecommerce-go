# Diagramas de Arquitectura y Clases — Sistema de E-Commerce

> Estos diagramas están escritos en sintaxis [Mermaid](https://mermaid.js.org/) y renderizan automáticamente en GitHub.

---

## 1. Diagrama de Clases POO (Encapsulación)

```mermaid
classDiagram
    direction TB

    class Producto {
        -codigo : string
        -nombre : string
        -precio : float64
        -cantidadDisponible : int
        +NuevoProducto(codigo, nombre, precio, cantidad) Producto
        +GetCodigo() string
        +GetNombre() string
        +GetPrecio() float64
        +GetCantidadDisponible() int
        +SetNombre(nombre) error
        +SetPrecio(precio) error
        +SetCantidadDisponible(cantidad) error
        +TieneStockSuficiente(cantidad) bool
        +DescontarStock(cantidad) error
        +MarshalJSON() bytes
        +UnmarshalJSON(bytes) error
    }

    class Cliente {
        -identificacion : string
        -nombre : string
        -correo : string
        +NuevoCliente(id, nombre, correo) Cliente
        +GetIdentificacion() string
        +GetNombre() string
        +GetCorreo() string
        +SetNombre(nombre) error
        +SetCorreo(correo) error
        +MostrarFicha()
        +MarshalJSON() bytes
        +UnmarshalJSON(bytes) error
    }

    class ElementoCarrito {
        -producto : Producto
        -cantidad : int
        -subtotal : float64
        +NuevoElementoCarrito(producto, cantidad) ElementoCarrito
        +GetProducto() Producto
        +GetCantidad() int
        +GetSubtotal() float64
        +SetCantidad(cantidad) error
    }

    class Carrito {
        -cliente : *Cliente
        -elementos : []ElementoCarrito
        -total : float64
        +NuevoCarrito() Carrito
        +SetCliente(cliente)
        +GetCliente() *Cliente
        +GetElementos() []ElementoCarrito
        +GetTotal() float64
        +EsVacio() bool
        +AgregarElemento(producto, cantidad) error
        +EliminarElemento(codigo) bool
        +Vaciar()
        +CalcularTotal() float64
    }

    class Pedido {
        -codigo : string
        -cliente : Cliente
        -elementos : []ElementoCarrito
        -total : float64
        -fecha : string
        +NuevoPedido(codigo, cliente, elementos, fecha) Pedido
        +GetCodigo() string
        +GetCliente() Cliente
        +GetElementos() []ElementoCarrito
        +GetTotal() float64
        +GetFecha() string
    }

    class Mostrable {
        <<interface>>
        +MostrarFicha()
    }

    Carrito "1" *-- "0..*" ElementoCarrito : contiene
    Carrito "1" o-- "0..1" Cliente : asocia
    ElementoCarrito "1" *-- "1" Producto : referencia
    Pedido "1" *-- "1..*" ElementoCarrito : copia inmutable
    Pedido "1" *-- "1" Cliente : comprador
    Cliente ..|> Mostrable : implementa

    note for Pedido "INMUTABLE: Solo getters\nSin setters después de creación"
    note for Carrito "GetElementos() retorna COPIA\npara proteger encapsulación"
```

---

## 2. Diagrama de Arquitectura por Capas

```mermaid
flowchart TB
    subgraph ENTRADA["🔄 Entrada del Sistema"]
        CONSOLA["🖥️ Consola Interactiva<br/>main.go + utilidades/menu.go"]
        HTTP["🌐 Servidor HTTP :8080<br/>api/servidor.go"]
    end

    subgraph API["📡 Capa API REST (8 Endpoints)"]
        H1["GET /api/productos"]
        H2["GET /api/productos/id"]
        H3["POST /api/productos"]
        H4["GET /api/clientes"]
        H5["POST /api/clientes"]
        H6["POST /api/carrito/agregar"]
        H7["POST /api/pedidos/checkout"]
        H8["GET /api/pedidos"]
    end

    subgraph SERVICIOS["⚙️ Capa de Servicios (Lógica de Negocio)"]
        SP["servicios/productos.go<br/>CRUD Productos"]
        SC["servicios/clientes.go<br/>CRUD Clientes"]
        SK["servicios/carrito.go<br/>Carrito + Checkout"]
        SI["servicios/inventario.go<br/>Consultas Stock"]
        SPD["servicios/pedidos.go<br/>Historial Pedidos"]
    end

    subgraph MODELOS["📦 Capa de Dominio (POO Encapsulado)"]
        MP["Producto"]
        MC["Cliente"]
        MCA["Carrito + ElementoCarrito"]
        MPE["Pedido"]
    end

    subgraph CONCURRENCIA["⚡ Capa de Concurrencia"]
        WP["Worker Pool<br/>3 Goroutines"]
        CH["Channel<br/>chan SolicitudPedido"]
        MU["sync.RWMutex<br/>Protección de Stock"]
    end

    subgraph PERSISTENCIA["💾 Capa de Persistencia"]
        UJ["utilidades/archivos.go<br/>GuardarJSON / LeerJSON"]
        DJ["datos/productos.json"]
        DC["datos/clientes.json"]
        DP["datos/pedidos.json"]
    end

    CONSOLA --> SERVICIOS
    HTTP --> API
    API --> SERVICIOS
    H7 --> CH
    CH --> WP
    WP --> MU
    MU --> SERVICIOS
    SERVICIOS --> MODELOS
    SERVICIOS --> PERSISTENCIA
    UJ --> DJ
    UJ --> DC
    UJ --> DP
```

---

## 3. Diagrama de Flujo de Concurrencia (Checkout)

```mermaid
sequenceDiagram
    participant C as Cliente HTTP
    participant H as CheckoutHandler
    participant CH as Canal Solicitudes
    participant W as Worker Goroutine
    participant MX as sync.RWMutex
    participant INV as Inventario
    participant JSON as datos/*.json

    C->>H: POST /api/pedidos/checkout
    H->>H: Validar carrito no vacío
    H->>H: Buscar cliente (RLock)
    H->>CH: Enviar SolicitudPedido
    Note over CH: El canal transfiere la<br/>solicitud al Worker Pool
    CH->>W: Worker recibe solicitud
    W->>W: Simular validación de pago (500ms)
    W->>MX: Lock() - Acceso exclusivo
    MX->>INV: Verificar stock suficiente
    INV-->>MX: Stock confirmado
    MX->>INV: DescontarStock()
    INV-->>MX: Stock actualizado
    MX->>JSON: GuardarProductos()
    W->>W: Crear Pedido inmutable (NuevoPedido)
    MX->>JSON: GuardarPedidos()
    W->>MX: Unlock() - Liberar acceso
    W-->>H: ResultadoPedido (por canal respuesta)
    H->>H: Vaciar carrito
    H-->>C: 201 Created + Comprobante JSON
```

---

## 4. Diagrama de Paquetes y Dependencias

```mermaid
flowchart LR
    MAIN["main.go"] --> MODELOS["modelos/"]
    MAIN --> SERVICIOS["servicios/"]
    MAIN --> UTILIDADES["utilidades/"]
    MAIN --> API["api/"]

    SERVICIOS --> MODELOS
    SERVICIOS --> UTILIDADES

    API --> MODELOS
    API --> SERVICIOS

    MODELOS -.->|"encoding/json"| STDLIB["Librería Estándar Go"]
    UTILIDADES -.->|"encoding/json<br/>bufio<br/>os"| STDLIB
    API -.->|"net/http<br/>sync<br/>time"| STDLIB

    style MAIN fill:#4CAF50,color:white
    style MODELOS fill:#2196F3,color:white
    style SERVICIOS fill:#FF9800,color:white
    style UTILIDADES fill:#9C27B0,color:white
    style API fill:#F44336,color:white
    style STDLIB fill:#607D8B,color:white
```
