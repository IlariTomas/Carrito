# ====================================
# CHEQUEOS PARA USUARIOS
# ====================================
# === Crear un usuario ===
POST {{host}}/users
Content-Type: application/json

{
  "nombre_usuario": "John Doe",
  "email": "john.doe@example.com"
}

HTTP 201
[Asserts]
jsonpath "$.nombre_usuario" == "John Doe"
jsonpath "$.email" == "john.doe@example.com"
[Captures]
userId: jsonpath "$.id_usuario"

# === Listar todos los usuarios ===
GET {{host}}/users
HTTP 200
[Asserts]
jsonpath "$[*].nombre_usuario" exists

# === Obtener un usuario por ID ===
GET {{host}}/user/{{userId}}
HTTP 200
[Asserts]
jsonpath "$.nombre_usuario" exists
jsonpath "$.email" exists

# === Actualizar un usuario ===
PUT {{host}}/user/{{userId}}
Content-Type: application/json

{
  "nombre_usuario": "Johnny Doe",
  "email": "johnny.doe@example.com"
}

HTTP 200
[Asserts]
jsonpath "$.nombre_usuario" == "Johnny Doe"
jsonpath "$.email" == "johnny.doe@example.com"

# === Eliminar un usuario ===
DELETE {{host}}/user/{{userId}}
HTTP 204

# === Intentar obtener un usuario eliminado ===
GET {{host}}/user/{{userId}}
HTTP 404

# === Crear un usuario para ventas ===
POST {{host}}/users
Content-Type: application/json

{
  "nombre_usuario": "Usuario Venta",
  "email": "usuario.venta@example.com"
}

HTTP 201
[Asserts]
jsonpath "$.nombre_usuario" == "Usuario Venta"
jsonpath "$.email" == "usuario.venta@example.com"
[Captures]
secondUserId: jsonpath "$.id_usuario"

# ====================================
# CHEQUEOS PARA PRODUCTOS
# ====================================

# === Crear un Producto ===
POST {{host}}/products
Content-Type: application/json

{
  "nombre_producto": "Laptop Gamer X",
  "descripcion": "La laptop más potente del mercado.",
  "stock": 10,
  "precio": "1200.00",
  "categoria": "Electrónica",
  "imagen": "https://www.crucial.mx/content/dam/crucial/articles/for-pc-builders/new025-how-to-upgrade-your-pc/modern-gaming-pc.jpg.transform/medium-jpg/img.jpg"
}

HTTP 201
[Asserts]
jsonpath "$.nombre_producto" == "Laptop Gamer X"
jsonpath "$.descripcion" == "La laptop más potente del mercado."
jsonpath "$.stock" == 10
jsonpath "$.precio" == "1200.00"
jsonpath "$.categoria" == "Electrónica"
jsonpath "$.imagen" == "https://www.crucial.mx/content/dam/crucial/articles/for-pc-builders/new025-how-to-upgrade-your-pc/modern-gaming-pc.jpg.transform/medium-jpg/img.jpg"
[Captures]
productId: jsonpath "$.id_producto"

# === Listar todos los Productos ===
GET {{host}}/products
HTTP 200
[Asserts]
jsonpath "$[*].nombre_producto" exists
jsonpath "$[*].descripcion" exists
jsonpath "$[*].stock" exists
jsonpath "$[*].precio" exists
jsonpath "$[*].categoria" exists
jsonpath "$[*].imagen" exists


# === Obtener un Producto por ID (Ejemplo con ID 1) ===
GET {{host}}/product/{{productId}}
HTTP 200
[Asserts]
jsonpath "$.nombre_producto" exists
jsonpath "$.descripcion" exists
jsonpath "$.stock" exists
jsonpath "$.precio" exists
jsonpath "$.categoria" exists
jsonpath "$.imagen" exists 


# === Actualizar un Producto (Ejemplo con ID 1) ===
PUT {{host}}/product/{{productId}}
Content-Type: application/json

{
  "nombre_producto": "Laptop Pro Max",
  "descripcion": "La mejor laptop para desarrollo.",
  "stock": 6,
  "precio": "1500.00",
  "categoria": "Tecnología",
  "imagen": "https://www.crucial.mx/content/dam/crucial/articles/for-pc-builders/new025-how-to-upgrade-your-pc/modern-gaming-pc.jpg.transform/medium-jpg/img.jpg"
}


