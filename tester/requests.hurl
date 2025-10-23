# === Crear un usuario ===
POST {{host}}/users
Content-Type: application/json

{
  "nombre": "John Doe",
  "email": "john.doe@example.com"
}

HTTP 201
[Asserts]
jsonpath "$.name" == "John Doe"
jsonpath "$.email" == "john.doe@example.com"

# === Listar todos los usuarios ===
GET {{host}}/users
HTTP 200
[Asserts]
jsonpath "$[*].name" exists

# === Obtener un usuario por ID ===
GET {{host}}/users/1
HTTP 200
[Asserts]
#jsonpath "$.id" == 1
jsonpath "$.nombre" exists
jsonpath "$.email" exists

# === Actualizar un usuario ===
PUT {{host}}/users/1
Content-Type: application/json

{
  "nombre": "Johnny Doe",
  "email": "johnny.doe@example.com"
}

HTTP 200
[Asserts]
jsonpath "$.nombre" == "Johnny Doe"
jsonpath "$.email" == "johnny.doe@example.com"

# === Eliminar un usuario ===
DELETE {{host}}/users/1
HTTP 204

# === Intentar obtener un usuario eliminado ===
GET {{host}}/users/1
HTTP 404