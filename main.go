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

	// Página /about

	connStr := "postgres://postgres:postgres@db:5432/apirest?sslmode=disable"
	db, err1 := sql.Open("postgres", connStr)

	if err1 != nil {
		log.Fatalf("failed to connect to DB: %v", err1)
	}
	defer db.Close()

	queries := sqlc.New(db)

	//Rutas
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "about.html")
	})
	mux.HandleFunc("/", handle.IndexPageHandler(queries))
	mux.HandleFunc("/login", handle.LoginHandler(queries))
	mux.HandleFunc("/register", handle.RegisterHandler(queries))
	mux.HandleFunc("/logout", handle.LogoutHandler())
	mux.HandleFunc("/users", handle.UsersHandler(queries))
	mux.HandleFunc("/products", handle.ProductsHandler(queries))
	mux.HandleFunc("/carrito", handle.CartHandler(queries))
	mux.HandleFunc("/carrito/items/", handle.CartItemHandler(queries))
	mux.HandleFunc("/list-products", handle.ListProductsHandler(queries))

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
