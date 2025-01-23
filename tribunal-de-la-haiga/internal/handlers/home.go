package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Ruta de tu plantilla
	tmplPath := filepath.Join("templates", "home.html")
	tmpl := template.Must(template.ParseFiles(tmplPath))

	tmpl.Execute(w, nil)
}
