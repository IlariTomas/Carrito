package handle

import (
	"log"
	"net/http"
	"strconv"
	"time"

	sqlc "carrito.com/db/sqlc"
	"carrito.com/views"
)

// --- LOGIN ---
func LoginHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ProcessLoginHandler(queries)(w, r) // Procesar el formulario (POST)
		} else {
			getLoginHandler()(w, r) // Mostrar el formulario (GET)
		}
	}
}

func getLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		views.LoginPage().Render(r.Context(), w)
	}
}

func ProcessLoginHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			views.AlertError("Error leyendo datos del formulario").Render(r.Context(), w)
			return
		}

		email := r.FormValue("email")
		usuario := r.FormValue("usuario")

		// Validación Mock
		usuarioValido := false
		if (email == "admin@test.com" && usuario == "1234") || email != "" {
			usuarioValido = true
		}

		// Verificar si el usuario existe en la base de datos
		user, err := queries.GetUserByEmail(r.Context(), email)
		if err != nil || user.Email == "" {
			views.AlertError("Usuario no encontrado. Regístrate primero.").Render(r.Context(), w)
			return
		}

		if !usuarioValido {
			views.AlertError("Credenciales incorrectas. Intenta de nuevo.").Render(r.Context(), w)
			return
		}

		CrearSesion(w, user)

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now().Add(-1 * time.Hour),
			Path:    "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// --- REGISTRO ---
func RegisterHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ProcessRegisterHandler(queries)(w, r) // Procesar registro (Insert en BD)
		} else {
			RegisterPageHandler()(w, r) // GET Mostrar página de registro
		}
	}
}

func RegisterPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		views.RegisterPage().Render(r.Context(), w)
	}
}

func ProcessRegisterHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			views.AlertError("Error leyendo formulario").Render(r.Context(), w)
			return
		}

		nombre := r.FormValue("usuario")
		email := r.FormValue("email")

		if nombre == "" || email == "" {
			views.AlertError("Nombre y Email son requeridos").Render(r.Context(), w)
			return
		}

		params := sqlc.CreateUserParams{
			NombreUsuario: nombre,
			Email:         email,
		}

		user, err := queries.CreateUser(r.Context(), params)
		if err != nil {
			views.AlertError("Error al registrar: prueba con otro usuario/email.").Render(r.Context(), w)
			return
		}
		log.Print(user)
		CrearSesion(w, user)

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}

func CrearSesion(w http.ResponseWriter, usuario sqlc.Usuario) {
	expiration := time.Now().Add(24 * time.Hour)

	cookie := http.Cookie{
		Name:     "session_token",
		Value:    strconv.Itoa(int(usuario.IDUsuario)),
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}
