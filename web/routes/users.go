package routes

import (
	"html/template"
	"net/http"

	"faisonz.net/cms/internal/users"
	"faisonz.net/cms/web/mux"
	"golang.org/x/crypto/bcrypt"
)

func UserRouteHandlers(mux *mux.AuthMux) {
	registerTempl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/users/register.html"))
	loginTempl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/users/login.html"))

	mux.NoAuthOnlyHandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		registerTempl.Execute(w, nil)
	})

	mux.NoAuthOnlyHandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		rawPass := r.FormValue("password")

		passBytes, err := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		newUser := users.User{
			Username: username,
			Password: string(passBytes),
		}

		if err := users.SaveNew(&newUser, mux.DB); err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	mux.NoAuthOnlyHandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		loginTempl.Execute(w, nil)
	})

	mux.NoAuthOnlyHandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		rawPass := r.FormValue("password")

		user, err := users.GetUserByUsername(username, mux.DB)
		if err != nil || user == nil {
			http.Error(w, "Username or Password incorrect", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPass)); err != nil {
			http.Error(w, "Username or Password incorrect", http.StatusUnauthorized)
			return
		}

		mux.Session.PutUserID(r.Context(), user.ID)

		http.Redirect(w, r, "/animals", http.StatusFound)
	})

	mux.HandleFunc("GET /logout", func(w http.ResponseWriter, r *http.Request) {
		mux.Session.Destroy(r.Context())
		http.Redirect(w, r, "/", http.StatusFound)
	})
}
