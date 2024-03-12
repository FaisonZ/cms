package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"faisonz.net/cms/internal/db"
	"faisonz.net/cms/internal/sessions"
	"faisonz.net/cms/web/mux"
	"faisonz.net/cms/web/routes"
	"github.com/joho/godotenv"
)

func StartServer() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("godotenv.Load(): ", err)
	}

	dbm, err := db.NewDBManager()
	if err != nil {
		log.Fatal(err)
	}

	sessionManager := sessions.New(dbm.Main)
	authMux := mux.NewAuthMux(sessionManager, dbm)

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html"))

	authMux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		templateData := authMux.GetTemplateData(r.Context())
		tmpl.Execute(w, templateData)
	})

	authMux.HandleFunc("GET /static/styles/{file}", func(w http.ResponseWriter, r *http.Request) {
		file := r.PathValue("file")

		if !strings.HasSuffix(file, ".css") {
			routes.Serve404(w, r)
			return
		}

		fp := filepath.Join("static", "styles", file)
		fmt.Println(fp)
		_, err := os.Stat(fp)
		if err != nil {
			fmt.Println(err)
			routes.Serve404(w, r)
			return
		}

		http.ServeFile(w, r, fp)
	})

	authMux.ReceiveRouteHandlers(routes.UserRouteHandlers)
	authMux.ReceiveRouteHandlers(routes.AnimalRouteHandlers)

	authMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		routes.Serve404(w, r)
	})

	log.Println("Server started on port 3000")

	log.Fatal(http.ListenAndServe(":3000", sessionManager.LoadAndSave(authMux)))
}
