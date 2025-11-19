-- name: CreateProd :one
INSERT INTO producto (nombre_producto, descripcion, precio, stock, categoria, imagen) VALUES ($1,$2, $3, $4, $5, $6) RETURNING *;

-- name: CreateUser :one
INSERT INTO usuario (nombre_usuario, email) VALUES ($1, $2) RETURNING id_usuario, nombre_usuario, email;

-- name: CreateVenta :one
INSERT INTO venta (id_producto,id_usuario, cantidad, total, fecha) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetProd :one
SELECT * FROM producto WHERE id_producto = $1;

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

-- name: ListVentasUsuario :many
SELECT * FROM venta WHERE id_usuario = $1;

-- name: ListVentas :many
SELECT * FROM venta ORDER BY fecha;

-- name: UpdateProducto :exec
UPDATE producto SET nombre_producto = $2, descripcion = $3, stock = $4, precio = $5, categoria = $6, imagen = $7 WHERE id_producto = $1;

-- name: UpdateProductoPrecio :exec
UPDATE producto SET precio = $2 WHERE id_producto = $1;

-- name: UpdateProductoStock :exec
UPDATE producto SET stock = $2 WHERE id_producto = $1;

-- name: UpdateUser :exec
UPDATE usuario SET nombre_usuario = $2, email = $3 WHERE id_usuario = $1;

-- name: UpdateVenta :exec
UPDATE venta SET cantidad = $2, total = $3, fecha = $4 WHERE id_venta = $1;

-- name: DeleteProd :exec
DELETE FROM producto WHERE id_producto = $1;

-- name: DeleteUser :exec
DELETE FROM usuario WHERE id_usuario = $1;

-- name: DeleteVenta :exec
DELETE FROM venta WHERE id_venta = $1;

-- name: ListProductsByPriceAsc :many
SELECT id_producto, nombre_producto, descripcion, precio, stock, categoria, imagen FROM producto ORDER BY precio ASC;

-- name: ListProductsByPriceDesc :many
SELECT id_producto, nombre_producto, descripcion, precio, stock, categoria, imagen FROM producto ORDER BY precio DESC;

-- name: AddToCart :one
INSERT INTO carrito (id_usuario, id_producto, cantidad) VALUES ($1, $2, $3) RETURNING *;

-- name: RemoveCartItems :exec
DELETE FROM carrito WHERE id_item = $1;

-- name: RemoveCart :exec
DELETE FROM carrito WHERE id_usuario = $1;

-- name: UpdateCartItem :exec
UPDATE carrito SET cantidad = $2 WHERE id_item = $1;

-- name: GetCartItems :many
SELECT c.*, p.nombre_producto, p.precio FROM carrito c JOIN producto p ON c.id_producto = p.id_producto WHERE c.id_usuario = $1;

-- name: GetCartItemByUserAndProduct :one
SELECT * FROM carrito WHERE id_usuario = $1 AND id_producto = $2;