HTTP 200
[Asserts]
jsonpath "$.nombre_producto" == "Laptop Pro Max"
jsonpath "$.descripcion" == "La mejor laptop para desarrollo."
jsonpath "$.stock" == 6
jsonpath "$.precio" == "1500.00"
jsonpath "$.categoria" == "Tecnología"
jsonpath "$.imagen" == "https://www.crucial.mx/content/dam/crucial/articles/for-pc-builders/new025-how-to-upgrade-your-pc/modern-gaming-pc.jpg.transform/medium-jpg/img.jpg"


# === Eliminar un Producto (Ejemplo con ID 1) ===
DELETE {{host}}/product/{{productId}}
HTTP 204


# === Intentar obtener un Producto eliminado ===
GET {{host}}/product/{{productId}}
HTTP 404

# === Crear un Producto para venta ===
POST {{host}}/products
Content-Type: application/json

{
  "nombre_producto": "Mouse Inalámbrico",
  "descripcion": "Un mouse inalámbrico ergonómico.",
  "stock": 25,
  "precio": "50.00",
  "categoria": "Accesorios",
  "imagen": "https://www.crucial.mx/content/dam/crucial/articles/for-pc-builders/new025-how-to-upgrade-your-pc/modern-gaming-pc.jpg.transform/medium-jpg/img.jpg"
} 

HTTP 201
[Asserts]
jsonpath "$.nombre_producto" == "Mouse Inalámbrico"
jsonpath "$.descripcion" == "Un mouse inalámbrico ergonómico."
jsonpath "$.stock" == 25
jsonpath "$.precio" == "50.00"
jsonpath "$.categoria" == "Accesorios"
jsonpath "$.imagen" == "https://www.crucial.mx/content/dam/crucial/articles/for-pc-builders/new025-how-to-upgrade-your-pc/modern-gaming-pc.jpg.transform/medium-jpg/img.jpg"
[Captures]
secondProductId: jsonpath "$.id_producto"

# ====================================
# CHEQUEOS PARA VENTAS
# ====================================

# === Crear una Venta ===
POST {{host}}/sales
Content-Type: application/json

{
  "id_producto": {{secondProductId}},
  "id_usuario": {{secondUserId}},
  "cantidad": 5,
  "total": "250.00",
  "fecha": "2025-10-24T18:00:00Z"
}

HTTP 201
[Asserts]
jsonpath "$.id_producto" == {{secondProductId}}
jsonpath "$.id_usuario" == {{secondUserId}}
jsonpath "$.cantidad" == 5
jsonpath "$.total" == "250.00"
#jsonpath "$.fecha" == "2025-10-24T18:00:00Z"
[Captures]
saleId: jsonpath "$.id_venta"


# === Listar todas las Ventas ===
GET {{host}}/sales
HTTP 200
[Asserts]
jsonpath "$[*].id_producto" exists
jsonpath "$[*].id_usuario" exists
jsonpath "$[*].id_venta" exists
jsonpath "$[*].cantidad" exists
jsonpath "$[*].total" exists
#jsonpath "$[*].fecha" exists


# === Obtener una Venta por ID (Ejemplo con ID 1) ===
GET {{host}}/sale/{{saleId}}
HTTP 200
[Asserts]
jsonpath "$.id_producto" == {{secondProductId}}
jsonpath "$.id_usuario" == {{secondUserId}}
jsonpath "$.cantidad" exists
jsonpath "$.total" exists
#jsonpath "$.fecha" == "2025-10-24T18:00:00Z"


# === Actualizar una Venta (Ejemplo con ID 1) ===
PUT {{host}}/sale/{{saleId}}
Content-Type: application/json

{
  "cantidad": 10,
  "total": "500.00",
  "fecha": "2025-10-25T09:00:00Z"
}

HTTP 200
[Asserts]
jsonpath "$.cantidad" == 10
jsonpath "$.total" == "500.00"
#jsonpath "$.fecha" == "2025-10-25T09:00:00Z"


# === Eliminar una Venta (Ejemplo con ID 1) ===
DELETE {{host}}/sale/{{saleId}}
HTTP 204


# === Intentar obtener una Venta eliminada ===
GET {{host}}/sale/{{saleId}}
HTTP 404

# === Eliminar un Producto ===
DELETE {{host}}/product/{{secondProductId}}
HTTP 204

# === Eliminar un Usuario ===
DELETE {{host}}/user/{{secondUserId}}
HTTP 204