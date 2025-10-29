package handle

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	sqlc "carrito.com/db/sqlc"
	"carrito.com/views"
	
)

// -----------------------------------------------------
// MANEJADOR PARA /products (API de JSON)
// -----------------------------------------------------

func ProductsHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listProdHandler(queries)(w, r)
		case http.MethodPost:
			createProdHandler(queries)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// 1. DTO (struct simple) para la PETICIÓN (Request)
//    (Lo mantenemos para manejar 'snake_case' y 'precio' como string)
type createProductRequest struct {
	NombreProducto string `json:"nombre_producto"`
	Descripcion    string `json:"descripcion"`
	Precio         string `json:"precio"`
	Stock          int32  `json:"stock"`
	Categoria      string `json:"categoria"`
	Imagen         string `json:"imagen"`
}

// 2. createProdHandler (POST /products) - SIMPLIFICADO
func createProdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Decodifica en el DTO simple
		var req createProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Mapea al struct de sqlc (¡que ahora es simple!)
		params := sqlc.CreateProdParams{
			NombreProducto: req.NombreProducto,
			Stock:          req.Stock,
			Precio:         req.Precio,
			Descripcion:    req.Descripcion, // 👈 Ahora es string
			Categoria:      req.Categoria,   // 👈 Ahora es string
			Imagen:         req.Imagen,      // 👈 Ahora es string
		}

		// Llama a la DB
		// 'producto' (el struct de respuesta) ahora es simple
		producto, err := queries.CreateProd(r.Context(), params)
		if err != nil {
			http.Error(w, "Error al crear producto: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. Devuelve el struct de sqlc DIRECTAMENTE
		//    (Ya no tiene 'sql.NullString', así que el JSON es limpio)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(producto)
	}
}

// 3. listProdHandler (GET /products) - SIMPLIFICADO
func listProdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := queries.ListProd(context.Background())
		if err != nil {
			http.Error(w, "Error al listar productos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Devuelve el struct de sqlc DIRECTAMENTE
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

// -----------------------------------------------------
// MANEJADOR PARA /product/{id} (API de JSON)
// -----------------------------------------------------

func ProductHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProdHandler(queries)(w, r)
		case http.MethodPut:
			updateProdHandler(queries)(w, r)
		case http.MethodDelete:
			deleteProdHandler(queries)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// 1. DTO para la petición de ACTUALIZAR
type updateProductRequest struct {
	NombreProducto string `json:"nombre_producto"`
	Descripcion    string `json:"descripcion"`
	Precio         string `json:"precio"`
	Stock          int32  `json:"stock"`
	Categoria      string `json:"categoria"`
	Imagen         string `json:"imagen"`
}

// 2. updateProdHandler (PUT /product/{id}) - SIMPLIFICADO
func updateProdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtiene el ID
		idStr := r.URL.Path[len("/product/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de producto inválido", http.StatusBadRequest)
			return
		}

		// Decodifica en el DTO (acepta snake_case)
		var req updateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Mapea al struct de sqlc (¡que ahora es simple!)
		params := sqlc.UpdateProductoParams{
			IDProducto:     int32(id),
			NombreProducto: req.NombreProducto,
			Descripcion:    req.Descripcion, // 👈 Ahora es string
			Stock:          req.Stock,
			Precio:         req.Precio,
			Categoria:      req.Categoria,   // 👈 Ahora es string
			Imagen:         req.Imagen,      // 👈 Ahora es string
		}

		// Llama a la DB
		err = queries.UpdateProducto(r.Context(), params)
		if err != nil {
			http.Error(w, "Error al actualizar producto: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Devuelve el producto actualizado (para el HTTP 200 que espera hurl)
		producto, err := queries.GetProd(r.Context(), int32(id))
		if err != nil {
			http.Error(w, "Error al obtener producto actualizado: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200 OK
		json.NewEncoder(w).Encode(producto)
	}
}

// 3. getProdHandler (GET /product/{id}) - SIMPLIFICADO
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
		// Devuelve el struct de sqlc DIRECTAMENTE
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

// 4. deleteProdHandler (DELETE /product/{id}) - (Estaba bien)
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

// -----------------------------------------------------
// MANEJADOR PARA /list-products (HTMX/templ)
// -----------------------------------------------------

// (Este handler estaba bien porque solo leía de 'sqlc.Producto',
//  que ya está simple gracias a 'sqlc generate')
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