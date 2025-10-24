package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	sqlc "carrito.com/db/sqlc" // generado por sqlc
	"carrito.com/handle"
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
	//Rutas
	mux.HandleFunc("/users", handle.UsersHandler(queries))
	mux.HandleFunc("/user/", handle.UserHandler(queries))

	mux.HandleFunc("/products", handle.ProductsHandler(queries))
	mux.HandleFunc("/product/", handle.ProductHandler(queries))

	mux.HandleFunc("/sales", handle.SalesHandler(queries))
	mux.HandleFunc("/sale/", handle.SaleHandler(queries))

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
