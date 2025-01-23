package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pinlinsan/tribunaldelahaiga/internal/handlers" // Ajusta la ruta según tu "module" en go.mod
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	// Rutas:
	// GET "/" -> muestra el formulario
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	// POST "/enviar-formulario" -> procesa los datos
	r.HandleFunc("/enviar-formulario", handlers.FormHandler).Methods("POST")
	// GET "/sentencia/{id}" -> muestra la sentencia generada
	r.HandleFunc("/sentencia/{id}", handlers.SentenceHandler).Methods("GET")

	log.Printf("Servidor escuchando en puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
