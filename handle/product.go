package handle

import (
	"log"
	"net/http"
	"strconv"

	sqlc "carrito.com/db/sqlc"
	"carrito.com/views"
)

// Handler principal con proteccion de login
func IndexPageHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		idInt, err := strconv.Atoi(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		_, err = queries.GetUser(r.Context(), int32(idInt))
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Renderizar vista lista
		views.Layout().Render(r.Context(), w)
	}
}

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
		// Recargar la lista de productos luego de crear uno
		productos, err := queries.ListProd(r.Context())
		if err != nil {
			http.Error(w, "Error cargando productos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		views.ProductListDelete(productos).Render(r.Context(), w)
	}
}

// Producto: GET /products
func listProdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		views.ProductView().Render(r.Context(), w)
	}
}

// PRODUCTOS INDIVIDAULES
func ProductHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			deleteProdHandler(queries)(w, r) // DELETE /products/{id}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// Producto: DELETE /products/{id}
func deleteProdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/products/"):]
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

		productos, err := queries.ListProd(r.Context())
		if err != nil {
			http.Error(w, "Error cargando productos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		views.ProductListDelete(productos).Render(r.Context(), w)
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

func ListProductsViewHandler(queries *sqlc.Queries) http.HandlerFunc {
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

		componente := views.ProductListDelete(productos)
		componente.Render(r.Context(), w)
	}
}

func LayoutHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		views.Layout().Render(r.Context(), w)
	}
}
