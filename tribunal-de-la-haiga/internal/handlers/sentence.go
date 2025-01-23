package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pinlinsan/tribunaldelahaiga/internal/services"
)

func FormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusBadRequest)
		return
	}

	falta := r.FormValue("falta_ortografia")
	demandado := r.FormValue("nombre_demandado")
	demandante := r.FormValue("nombre_demandante")
	fecha := r.FormValue("fecha")

	// Llamada a la API de ChatGPT
	sentencia, err := services.GenerarSentencia(falta, demandado, demandante, fecha)
	if err != nil {
		// Imprime en la consola
		log.Printf("Error llamando a OpenAI: %v", err)

		// También podrías mostrarlo al usuario (aunque no siempre es deseable por seguridad):
		http.Error(w, "Error al generar una sentencia: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Guardar la sentencia en memoria y obtener ID
	id := services.GuardarSentencia(sentencia)

	// Redirigir a /sentencia/{id}
	http.Redirect(w, r, fmt.Sprintf("/sentencia/%s", id), http.StatusSeeOther)
}

func SentenceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sentencia := services.ObtenerSentenciaPorID(id)
	if sentencia == "" {
		http.Error(w, "Sentencia no encontrada", http.StatusNotFound)
		return
	}

	// Mostrar la sentencia directamente.
	fmt.Fprintf(w, "<h1>Sentencia del Tribunal de la Haiga</h1><p>%s</p>", sentencia)
}
