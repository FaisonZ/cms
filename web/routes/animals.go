package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"

	"faisonz.net/cms/internal/animals"
	"faisonz.net/cms/web/mux"
)

type indexPage struct {
	Animals []animals.Animal
}
type animalPage struct {
	Animal animals.Animal
}
type editAnimalPage struct {
	Animal      animals.Animal
	AnimalTypes []animals.AnimalType
}
type newAnimal struct {
	AnimalTypes []animals.AnimalType
}

func AnimalRouteHandlers(mux *mux.AuthMux) {
	animalTypes, err := animals.GetAnimalTypes(mux.DB)
	if err != nil {
		log.Fatal(err)
	}

	indexTmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/animals/index.html"))
	animalTmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/animals/show.html"))
	editAnimalTmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/animals/edit.html"))
	newAnimalTmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/animals/new.html"))

	mux.ProtectHandleFunc("GET /animals", func(w http.ResponseWriter, r *http.Request) {
		anmls, err := animals.GetAll(mux.DB)
		if err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}
		if err := animals.FillTypeForAnimals(anmls, animalTypes); err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}
		indexTmpl.Execute(w, indexPage{Animals: anmls})
	})

	mux.ProtectHandleFunc("GET /animals/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 0)
		if err != nil {
			Serve404(w, r)
			return
		}

		anml, err := animals.Get(int(id), mux.DB)
		if err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		if anml == nil {
			Serve404(w, r)
			return
		}

		if err := animals.FillTypeForAnimal(anml, animalTypes); err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}
		animalTmpl.Execute(w, animalPage{Animal: *anml})
	})

	mux.ProtectHandleFunc("GET /animals/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 0)
		if err != nil {
			Serve404(w, r)
			return
		}

		anml, err := animals.Get(int(id), mux.DB)
		if err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		if anml == nil {
			Serve404(w, r)
			return
		}

		if err := animals.FillTypeForAnimal(anml, animalTypes); err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}
		editAnimalTmpl.Execute(w, editAnimalPage{Animal: *anml, AnimalTypes: animalTypes})
	})

	mux.ProtectHandleFunc("POST /animals/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 0)
		if err != nil {
			Serve404(w, r)
			return
		}

		anml, err := animals.Get(int(id), mux.DB)
		if err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		if anml == nil {
			Serve404(w, r)
			return
		}

		anml.Name = r.FormValue("name")

		if err := animals.Update(*anml, mux.DB); err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/animals/%d", anml.ID), http.StatusSeeOther)
	})

	mux.ProtectHandleFunc("GET /animals/new", func(w http.ResponseWriter, r *http.Request) {
		newAnimalTmpl.Execute(w, newAnimal{AnimalTypes: animalTypes})
	})

	mux.ProtectHandleFunc("POST /animals", func(w http.ResponseWriter, r *http.Request) {
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

		if err := animals.SaveMany(animal, int(total), mux.DB); err != nil {
			http.Error(w, "Unexepected error saving animals", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/animals", http.StatusFound)
	})
}
