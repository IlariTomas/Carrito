package handle

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	sqlc "carrito.com/db/sqlc"
	"carrito.com/views"
)

// CartHandler maneja las rutas para GET, DELETE en /carrito/{id}
func CartHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCartHandler(queries)(w, r) // GET /carrito/{id}
		case http.MethodDelete:
			deleteCartHandler(queries)(w, r) // DELETE /carrito/{id}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// getCartHandler obtiene los items de un carrito por su ID
func getCartHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraer ID del carrito desde la URL
		idStr := r.URL.Path[len("/carrito/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de carrito inválido", http.StatusBadRequest)
			return
		}

		// Obtener todos los items del carrito
		carritoItems, err := queries.GetCartItems(r.Context(), int32(id))
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Error al obtener productos: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Renderizar la vista templ
		views.CarritoList(carritoItems).Render(r.Context(), w)
	}
}

func deleteCartHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/carrito/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de carrito inválido", http.StatusBadRequest)
			return
		}

		err = queries.RemoveCart(r.Context(), int32(id))
		if err != nil {
			http.Error(w, "Error al eliminar carrito: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func CartItemHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addCartHandler(queries)(w, r) // POST /carrito/items
		case http.MethodPut:
			updateItemHandler(queries)(w, r) // PUT /carrito/items/{id}
		case http.MethodDelete:
			deleteCartHandler(queries)(w, r) // DELETE /carrito/items/{id}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func addCartHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error leyendo formulario", http.StatusBadRequest)
			return
		}

		// Obtener valores del formulario
		nombre := r.FormValue("nombre_producto")
		descripcion := r.FormValue("descripcion")
		precio := r.FormValue("precio")
		stockStr := r.FormValue("stock")
		categoria := r.FormValue("categoria")
		imagen := r.FormValue("imagen")

		// Validación básica
		if nombre == "" || precio == "" {
			http.Error(w, "Nombre y Precio son requeridos", http.StatusBadRequest)
			return
		}

		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			http.Error(w, "Stock inválido", http.StatusBadRequest)
			return
		}

		// Crear parámetros para sqlc
		req := sqlc.CreateProdParams{
			NombreProducto: nombre,
			Descripcion:    descripcion,
			Precio:         precio,
			Stock:          int32(stock),
			Categoria:      categoria,
			Imagen:         imagen,
		}

		// Crear producto en DB
		_, err = queries.CreateProd(r.Context(), req)
		if err != nil {
			http.Error(w, "Error al crear producto: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Respuesta JSON
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Producto creado correctamente"))
	}
}

func updateItemHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraer ID del item desde la URL
		idStr := r.URL.Path[len("/carrito/items/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID del item inválido", http.StatusBadRequest)
			return
		}

		// Parsear JSON del body
		var req sqlc.UpdateCartItemParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// Asignar el ID de la URL
		req.IDItem = int32(id)

		// Ejecutar la actualización en la base de datos
		if err := queries.UpdateCartItem(r.Context(), req); err != nil {
			http.Error(w, "Error al actualizar item: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Respuesta exitosa
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}
}
