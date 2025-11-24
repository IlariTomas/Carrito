
# Carrito 

Este proyecto implementa un **Carrito de Compras** utilizando **Go (Golang)** junto con SQL y recursos estáticos.  
El objetivo es gestionar productos, usuarios y operaciones de compra, integrando frontend y backend de forma sencilla.

---

## Estructura del Proyecto

```
Carrito/
├── api/                # Lógica principal del servidor Go
│
├── db/                 # Configuración y acceso a la base de datos
│ ├── schema/           # Esquemas SQL (tablas, constraints, etc.)
│ │ └── schema.sql
│ ├── queries/          # Consultas SQL definidas para sqlc
│ │ └── queries.sql
│ └── sqlc/             # Código Go generado automáticamente por sqlc
│   ├── db.go
│   ├── models.go
│   └── queries.sql.go
│
│── handle/             # Código Go encargado de manejar operaciones CRUD
│ ├── auth.go
│ ├── carrito.go
| ├── product.go
│ └── ventas.go
|
│── views/             # Plantillas (templates) usadas por el servidor Go│ 
| ├── carrito.templ
| ├── ventas_view.templ
| ├── components.templ
| ├── layout_productos.templ
│ ├── layout_user.templ
| ├── layout.templ
| ├── login_view.templ
│ ├── new_producto.templ
│ ├── new_usuario.templ
│ ├── productos_view.templ
│ └── register_view.templ
|
├── static/             # Archivos estáticos (CSS, imágenes, etc.)
├── sqlc.yaml           # Configuración de sqlc
├── main.go         
├── docker-compose.yml  # Configuración para levantar la app con Docker
├── go.mod              # Módulo Go
├── go.sum              # Dependencias verificadas de Go
├── about.html          # Página About
├── README.md           # Documentación del proyecto
├── .dockerignore       # Archivos y carpetas a ignorar por Docker
└── .gitignore          # Archivos y carpetas a ignorar por Git
```

---

##  Ejecución del Proyecto

1. **Clonar el repositorio:**
   git clone https://github.com/IlariTomas/Carrito
   cd Carrito
   

2. **Ejecutar el servidor:**
   make full-reset -- correr el servidor creando archivos templ y sqlc  
   make down       -- detiene los contenedores  
   - En caso de ser la primera ejecucion ejecutar el comando make setup para instalar templ y sqlc

3. **Abrir en el navegador:**  
   Acceder a [http://localhost:8080](http://localhost:8080)

---

##  Dominio de la Aplicación

El sistema desarrollado corresponde a un **Carrito de Compras**, cuyo dominio incluye:  

- **Productos:** artículos disponibles para la compra.  
- **Usuarios:** clientes que interactúan con el sistema.  
- **Carrito:** donde los usuarios agregan productos antes de confirmar la compra.  
- **Persistencia en Base de Datos:** utilizando **sqlc** para generar código Go desde SQL.  
- **Frontend simple:** con HTML, CSS , integrando HTMX para la interacción con el usuario.  
- **Backend en Go:** que gestiona peticiones HTTP y conexión con la base de datos.

---

## Aclaracion de implementacion: Entorno de Ejecución con Docker

Para asegurar la portabilidad y consistencia del proyecto, todo el entorno está gestionado por Docker Compose, el cual levanta los siguientes servicios:

api: Un contenedor con Go que compila y ejecuta la API REST.

db: Un contenedor con Postgres que sirve como base de datos.

Decidimos implentarlo de esta manera ya que garantiza que todos los desarrolladores utilicen exactamente las mismas versiones y configuraciones de software, eliminando la necesidad de instalar dependencias globalmente.

----

##  Requisitos

Instalar docker compose
Aclaracion: docker compose o docker-compose (CON GUION) segun la version de docker Compose instalada. 
          - en caso de tener la version CON GUION modificar en el makefile para poder ejecutar.

Tomas Ilari, Martino Masson, Juan Abraham

Repositorio: [GitHub - Carrito](https://github.com/IlariTomas/Carrito)
