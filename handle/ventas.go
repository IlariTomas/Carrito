package handle

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	sqlc "carrito.com/db/sqlc"
)

func SalesHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listVentasHandler(queries)(w, r) // GET /sales (Todas las ventas)
		case http.MethodPost:
			createVentaHandler(queries)(w, r) // POST /sales
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func createVentaHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req sqlc.CreateVentaParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if req.IDProducto == 0 || req.Total == "" || req.Cantidad == 0 {
			http.Error(w, "ID de Producto, Cantidad y Total son requeridos", http.StatusBadRequest)
			return
		}

		// CreateVenta ahora usa :one y devuelve CreateVentaRow
		venta, err := queries.CreateVenta(r.Context(), req)
		if err != nil {
			http.Error(w, "Error al crear venta: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(venta)
	}
}

// Venta: GET /sales (Lista TODAS las ventas)
func listVentasHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ListVentas ahora lista todas las ventas sin un parámetro de usuario
		ventas, err := queries.ListVentas(context.Background())
		if err != nil {
			http.Error(w, "Error al listar ventas: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ventas)
	}
}

func SaleHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getVentaHandler(queries)(w, r) // GET /sales/{id_venta}
		case http.MethodPut:
			updateVentaHandler(queries)(w, r) // PUT /sales/{id_venta}
		case http.MethodDelete:
			deleteVentaHandler(queries)(w, r) // DELETE /sales/{id_venta}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// Venta: GET /sales/{id}
func getVentaHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/sales/"):]
		idVenta, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de venta inválido", http.StatusBadRequest)
			return
		}

		venta, err := queries.GetVenta(r.Context(), int32(idVenta))
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Error al obtener venta: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(venta)
	}
}

// Venta: PUT /sales/{id}
func updateVentaHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/sales/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de venta inválido", http.StatusBadRequest)
			return
		}

		var req sqlc.UpdateVentaParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// Asignar el ID de la URL al struct de parámetros
		req.IDVenta = int32(id)

		err = queries.UpdateVenta(r.Context(), req)
		if err != nil {
			http.Error(w, "Error al actualizar venta: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Obtener y devolver la venta actualizada
		venta, _ := queries.GetVenta(r.Context(), int32(id))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(venta)
	}
}

// Venta: DELETE /sales/{id}
func deleteVentaHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/sales/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de venta inválido", http.StatusBadRequest)
			return
		}

		err = queries.DeleteVenta(r.Context(), int32(id))
		if err != nil {
			http.Error(w, "Error al eliminar venta: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
