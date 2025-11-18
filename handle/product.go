package handle

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	sqlc "carrito.com/db/sqlc"
	"carrito.com/views"
)

func ProductsHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listProdHandler(queries)(w, r) // GET /products
		case http.MethodPost:
			createProdHandler(queries)(w, r) // POST /products
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// Producto: POST /products
func createProdHandler(queries *sqlc.Queries) http.HandlerFunc {
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

	}
}

// Producto: GET /products
func listProdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		views.ProductView().Render(r.Context(), w)
	}
}

func ProductHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProdHandler(queries)(w, r) // GET /products/{id}
		case http.MethodPut:
			updateProdHandler(queries)(w, r) // PUT /products/{id}
		case http.MethodDelete:
			deleteProdHandler(queries)(w, r) // DELETE /products/{id}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// Producto: GET /products/{id}
func getProdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/product/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de producto inválido", http.StatusBadRequest)
			return
		}

		product, err := queries.GetProd(r.Context(), int32(id))
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Error al obtener producto: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

// Producto: PUT /products/{id} (Usando UpdateProducto)
func updateProdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/product/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de producto inválido", http.StatusBadRequest)
			return
		}

		var req sqlc.UpdateProductoParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// Asignar el ID de la URL al struct de parámetros
		req.IDProducto = int32(id)

		err = queries.UpdateProducto(r.Context(), req)
		if err != nil {
			http.Error(w, "Error al actualizar producto: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Obtener y devolver el producto actualizado
		product, _ := queries.GetProd(r.Context(), int32(id))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

// Producto: DELETE /products/{id}
func deleteProdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/product/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de producto inválido", http.StatusBadRequest)
			return
		}

		err = queries.DeleteProd(r.Context(), int32(id))
		if err != nil {
			http.Error(w, "Error al eliminar producto: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// handler para  /list-products (HTMX/templ)

func ListProductsHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortBy := r.URL.Query().Get("sort")
		var (
			productos []sqlc.Producto
			err       error
		)

		switch sortBy {
		case "price-asc":
			productos, err = queries.ListProductsByPriceAsc(r.Context())
		case "price-desc":
			productos, err = queries.ListProductsByPriceDesc(r.Context())
		default:
			productos, err = queries.ListProd(r.Context())
		}

		if err != nil {
			log.Printf("Error al obtener productos para templ: %v", err)
			http.Error(w, "Error al obtener productos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		componente := views.ProductList(productos)
		componente.Render(r.Context(), w)
	}
}

func LayoutHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		views.Layout().Render(r.Context(), w)
	}
}

func CarritoHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listCartHandler(queries)(w, r) // GET /carrito
		case http.MethodPost:
			createCartdHandler(queries)(w, r) // POST /carrito
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func createCartdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req sqlc.AddToCartParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if req.IDUsuario == 0 || req.IDProducto == 0 {
			http.Error(w, "ID de usuario y ID de producto son requeridos", http.StatusBadRequest)
			return
		}

		// AddToCart ahora usa :one y devuelve AddToCartRow
		cartItem, err := queries.AddToCart(r.Context(), req)
		if err != nil {
			http.Error(w, "Error al agregar al carrito: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(cartItem)
	}
}

// Producto: GET /products
func listCartHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		views.CarritoList([]sqlc.Carrito{}).Render(r.Context(), w)
	}
}
