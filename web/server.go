package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"faisonz.net/cms/internal/animals"
	"faisonz.net/cms/internal/db"
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

	animalTypes, err := animals.GetAnimalTypes(db)
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.ParseFiles("web/templates/layout.html"))

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
		fmt.Println(fp)
		_, err := os.Stat(fp)
		if err != nil {
			fmt.Println(err)
			serve404(w, r)
			return
		}

		http.ServeFile(w, r, fp)
	})

	type newAnimal struct {
		AnimalTypes []animals.AnimalType
	}
	newAnimalTmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/animals/new.html"))
	http.HandleFunc("GET /animals/new", func(w http.ResponseWriter, r *http.Request) {
		newAnimalTmpl.Execute(w, newAnimal{AnimalTypes: animalTypes})
	})

	http.HandleFunc("POST /animals", func(w http.ResponseWriter, r *http.Request) {
		aTypeID, err := strconv.ParseInt(r.FormValue("animalType"), 10, 0)
		if err != nil {
			http.Error(w, "Invalid type id", http.StatusBadRequest)
			return
		}

		total, err := strconv.ParseInt(r.FormValue("total"), 10, 0)
		if err != nil || total < 1 || total > 20 {
			http.Error(w, "Invalid total", http.StatusBadRequest)
			return
		}

		typeIndex := slices.IndexFunc(animalTypes, func(aType animals.AnimalType) bool {
			return aType.ID == int(aTypeID)
		})

		if typeIndex == -1 {
			http.Error(w, "Invalid animal type", http.StatusBadRequest)
			return
		}

		animal := animals.Animal{
			AnimalType: animalTypes[typeIndex],
		}

		if err := animals.SaveMany(animal, int(total), db); err != nil {
			http.Error(w, "Unexepected error saving animals", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serve404(w, r)
	})

	log.Println("Server started on port 3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func serve404(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/errors/404.html"))
	w.WriteHeader(404)
	tmpl.Execute(w, nil)
}
