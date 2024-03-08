package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func StartServer() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("godotenv.Load(): ", err)
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html"))

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("GET /static/styles/{file}", func(w http.ResponseWriter, r *http.Request) {
		file := r.PathValue("file")

		if !strings.HasSuffix(file, ".css") {
			serve404(w, r)
			return
		}

		fp := filepath.Join("static", "styles", file)

		_, err := os.Stat(fp)
		if err != nil {
			fmt.Println(err)
			serve404(w, r)
			return
		}

		http.ServeFile(w, r, fp)
	})

	newAnimalTmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/animals/new.html"))
	http.HandleFunc("GET /animals/new", func(w http.ResponseWriter, r *http.Request) {
		newAnimalTmpl.Execute(w, nil)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serve404(w, r)
	})

	log.Println("Server started on port 3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func serve404(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/errors/404.html"))
	w.WriteHeader(404)
	tmpl.Execute(w, nil)
}
