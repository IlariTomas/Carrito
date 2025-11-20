package handle

import (
	"database/sql"
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
			getCartHandler(queries)(w, r) // GET /carrito
		case http.MethodDelete:
			deleteCartHandler(queries)(w, r) // DELETE /carrito
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// getCartHandler obtiene los items de un carrito por su ID
func getCartHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "No hay sesión activa", http.StatusUnauthorized)
			return
		}

		id, err := strconv.Atoi(cookie.Value)
		if err != nil {
			http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
			return
		}

		// Obtener items del carrito
		carritoItems, err := queries.GetCartItems(r.Context(), int32(id))
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Error al obtener productos: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		componente := views.CarritoList(carritoItems)
		componente.Render(r.Context(), w)
	}
}

func deleteCartHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, err := r.Cookie("session_token")
		id, err := strconv.Atoi(idStr.Value)
		if err != nil {
			http.Error(w, "ID de carrito inválido", http.StatusBadRequest)
			return
		}

		err = queries.DeleteCart(r.Context(), int32(id))
		if err != nil {
			http.Error(w, "Error al eliminar carrito: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Obtener items del carrito
		carritoItems, err := queries.GetCartItems(r.Context(), int32(id))
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Error al obtener productos: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		componente := views.CarritoList(carritoItems)
		componente.Render(r.Context(), w)

		w.WriteHeader(http.StatusNoContent)
	}
}

// HANDLERS PARA ITEMS DEL CARRITO

func CartItemHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addCartHandler(queries)(w, r) // POST /carrito/items/{id}
		case http.MethodPut:
			updateItemHandler(queries)(w, r) // PUT /carrito/items/{id}
		case http.MethodDelete:
			deleteCartItemsHandler(queries)(w, r) // DELETE /carrito/items/{id}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func addCartHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idUsuarioStr, err := r.Cookie("session_token")
		idUsuario, err := strconv.Atoi(idUsuarioStr.Value)

		idStr := r.URL.Path[len("/carrito/items/"):]
		idProducto, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de producto inválido", http.StatusBadRequest)
			return
		}

		// Intento obtener el item del carrito
		item, err := queries.GetCartItemByUserAndProduct(
			r.Context(),
			sqlc.GetCartItemByUserAndProductParams{
				IDUsuario:  int32(idUsuario),
				IDProducto: int32(idProducto),
			},
		)

		if err == nil {
			// Existe → sumo cantidad
			update := sqlc.UpdateCartItemParams{
				IDItem:   item.IDItem,
				Cantidad: item.Cantidad + 1,
			}

			if err := queries.UpdateCartItem(r.Context(), update); err != nil {
				http.Error(w, "Error al actualizar cantidad", http.StatusInternalServerError)
				return
			}

		} else {
			// No existe → creo nuevo
			req := sqlc.AddToCartParams{
				IDUsuario:  int32(idUsuario),
				IDProducto: int32(idProducto),
				Cantidad:   1,
			}

			if _, err := queries.AddToCart(r.Context(), req); err != nil {
				http.Error(w, "Error al agregar producto", http.StatusInternalServerError)
				return
			}
		}

		// 🔹 Renderizo solo el carrito actualizado
		carritoItems, err := queries.GetCartItems(r.Context(), int32(idUsuario))
		if err != nil {
			http.Error(w, "Error cargando carrito", http.StatusInternalServerError)
			return
		}

		componente := views.CarritoList(carritoItems)
		componente.Render(r.Context(), w)
	}
}

func updateItemHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/carrito/items/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID del item inválido", http.StatusBadRequest)
			return
		}
		var req sqlc.UpdateCartItemParams

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

func deleteCartItemsHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/carrito/items/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de carrito inválido", http.StatusBadRequest)
			return
		}

		err = queries.DeleteProdCarrito(r.Context(), int32(id))
		if err != nil {
			http.Error(w, "Error al eliminar producto del carrito: "+err.Error(), http.StatusInternalServerError)
			return
		}

		cookie, _ := r.Cookie("session_token")
		id_usuario, _ := strconv.Atoi(cookie.Value)

		// Obtener items del carrito
		carritoItems, err := queries.GetCartItems(r.Context(), int32(id_usuario))
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Error al obtener productos: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		componente := views.CarritoList(carritoItems)
		componente.Render(r.Context(), w)

		w.WriteHeader(http.StatusNoContent)
	}
}
