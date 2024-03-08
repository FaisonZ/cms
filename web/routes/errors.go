package routes

import (
	"html/template"
	"net/http"
)

func Serve404(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/errors/404.html"))
	w.WriteHeader(404)
	tmpl.Execute(w, nil)
}
