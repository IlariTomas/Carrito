package handle

import (
	"net/http"
	"time"

	sqlc "carrito.com/db/sqlc"
	"carrito.com/views"
)

// GET /login -> Muestra el formulario
func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Renderizamos la vista que creamos arriba
		views.LoginPage().Render(r.Context(), w)
	}
}

// POST /login -> Procesa los datos
func ProcessLoginHandler(queries *sqlc.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Leer formulario
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error leyendo datos", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		usuario := r.FormValue("usuario")

		usuarioValido := false
		if email == "admin@test.com" && usuario == "tomi" {
			usuarioValido = true
		}

		if !usuarioValido {
			http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
			return
		}

		// Si es válido, le damos una cookie al navegador
		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{
			Name:    "session_token",
			Value:   "usuario_autenticado_123", // En la realidad, esto sería un token seguro
			Expires: expiration,
			Path:    "/", // Valido para toda la pagina
		}
		http.SetCookie(w, &cookie)

		// 4. Redirigir al Home
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Logout: Borra la cookie
func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now().Add(-1 * time.Hour), // Fecha pasada = borrar
			Path:    "/",
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
