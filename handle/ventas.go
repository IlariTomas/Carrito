package handle

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	sqlc "carrito.com/db/sqlc"
	"carrito.com/views"
)

func SalesHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listVentasHandler(queries)(w, r) // GET
		case http.MethodPost:
			createVentaHandler(queries)(w, r) // POST
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func createVentaHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			views.AlertError("Debes iniciar sesión para comprar").Render(r.Context(), w)
			return
		}
		userID, err := strconv.Atoi(cookie.Value)
		if err != nil {
			views.AlertError("Debes iniciar sesión para comprar").Render(r.Context(), w)
			return
		}

		ctx := r.Context()
		cartItems, err := queries.GetCartItems(ctx, int32(userID))
		if err != nil || len(cartItems) == 0 {
			views.AlertError("El carrito está vacío o hubo un error").Render(ctx, w)
			return
		}

		for _, item := range cartItems {
			precioFloat, _ := strconv.ParseFloat(item.Precio, 64)
			totalLinea := precioFloat * float64(item.Cantidad)

			// Creamos la venta usando los parámetros de TU query
			ventaParams := sqlc.CreateVentaParams{
				IDProducto: item.IDProducto,
				IDUsuario:  int32(userID),
				Cantidad:   item.Cantidad,
				Total:      fmt.Sprintf("%.2f", totalLinea),
			}

			_, err := queries.CreateVenta(ctx, ventaParams)
			if err != nil {
				views.AlertError("Error procesando la compra").Render(ctx, w)
				return
			}
		}

		_ = queries.DeleteCart(ctx, int32(userID))
		views.AlertSuccess("¡Compra realizada con éxito!").Render(ctx, w)
	}
}

// Venta: GET /sales (Lista TODAS las ventas)
func listVentasHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_token")
		if err != nil {
			views.AlertError("Debes iniciar sesión para comprar").Render(r.Context(), w)
			return
		}
		userID, err := strconv.Atoi(cookie.Value)
		if err != nil {
			views.AlertError("Debes iniciar sesión para comprar").Render(r.Context(), w)
			return
		}

		ventas, err := queries.ListVentasUsuario(context.Background(), int32(userID))
		if err != nil {
			http.Error(w, "Error al listar ventas: "+err.Error(), http.StatusInternalServerError)
			return
		}

		views.SalesList(ventas).Render(r.Context(), w)
	}
}
