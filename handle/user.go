package handle

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	sqlc "carrito.com/db/sqlc"
	"carrito.com/views"
)

// usersHandler maneja todas las peticiones a /users
func UsersHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listUsersHandler(queries)(w, r) // GET /users
		case http.MethodPost:
			createUserHandler(queries)(w, r) // POST /users
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	}
}

// POST /users
func createUserHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error leyendo formulario", http.StatusBadRequest)
			return
		}
		nombreUsuario := r.FormValue("name_user")
		email := r.FormValue("email_user")

		req := sqlc.CreateUserParams{
			NombreUsuario: nombreUsuario,
			Email:         email,
		}

		if req.NombreUsuario == "" || req.Email == "" {
			http.Error(w, "Nombre y Email son requeridos", http.StatusBadRequest)
			return
		}

		_, err := queries.CreateUser(r.Context(), req)
		if err != nil {
			http.Error(w, "Error al crear usuario: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Recargar la lista de usuarios luego de crear uno
		users, err := queries.ListUsers(r.Context())
		if err != nil {
			http.Error(w, "Error cargando usuarios: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		views.ListUser(users).Render(r.Context(), w)
	}
}

// GET /users
func listUsersHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		views.UserView().Render(r.Context(), w)
	}
}

func UserHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func getUserHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/user/"):]
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
		idStr := r.URL.Path[len("/user/"):]
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
		idStr := r.URL.Path[len("/user/"):]
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
