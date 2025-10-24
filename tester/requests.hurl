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
GET {{host}}/users/1
HTTP 200
[Asserts]
#jsonpath "$.id" == 1
jsonpath "$.nombre_usuario" exists
jsonpath "$.email" exists

# === Actualizar un usuario ===
PUT {{host}}/users/1
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
DELETE {{host}}/users/1
HTTP 204

# === Intentar obtener un usuario eliminado ===
GET {{host}}/users/1
HTTP 404

# ====================================
# CHEQUEOS PARA PRODUCTOS
# ====================================

# === Crear un Producto ===
# Nota: La estructura CreateProdParams incluye nombre_producto, descripcion, precio, categoria
POST {{host}}/products
Content-Type: application/json

{
  "nombre_producto": "Laptop Gamer X",
  "descripcion": "La laptop más potente del mercado.",
  "precio": "1200.00",
  "categoria": "Electrónica"
}

HTTP 201
[Asserts]
jsonpath "$.nombre_producto" == "Laptop Gamer X"
jsonpath "$.precio" == "1200.00"


# === Listar todos los Productos ===
GET {{host}}/products
HTTP 200
[Asserts]
jsonpath "$[*].nombre_producto" exists
jsonpath "$[*].precio" exists


# === Obtener un Producto por ID (Ejemplo con ID 1) ===
GET {{host}}/products/1
HTTP 200
[Asserts]
jsonpath "$.nombre_producto" exists
jsonpath "$.precio" exists


# === Actualizar un Producto (Ejemplo con ID 1) ===
# Nota: Usa UpdateProductoParams
PUT {{host}}/products/1
Content-Type: application/json

{
  "nombre_producto": "Laptop Pro Max",
  "descripcion": "La mejor laptop para desarrollo.",
  "precio": "1500.00",
  "categoria": "Tecnología"
}

HTTP 200
[Asserts]
jsonpath "$.nombre_producto" == "Laptop Pro Max"
jsonpath "$.precio" == "1500.00"
jsonpath "$.categoria" == "Tecnología"


# === Eliminar un Producto (Ejemplo con ID 1) ===
DELETE {{host}}/products/1
HTTP 204


# === Intentar obtener un Producto eliminado ===
GET {{host}}/products/1
HTTP 404

# ====================================
# CHEQUEOS PARA VENTAS
# ====================================

# === Crear una Venta ===
POST {{host}}/sales
Content-Type: application/json

{
  "id_producto": 2,      
  "id_venta": 99,       
  "cantidad": 5,
  "total": "250.00",
  "fecha": "2025-10-24T18:00:00Z"
}

HTTP 201
[Asserts]
jsonpath "$.id_producto" == 2
jsonpath "$.cantidad" == 5
jsonpath "$.total" == "250.00"


# === Listar todas las Ventas ===
GET {{host}}/sales
HTTP 200
[Asserts]
jsonpath "$[*].id_producto" exists
jsonpath "$[*].total" exists


# === Obtener una Venta por ID (Ejemplo con ID 99) ===
GET {{host}}/sales/99
HTTP 200
[Asserts]
jsonpath "$.id_venta" == 99
jsonpath "$.cantidad" exists
jsonpath "$.total" exists


# === Actualizar una Venta (Ejemplo con ID 99) ===
# Nota: UpdateVentaParams solo actualiza cantidad, total, y fecha.
PUT {{host}}/sales/99
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


# === Eliminar una Venta (Ejemplo con ID 99) ===
DELETE {{host}}/sales/99
HTTP 204


# === Intentar obtener una Venta eliminada ===
GET {{host}}/sales/99
HTTP 404