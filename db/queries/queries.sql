-- name: CreateProd :exec
INSERT INTO producto (nombre_producto, descripcion, precio, categoria) VALUES ($1,$2, $3, $4);

-- name: CreateUser :one
INSERT INTO usuario (nombre_usuario, email) VALUES ($1, $2) RETURNING id_usuario, nombre_usuario, email;

-- name: CreateVenta :exec
INSERT INTO venta (id_producto, id_venta, cantidad, total, fecha) VALUES ($1,$2, $3, $4, $5);

-- name: GetProd :one
SELECT nombre_producto, descripcion, precio, categoria FROM producto WHERE id_producto = $1;

-- name: GetVenta :one
SELECT * FROM venta WHERE id_venta = $1;

-- name: GetVenta_usuario :one
SELECT * FROM venta WHERE id_usuario = $1;

-- name: GetUser :one
SELECT nombre_usuario, email FROM usuario WHERE id_usuario = $1;

-- name: ListProd :many
SELECT * FROM producto ORDER BY nombre_producto;

-- name: ListUsers :many
SELECT * FROM usuario ORDER BY nombre_usuario;

-- name: ListVentas :many
SELECT * FROM venta WHERE id_usuario = $1;

-- name: UpdateProductoPrecio :exec
UPDATE producto SET precio = $2 WHERE id_producto = $1;

-- name: UpdateProductoStock :exec
UPDATE producto SET stock = $2 WHERE id_producto = $1;

-- name: UpdateUser :exec
UPDATE usuario SET nombre_usuario = $2, email = $3 WHERE id_usuario = $1;

-- name: DeleteProd :exec
DELETE FROM producto WHERE id_producto = $1;

-- name: DeleteUser :exec
DELETE FROM usuario WHERE id_usuario = $1;
