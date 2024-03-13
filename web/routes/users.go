package routes

import (
	"html/template"
	"net/http"

	"faisonz.net/cms/internal/clients"
	"faisonz.net/cms/internal/db"
	"faisonz.net/cms/internal/users"
	"faisonz.net/cms/web/mux"
	"golang.org/x/crypto/bcrypt"
)

func UserRouteHandlers(m *mux.AuthMux) {
	registerTempl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/users/register.html"))
	loginTempl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/users/login.html"))

	m.NoAuthOnlyHandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		templateData := m.GetTemplateData(r.Context())
		registerTempl.Execute(w, templateData)
	})

	m.NoAuthOnlyHandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		clientName := r.FormValue("clientname")
		username := r.FormValue("username")
		rawPass := r.FormValue("password")

		if user, err := users.GetUserByUsername(username, m.DBM.Main); err != nil {
			http.Error(w, "Oops", http.StatusUnauthorized)
			return
		} else if user != nil {
			http.Error(w, "Username taken", http.StatusUnauthorized)
			return
		}

		passBytes, err := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		client := &clients.Client{
			Name:   clientName,
			Active: true,
		}

		clientID, err := clients.SaveClient(client, m.DBM.Main)
		if err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}
		client.ID = clientID

		newUser := users.User{
			ClientID: client.ID,
			Username: username,
			Password: string(passBytes),
		}

		if err := db.SetupClientDB(client.ID, m.DBM); err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		if err := users.SaveNew(&newUser, m.DBM.Main); err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	m.NoAuthOnlyHandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		templateData := m.GetTemplateData(r.Context())
		loginTempl.Execute(w, templateData)
	})

	m.NoAuthOnlyHandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		rawPass := r.FormValue("password")

		user, err := users.GetUserByUsername(username, m.DBM.Main)
		if err != nil || user == nil {
			http.Error(w, "Username or Password incorrect", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPass)); err != nil {
			http.Error(w, "Username or Password incorrect", http.StatusUnauthorized)
			return
		}

		m.Session.PutUserID(r.Context(), user.ID)

		http.Redirect(w, r, "/animals", http.StatusFound)
	})

	m.HandleFunc("GET /logout", func(w http.ResponseWriter, r *http.Request) {
		m.Session.Destroy(r.Context())
		http.Redirect(w, r, "/", http.StatusFound)
	})
}
