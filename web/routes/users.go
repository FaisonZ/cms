package routes

import (
	"database/sql"
	"html/template"
	"net/http"

	"faisonz.net/cms/internal/users"
	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserRoutes(mux *http.ServeMux, sessionManager *scs.SessionManager, db *sql.DB) error {
	registerTempl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/users/register.html"))
	loginTempl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/users/login.html"))

	mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		registerTempl.Execute(w, nil)
	})

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
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

		if err := users.SaveNew(&newUser, db); err != nil {
			http.Error(w, "Oops", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		loginTempl.Execute(w, nil)
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		rawPass := r.FormValue("password")

		user, err := users.GetUserByUsername(username, db)
		if err != nil || user == nil {
			http.Error(w, "Username or Password incorrect", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPass)); err != nil {
			http.Error(w, "Username or Password incorrect", http.StatusUnauthorized)
			return
		}

		sessionManager.Put(r.Context(), "user_id", user.ID)
		// sessionID, _, err := sessionManager.Commit(r.Context())
		// if err != nil {
		// 	http.Error(w, "Oops", http.StatusInternalServerError)
		// 	return
		// }

		// if err := sessions.LinkSessionWithUser(sessionID, user, db); err != nil {
		// 	http.Error(w, "Oops", http.StatusInternalServerError)
		// }

		http.Redirect(w, r, "/", http.StatusFound)
	})

	mux.HandleFunc("GET /logout", func(w http.ResponseWriter, r *http.Request) {
		sessionManager.Destroy(r.Context())
		http.Redirect(w, r, "/", http.StatusFound)
	})

	return nil
}
