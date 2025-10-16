# === Crear un usuario ===
POST http://localhost:8080/users
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
GET http://localhost:8080/users
HTTP 200
[Asserts]
jsonpath "$[*].name" exists

# === Obtener un usuario por ID ===
GET http://localhost:8080/users/25
HTTP 200
[Asserts]
#jsonpath "$.id" == 25
jsonpath "$.nombre" exists
jsonpath "$.email" exists

# === Actualizar un usuario ===
PUT http://localhost:8080/users/25
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
DELETE http://localhost:8080/users/25
HTTP 204

# === Intentar obtener un usuario eliminado ===
GET http://localhost:8080/users/25
HTTP 404
