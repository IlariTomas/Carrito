package handle

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time" 

	sqlc "carrito.com/db/sqlc"
)

// 1. DTO para la petición de CREAR una venta (acepta JSON 'snake_case')
type createSaleRequest struct {
	IDProducto int32  `json:"id_producto"`
	IDUsuario  int32  `json:"id_usuario"`
	Cantidad   int32  `json:"cantidad"`
	Total      string `json:"total"` // Recibimos como string
	Fecha      string `json:"fecha"` // Recibimos como string
}

// 2. DTO para la petición de ACTUALIZAR una venta
type updateSaleRequest struct {
	Cantidad int32  `json:"cantidad"`
	Total    string `json:"total"`
	Fecha    string `json:"fecha"`
}

// 3. DTO para la RESPUESTA de una venta (JSON limpio)
type saleResponse struct {
	IDVenta    int32  `json:"id_venta"`
	IDProducto int32  `json:"id_producto"`
	IDUsuario  int32  `json:"id_usuario"`
	Cantidad   int32  `json:"cantidad"`
	Total      string `json:"total"`
	Fecha      string `json:"fecha"` // Enviamos como string
}


func SalesHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listSalesHandler(queries)(w, r)
		case http.MethodPost:
			createSaleHandler(queries)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// ESTE ES EL HANDLER 'create' ARREGLADO
func createSaleHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Decodifica en el DTO simple (¡esto funciona!)
		var req createSaleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
			return
		}

		// 2. "Traduce" el string de fecha a sql.NullTime
		parsedTime, err := time.Parse(time.RFC3339, req.Fecha)
		if err != nil {
			http.Error(w, "Formato de fecha inválido, se espera RFC3339: "+err.Error(), http.StatusBadRequest)
			return
		}
		fechaNullTime := sql.NullTime{Time: parsedTime, Valid: true}

		// 3. Mapea al struct de sqlc (PascalCase)
		params := sqlc.CreateVentaParams{
			IDProducto: req.IDProducto,
			IDUsuario:  req.IDUsuario,
			Cantidad:   req.Cantidad,
			Total:      req.Total,
			Fecha:      fechaNullTime, 
		}

		// 4. Llama a la DB
		venta, err := queries.CreateVenta(r.Context(), params)
		if err != nil {
			http.Error(w, "Error al crear venta: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 5. Mapea la respuesta de la DB (compleja) al JSON (simple)
		respuesta := saleResponse{
			IDVenta:    venta.IDVenta,
			IDProducto: venta.IDProducto,
			IDUsuario:  venta.IDUsuario,
			Cantidad:   venta.Cantidad,
			Total:      venta.Total,
			Fecha:      venta.Fecha.Time.Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(respuesta)
	}
}

func listSalesHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ventas, err := queries.ListVentas(r.Context())
		if err != nil {
			http.Error(w, "Error al listar ventas: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Mapea la lista de structs complejos a una lista de JSONs simples
		respuestas := make([]saleResponse, len(ventas))
		for i, v := range ventas {
			respuestas[i] = saleResponse{
				IDVenta:    v.IDVenta,
				IDProducto: v.IDProducto,
				IDUsuario:  v.IDUsuario,
				Cantidad:   v.Cantidad,
				Total:      v.Total,
				Fecha:      v.Fecha.Time.Format(time.RFC3339),
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respuestas)
	}
}

// -----------------------------------------------------
// MANEJADOR PARA /sale/{id} (Individual)
// -----------------------------------------------------

func SaleHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getSaleHandler(queries)(w, r)
		case http.MethodPut:
			updateSaleHandler(queries)(w, r)
		case http.MethodDelete:
			deleteSaleHandler(queries)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getSaleHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/sale/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de venta inválido", http.StatusBadRequest)
			return
		}
		
		venta, err := queries.GetVenta(r.Context(), int32(id))
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
				return
			}
			http.Error(w, "Error al obtener venta: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Mapea la respuesta de la DB (compleja) al JSON (simple)
		respuesta := saleResponse{
			IDVenta:    venta.IDVenta,
			IDProducto: venta.IDProducto,
			IDUsuario:  venta.IDUsuario,
			Cantidad:   venta.Cantidad,
			Total:      venta.Total,
			Fecha:      venta.Fecha.Time.Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respuesta)
	}
}

// ESTE ES EL HANDLER 'update' ARREGLADO
func updateSaleHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/sale/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID de venta inválido", http.StatusBadRequest)
			return
		}
		
		var req updateSaleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
			return
		}
		
		parsedTime, err := time.Parse(time.RFC3339, req.Fecha)
		if err != nil {
			http.Error(w, "Formato de fecha inválido: "+err.Error(), http.StatusBadRequest)
			return
		}
		fechaNullTime := sql.NullTime{Time: parsedTime, Valid: true}

		params := sqlc.UpdateVentaParams{
			IDVenta:  int32(id),
			Cantidad: req.Cantidad,
			Total:    req.Total,
			Fecha:    fechaNullTime,
		}

		err = queries.UpdateVenta(r.Context(), params)
		if err != nil {
			http.Error(w, "Error al actualizar venta: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Obtiene y devuelve la venta actualizada (para hurl)
		venta, _ := queries.GetVenta(r.Context(), int32(id))
		respuesta := saleResponse{
			IDVenta:    venta.IDVenta,
			IDProducto: venta.IDProducto,
			IDUsuario:  venta.IDUsuario,
			Cantidad:   venta.Cantidad,
			Total:      venta.Total,
			Fecha:      venta.Fecha.Time.Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respuesta)
	}
}

func deleteSaleHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/sale/"):]
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