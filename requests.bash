curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"nombre":"John Doe","email":"john.doe@example.com"}'
    #Crea un usuario y devuelve su ID, nombre y email en formato JSON

curl http://localhost:8080/users
#Lista todos los usuarios en formato JSON

curl http://localhost:8080/users/1
#Obtiene un usuario por su ID en formato JSON

curl -X PUT http://localhost:8080/users/1 \
     -H "Content-Type: application/json" \
     -d '{"nombre":"Johnny Doe","email":"johnny.doe@example.com"}'
#Actualiza un usuario por su ID y devuelve el usuario actualizado en formato JSON

curl -X DELETE http://localhost:8080/users/1
#Elimina un usuario por su ID

curl http://localhost:8080/users/1
#Intenta obtener un usuario eliminado, debería devolver un error 404