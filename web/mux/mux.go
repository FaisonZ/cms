package mux

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"faisonz.net/cms/internal/users"
	"github.com/alexedwards/scs/v2"
)

type AuthMux struct {
	*http.ServeMux
	Session *scs.SessionManager
	DB      *sql.DB
}

type RouteHandlersFunc func(*AuthMux)

type key string

const authUserKey key = "user"

func NewAuthMux(session *scs.SessionManager, db *sql.DB) *AuthMux {
	return &AuthMux{
		ServeMux: http.NewServeMux(),
		Session:  session,
		DB:       db,
	}
}

func (mux *AuthMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("In the AuthMux")
	mux.withAuthedUser(mux.ServeMux).ServeHTTP(w, r)
}

func (mux *AuthMux) withAuthedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("In authed user")

		userID, ok := mux.Session.Get(r.Context(), "user_id").(int)
		if !ok {
			log.Println("No user id in session")
			next.ServeHTTP(w, r)
			return
		}

		user, err := users.GetUserByID(userID, mux.DB)
		if err != nil {
			log.Println("no user found for use id in session")
			next.ServeHTTP(w, r)
			return
		}

		log.Println("User found in session:", user.Username)

		ctx := context.WithValue(r.Context(), authUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mux *AuthMux) ProtectHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
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

func (mux *AuthMux) NoAuthOnlyHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
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

func (mux *AuthMux) ReceiveRouteHandlers(rff RouteHandlersFunc) {
	rff(mux)
}

func GetUserFromContext(ctx context.Context) (*users.User, bool) {
	user, ok := ctx.Value(authUserKey).(*users.User)
	return user, ok
}
