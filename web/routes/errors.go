package routes

import (
	"html/template"
	"net/http"

	"faisonz.net/cms/web/mux"
)

func ErrorRouteHandlers(m *mux.AuthMux) {
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Serve404(m, w, r)
	})
}

func Serve404(m *mux.AuthMux, w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/errors/404.html"))
	w.WriteHeader(404)
	tmpl.Execute(w, m.GetTemplateData(r.Context()))
}
