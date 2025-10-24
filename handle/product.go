package handle

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	sqlc "carrito.com/db/sqlc"
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
		var req sqlc.CreateProdParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if req.NombreProducto == "" || req.Precio == "" {
			http.Error(w, "Nombre y Precio son requeridos", http.StatusBadRequest)
			return
		}

		// CreateProd ahora usa :one y devuelve CreateProdRow
		product, err := queries.CreateProd(r.Context(), req)
		if err != nil {
			http.Error(w, "Error al crear producto: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}

// Producto: GET /products
func listProdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := queries.ListProd(context.Background())
		if err != nil {
			http.Error(w, "Error al listar productos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
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
