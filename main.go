package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	sqlc "carrito.com/db/sqlc" // generado por sqlc
	_ "github.com/lib/pq"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Servir el archivo index.html en la ruta "/"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// Si la ruta no existe, mostrar 404
			notFoundHandler(w, r)
			return
		}
		http.ServeFile(w, r, "index.html")
	})

	// Página /about
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!DOCTYPE html><html><head>
        <title>Acerca de</title></head><body>
        <h1>Acerca del servidor</h1>
        <p>Este es un servidor web básico escrito en Go.</p>
        <ul>
            <li>Host localhost:8080 was resolved.</li>
            <li>IPv6: ::1</li>
            <li>IPv4: 127.0.0.1</li>
            <li>Connected to localhost (::1) port 8080</li>
            <li>GET / HTTP/1.1</li>
            <li>Host: localhost:8080</li>
            <li>Content-Type: text/html; charset=utf-8</li>
        </ul>
        </body></html>`)
	})

	connStr := "postgres://postgres:postgres@db:5432/apirest?sslmode=disable"
	db, err1 := sql.Open("postgres", connStr)

	if err1 != nil {
		log.Fatalf("failed to connect to DB: %v", err1)
	}
	defer db.Close()

	queries := sqlc.New(db)
	//ctx := context.Background()

	// Rutas NUEVAS para USUARIOS
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listUsersHandler(queries)(w, r) // GET /users
		case http.MethodPost:
			createUserHandler(queries)(w, r) // POST /users
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUserHandler(queries)(w, r) // GET /users/{id}
		case http.MethodPut:
			updateUserHandler(queries)(w, r) // PUT /users/{id}
		case http.MethodDelete:
			deleteUserHandler(queries)(w, r) // DELETE /users/{id}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Rutas NUEVAS para PRODUCTOS
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listProdHandler(queries)(w, r) // GET /products
		case http.MethodPost:
			createProdHandler(queries)(w, r) // POST /products
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Rutas NUEVAS para VENTAS
	mux.HandleFunc("/sales", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listVentasHandler(queries)(w, r) // GET /sales (Todas las ventas)
		case http.MethodPost:
			createVentaHandler(queries)(w, r) // POST /sales
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Esta ruta maneja GET, PUT, DELETE para /sales/{id_venta}
	mux.HandleFunc("/sales/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// ... (El código http.ListenAndServe)

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}

}

// Función para manejar errores 404
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, `<!DOCTYPE html><html><head>
    <title>404 - No encontrado</title></head><body>
    <h1>404 - Página no encontrada</h1>
    <p>Lo sentimos, la página que buscas no existe.</p>
    </body></html>`)
}

//-----------------------------------------------HANDLERS USUARIOS ---------------------------------------------

func createUserHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req sqlc.CreateUserParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.NombreUsuario == "" || req.Email == "" {
			http.Error(w, "Nombre y Email son requeridos", http.StatusBadRequest)
			return
		}

		user, err := queries.CreateUser(r.Context(), req)
		if err != nil {
			http.Error(w, "Error al crear usuario: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func listUsersHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		users, err := queries.ListUsers(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users) // Esto es clave
	}
}

func getUserHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/users/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		user, err := queries.GetUser(r.Context(), int32(id))
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Error al obtener usuario: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func updateUserHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/users/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		var req sqlc.CreateUserParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if req.NombreUsuario == "" || req.Email == "" {
			http.Error(w, "Nombre y Email son requeridos", http.StatusBadRequest)
			return
		}

		err = queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
			IDUsuario:     int32(id),
			NombreUsuario: req.NombreUsuario,
			Email:         req.Email,
		})
		if err != nil {
			http.Error(w, "Error al actualizar usuario: "+err.Error(), http.StatusInternalServerError)
			return
		}

		user, _ := queries.GetUser(r.Context(), int32(id))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func deleteUserHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/users/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}

		err = queries.DeleteUser(r.Context(), int32(id))
		if err != nil {
			http.Error(w, "Error al eliminar usuario: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ------------------------------------
// HANDLERS PARA PRODUCTOS
// ------------------------------------

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

// Producto: GET /products/{id}
func getProdHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/products/"):]
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
		idStr := r.URL.Path[len("/products/"):]
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

		w.WriteHeader(http.StatusNoContent)
	}
}

// ------------------------------------
// HANDLERS PARA VENTAS
// ------------------------------------

// Venta: POST /sales
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
