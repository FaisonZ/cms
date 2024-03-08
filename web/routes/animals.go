package routes

import (
	"database/sql"
	"html/template"
	"net/http"
	"slices"
	"strconv"

	"faisonz.net/cms/internal/animals"
)

func RegisterAnimalRoutes(db *sql.DB) error {
	type indexPage struct {
		Animals []animals.Animal
	}
	type newAnimal struct {
		AnimalTypes []animals.AnimalType
	}

	animalTypes, err := animals.GetAnimalTypes(db)
	if err != nil {
		return err
	}

	indexTmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/animals/index.html"))
	newAnimalTmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/animals/new.html"))

	http.HandleFunc("GET /animals", func(w http.ResponseWriter, r *http.Request) {
		anmls, err := animals.GetAll(db, animalTypes)
		if err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}
		indexTmpl.Execute(w, indexPage{Animals: anmls})
	})

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

		http.Redirect(w, r, "/animals", http.StatusFound)
	})

	return nil
}
