package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pinlinsan/tribunaldelahaiga/internal/services"
)

// Template para renderizar la sentencia
var sentenciaTemplate = template.Must(template.New("sentence").Parse(`
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Sentencia del Tribunal de la Haiga</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            margin: 0 auto;
            padding: 20px;
            max-width: 800px;
            background-color: #f9f9f9;
        }
        h1 {
            color: #333;
            text-align: center;
        }
        .content {
            margin-top: 20px;
            background: #fff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
    </style>
</head>
<body>
    <h1>Sentencia del Tribunal de la Haiga</h1>
    <div class="content">{{ . }}</div>
</body>
</html>
`))

// Maneja el formulario enviado desde la página principal
func FormHandler(w http.ResponseWriter, r *http.Request) {
	// Parsear los datos del formulario
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error al parsear el formulario", http.StatusBadRequest)
		return
	}

	// Extraer los valores del formulario
	falta := r.FormValue("falta_ortografia")
	demandado := r.FormValue("nombre_demandado")
	demandante := r.FormValue("nombre_demandante")
	fecha := r.FormValue("fecha")

	// Log para verificar los datos ingresados
	log.Printf("Formulario recibido - Falta: %s, Demandado: %s, Demandante: %s, Fecha: %s", falta, demandado, demandante, fecha)

	// Llamada a la API de ChatGPT para generar la sentencia
	sentencia, err := services.GenerarSentencia(falta, demandado, demandante, fecha)
	if err != nil {
		log.Printf("Error llamando a OpenAI: %v", err)
		http.Error(w, "Error al generar una sentencia: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Guardar la sentencia en memoria y obtener su ID
	id := services.GuardarSentencia(sentencia)

	// Log para verificar el ID generado
	log.Printf("Sentencia guardada con ID: %s", id)

	// Redirigir al endpoint de la sentencia generada
	http.Redirect(w, r, "/sentencia/"+id, http.StatusSeeOther)
}

// Maneja la visualización de una sentencia por ID
func SentenceHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID de la sentencia desde la URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Obtener la sentencia por ID
	sentencia := services.ObtenerSentenciaPorID(id)

	// Log para verificar el contenido de la sentencia
	log.Printf("Sentencia obtenida para ID %s: %s", id, sentencia)

	// Verificar si la sentencia existe
	if sentencia == "" {
		http.Error(w, "Sentencia no encontrada", http.StatusNotFound)
		return
	}

	// Convertir la sentencia en template.HTML para renderizar correctamente etiquetas HTML
	renderedSentencia := template.HTML(sentencia)

	// Configurar el encabezado de la respuesta
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Renderizar la sentencia con el template
	err := sentenciaTemplate.Execute(w, renderedSentencia)
	if err != nil {
		http.Error(w, "Error al renderizar la sentencia", http.StatusInternalServerError)
		log.Printf("Error al ejecutar el template: %v", err)
	}
}
