package mux

import (
	"context"
	"log"
	"net/http"

	"faisonz.net/cms/internal/db"
	"faisonz.net/cms/internal/sessions"
	"faisonz.net/cms/internal/users"
)

type AuthMux struct {
	*http.ServeMux
	Session *sessions.Session
	DBM     *db.DBManager
}

type RouteHandlersFunc func(*AuthMux)

type TemplateData struct {
	LoggedIn bool
	User     users.User
}

type key string

const authUserKey key = "user"

func NewAuthMux(session *sessions.Session, dbm *db.DBManager) *AuthMux {
	return &AuthMux{
		ServeMux: http.NewServeMux(),
		Session:  session,
		DBM:      dbm,
	}
}

func (m *AuthMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("In the AuthMux")

	userID := m.Session.GetUserID(r.Context())
	if userID == -1 {
		log.Println("No user id in session")
		m.ServeMux.ServeHTTP(w, r)
		return
	}

	user, err := users.GetUserByID(userID, m.DBM.Main)
	if err != nil {
		log.Println("no user found for use id in session")
		m.ServeMux.ServeHTTP(w, r)
		return
	}

	log.Println("User found in session:", user.Username)

	ctx := context.WithValue(r.Context(), authUserKey, *user)
	m.ServeMux.ServeHTTP(w, r.WithContext(ctx))
}

func (m *AuthMux) ProtectHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking auth status...")
		_, loggedIn := GetUserFromContext(r.Context())
		if !loggedIn {
			log.Println("Not authenticated, blocking")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		log.Println("Authenticated, letting in")

		handler(w, r)
	})
}

func (m *AuthMux) NoAuthOnlyHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking auth status...")
		_, loggedIn := GetUserFromContext(r.Context())
		if loggedIn {
			log.Println("authenticated, redirecting")
			http.Redirect(w, r, "/animals", http.StatusFound)
			return
		}

		log.Println("Not authenticated, letting in")

		handler(w, r)
	})
}

func (m *AuthMux) ReceiveRouteHandlers(rff RouteHandlersFunc) {
	rff(m)
}

func (m *AuthMux) GetTemplateData(ctx context.Context) TemplateData {
	var data TemplateData
	data.User, data.LoggedIn = GetUserFromContext(ctx)
	return data
}

func GetUserFromContext(ctx context.Context) (users.User, bool) {
	user, ok := ctx.Value(authUserKey).(users.User)

	// Don't accidentally leak a pass hash through a template
	user.Password = ""

	return user, ok
}
