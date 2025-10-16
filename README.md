
# Carrito 

Este proyecto implementa un **Carrito de Compras** utilizando **Go (Golang)** junto con SQL y recursos estáticos.  
El objetivo es gestionar productos, usuarios y operaciones de compra, integrando frontend y backend de forma sencilla.

---

## Estructura del Proyecto

```
Carrito/
│── db/               # Configuración de base de datos
│   ├── queries/      # Consultas SQL
│   ├── schema/       # Esquemas de la base de datos
│   └── sqlc/         # Código generado por SQLC
│
│── js/               # Archivos JavaScript para frontend
│── static/           # Archivos estáticos (CSS, imágenes, etc.)
│── .gitignore        # Archivos a ignorar por Git
│── go.mod            # Dependencias del proyecto Go
│── index.html        # Página principal del frontend
│── main.go           # Punto de entrada del servidor Go
│── README.md         # Documentación del proyecto
│── sqlc.yaml         # Configuración de SQLC
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

- [Go 1.20+](https://go.dev/dl/)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
- Navegador web actualizado

---

Tomas Ilari, Martino Masson, Juan Abraham

Repositorio: [GitHub - Carrito](https://github.com/IlariTomas/Carrito)

Para ejecutar las pruebas de hurl usar el siguiente comando: hurl --test requests.hurl