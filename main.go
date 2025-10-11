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

	connStr := "user=postgres password=XYZ dbname=tpcarrito host=localhost port=5432 sslmode=disable"
	db, err1 := sql.Open("postgres", connStr)
	if err1 != nil {
		log.Fatalf("failed to connect to DB: %v", err1)
	}
	defer db.Close()
	queries := sqlc.New(db)
	//ctx := context.Background()

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

func createUserHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req sqlc.CreateUserParams
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.Nombre == "" || req.Email == "" {
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

		if req.Nombre == "" || req.Email == "" {
			http.Error(w, "Nombre y Email son requeridos", http.StatusBadRequest)
			return
		}

		err = queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
			IDUsuario: int32(id),
			Nombre:    req.Nombre,
			Email:     req.Email,
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
