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

	db, err := db.LoadDB()
	if err != nil {
		log.Fatal(db)
	}

	sessionManager := sessions.New(db)

	mux := mux.NewAuthMux(sessionManager, db)
	hndl := sessionManager.LoadAndSave(
		mux,
	)

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html"))

	type homeData struct {
		sessions.SessionData
	}

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		sessionData := sessions.GetSessionData(mux.Session, r.Context())
		tmpl.Execute(w, homeData{SessionData: sessionData})
	})

	mux.HandleFunc("GET /static/styles/{file}", func(w http.ResponseWriter, r *http.Request) {
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

	mux.ReceiveRouteHandlers(routes.UserRouteHandlers)
	mux.ReceiveRouteHandlers(routes.AnimalRouteHandlers)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		routes.Serve404(w, r)
	})

	log.Println("Server started on port 3000")

	log.Fatal(http.ListenAndServe(":3000", hndl))
}
