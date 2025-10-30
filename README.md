
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
│ ├── product.go
│ ├── user.go
│ └── ventas.go
│
├── js/                 # Archivos JavaScript para el frontend
├── static/             # Archivos estáticos (CSS, imágenes, etc.)
├── tester/             # Pruebas y archivos de testing (por ejemplo requests.hurl)
│ ├── Dockerfile
│ └── requests.hurl
│
├── sqlc.yaml           # Configuración de sqlc
├── main.go         
├── docker-compose.yml  # Configuración para levantar la app con Docker
├── go.mod              # Módulo Go
├── go.sum              # Dependencias verificadas de Go
├── index.html          # Página principal del frontend
├── README.md           # Documentación del proyecto
└── .gitignore          # Archivos y carpetas a ignorar por Git
```

---

##  Ejecución del Proyecto

1. **Clonar el repositorio:**
   git clone https://github.com/IlariTomas/Carrito
   cd Carrito
   

2. **Ejecutar el servidor:**
   go run main.go
   

3. **Abrir en el navegador:**  
   Acceder a [http://localhost:8080](http://localhost:8080)

---

##  Dominio de la Aplicación

El sistema desarrollado corresponde a un **Carrito de Compras**, cuyo dominio incluye:  

- **Productos:** artículos disponibles para la compra.  
- **Usuarios:** clientes que interactúan con el sistema.  
- **Carrito:** donde los usuarios agregan productos antes de confirmar la compra.  
- **Persistencia en Base de Datos:** utilizando **sqlc** para generar código Go desde SQL.  
- **Frontend simple:** con HTML, CSS y JS para la interacción con el usuario.  
- **Backend en Go:** que gestiona peticiones HTTP, operaciones CRUD y conexión con la base de datos.

---

##  Requisitos

Instalar docker-compose

Ejecutar el siguiente comando para levantar los contenedores y correr los test Hurl
docker-compose up --build


Tomas Ilari, Martino Masson, Juan Abraham

Repositorio: [GitHub - Carrito](https://github.com/IlariTomas/Carrito)
