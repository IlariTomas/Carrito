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

# === Listar todos los usuarios ===
GET {{host}}/users
HTTP 200
[Asserts]
jsonpath "$[*].nombre_usuario" exists

# === Obtener un usuario por ID ===
GET {{host}}/user/1
HTTP 200
[Asserts]
jsonpath "$.nombre_usuario" exists
jsonpath "$.email" exists

# === Actualizar un usuario ===
PUT {{host}}/user/1
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
DELETE {{host}}/user/1
HTTP 204

# === Intentar obtener un usuario eliminado ===
GET {{host}}/user/1
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

# ====================================
# CHEQUEOS PARA PRODUCTOS
# ====================================

# ====================================
# CHEQUEOS PARA PRODUCTOS
# ====================================

# === Crear un Producto ===
# ARREGLO: El cuerpo del POST debe estar en 'snake_case'

POST {{host}}/products
Content-Type: application/json

{
  "nombre_producto": "Laptop Gamer X",
  "descripcion": "La laptop más potente del mercado.",
  "stock": 10,
  "precio": "1200.00",
  "categoria": "Electrónica"
}

HTTP 201
[Asserts]
jsonpath "$.nombre_producto" == "Laptop Gamer X"
jsonpath "$.descripcion" == "La laptop más potente del mercado."
jsonpath "$.stock" == 10
jsonpath "$.precio" == "1200.00"
jsonpath "$.categoria" == "Electrónica"


# === Listar todos los Productos ===
GET {{host}}/products
HTTP 200
[Asserts]
jsonpath "$[*].nombre_producto" exists
jsonpath "$[*].descripcion" exists
jsonpath "$[*].stock" exists
jsonpath "$[*].precio" exists
jsonpath "$[*].categoria" exists


# === Obtener un Producto por ID (Ejemplo con ID 1) ===
GET {{host}}/product/1
HTTP 200
[Asserts]
jsonpath "$.nombre_producto" exists
jsonpath "$.descripcion" exists
jsonpath "$.stock" exists
jsonpath "$.precio" exists
jsonpath "$.categoria" exists


# === Actualizar un Producto (Ejemplo con ID 1) ===
PUT {{host}}/product/1
Content-Type: application/json

{
  "nombre_producto": "Laptop Pro Max",
  "descripcion": "La mejor laptop para desarrollo.",
  "stock": 6,
  "precio": "1500.00",
  "categoria": "Tecnología",
  "imagen": "" 
}

HTTP 200

[Asserts]
jsonpath "$.nombre_producto" == "Laptop Pro Max"
jsonpath "$.descripcion" == "La mejor laptop para desarrollo."
jsonpath "$.stock" == 6
jsonpath "$.precio" == "1500.00"
jsonpath "$.categoria" == "Tecnología"

# === Eliminar un Producto (Ejemplo con ID 1) ===
DELETE {{host}}/product/1
HTTP 204


# === Intentar obtener un Producto eliminado ===
GET {{host}}/product/1
HTTP 404

# === Crear un Producto para venta ===
# ARREGLO: El cuerpo del POST debe estar en 'snake_case'.

POST {{host}}/products
Content-Type: application/json

{
  "nombre_producto": "Mouse Inalámbrico",
  "descripcion": "Un mouse inalámbrico ergonómico.",
  "stock": 25,
  "precio": "50.00",
  "categoria": "Accesorios"
}

HTTP 201
[Asserts]
jsonpath "$.nombre_producto" == "Mouse Inalámbrico"
jsonpath "$.descripcion" == "Un mouse inalámbrico ergonómico."
jsonpath "$.stock" == 25
jsonpath "$.precio" == "50.00"
jsonpath "$.categoria" == "Accesorios"

# ====================================
# CHEQUEOS PARA VENTAS
# ====================================

# === Crear una Venta ===
POST {{host}}/sales
Content-Type: application/json

{
  "id_producto": 2,
  "id_usuario": 2,
  "cantidad": 5,
  "total": "250.00",
  "fecha": "2025-10-24T18:00:00Z"
}

HTTP 201
[Asserts]
jsonpath "$.id_producto" == 2
jsonpath "$.id_usuario" == 2
jsonpath "$.cantidad" == 5
jsonpath "$.total" == "250.00"
#jsonpath "$.fecha" == "2025-10-24T18:00:00Z"


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
GET {{host}}/sale/1
HTTP 200
[Asserts]
jsonpath "$.id_producto" == 2
jsonpath "$.id_usuario" == 2
jsonpath "$.cantidad" exists
jsonpath "$.total" exists
#jsonpath "$.fecha" == "2025-10-24T18:00:00Z"


# === Actualizar una Venta (Ejemplo con ID 1) ===
# Nota: UpdateVentaParams solo actualiza cantidad, total, y fecha.
PUT {{host}}/sale/1
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
DELETE {{host}}/sale/1
HTTP 204


# === Intentar obtener una Venta eliminada ===
GET {{host}}/sale/1
HTTP 404

# === Eliminar un Producto ===
DELETE {{host}}/product/2
HTTP 204

# === Eliminar un Usuario ===
DELETE {{host}}/user/2
HTTP 204