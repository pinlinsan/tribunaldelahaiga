package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pinlinsan/tribunaldelahaiga/internal/services"
)

var sentenciaTemplate = template.Must(template.New("sentence").Parse(`
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sentencia del Tribunal de la Haiga</title>
</head>
<body>
    <h1>Sentencia del Tribunal de la Haiga</h1>
    <div>{{ . }}</div>
</body>
</html>
`))

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
	http.Redirect(w, r, "/sentencia/"+id, http.StatusSeeOther)
}

func SentenceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Obtener la sentencia por ID
	sentencia := services.ObtenerSentenciaPorID(id)
	if sentencia == "" {
		http.Error(w, "Sentencia no encontrada", http.StatusNotFound)
		return
	}

	// Renderizar la sentencia con el template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := sentenciaTemplate.Execute(w, sentencia)
	if err != nil {
		http.Error(w, "Error al renderizar la sentencia", http.StatusInternalServerError)
		log.Printf("Error al ejecutar el template: %v", err)
	}
}
