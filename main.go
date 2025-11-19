package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	sqlc "carrito.com/db/sqlc" // generado por sqlc
	"carrito.com/handle"
	_ "github.com/lib/pq"

	"github.com/rs/cors"
)

func main() {

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Página /about
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "about.html")
	})

	connStr := "postgres://postgres:postgres@db:5432/apirest?sslmode=disable"
	db, err1 := sql.Open("postgres", connStr)

	if err1 != nil {
		log.Fatalf("failed to connect to DB: %v", err1)
	}
	defer db.Close()

	queries := sqlc.New(db)

	//Rutas
	mux.HandleFunc("/", handle.IndexPageHandler(queries))

	// 2. Rutas de Autenticación (NUEVAS)
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Procesar el formulario (POST)
			handle.ProcessLoginHandler(queries)(w, r)
		} else {
			// Mostrar el formulario (GET)
			handle.LoginHandler()(w, r)
		}
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Procesar registro (Insert en BD)
			handle.ProcessRegisterHandler(queries)(w, r)
		} else {
			// Mostrar página de registro
			handle.RegisterPageHandler()(w, r)
		}
	})

	mux.HandleFunc("/logout", handle.LogoutHandler())
	mux.HandleFunc("/products", handle.ProductsHandler(queries))
	mux.HandleFunc("/carrito", handle.CartHandler(queries))
	mux.HandleFunc("/list-products", handle.ListProductsHandler(queries))

	// 1. Crear una instancia de CORS que permite cualquier origen (*).
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// 2. Envolver el router (mux) con el handler de CORS.
	handler := c.Handler(mux)

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	// 3. Pasar el handler envuelto (handler) al ListenAndServe
	err := http.ListenAndServe(port, handler)
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
