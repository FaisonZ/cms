package sessions

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"faisonz.net/cms/internal/users"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

type Session struct {
	sessionManager *scs.SessionManager
}

type SessionData struct {
	User     *users.User
	LoggedIn bool
}

func New(db *sql.DB) *Session {
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = time.Hour * 4

	return &Session{
		sessionManager: sessionManager,
	}
}

func (s *Session) LoadAndSave(next http.Handler) http.Handler {
	return s.sessionManager.LoadAndSave(next)
}

func (s *Session) Destroy(ctx context.Context) error {
	return s.sessionManager.Destroy(ctx)
}

func (s *Session) GetUserID(ctx context.Context) int {
	userID, ok := s.sessionManager.Get(ctx, "user_id").(int)
	if !ok {
		return -1
	}

	return userID
}

func (s *Session) PutUserID(ctx context.Context, id int) {
	s.sessionManager.Put(ctx, "user_id", id)
	// TODO: Connect session id with user id
	// sessionID, _, err := mux.Session.Commit(r.Context())
	// if err != nil {
	// 	http.Error(w, "Oops", http.StatusInternalServerError)
	// 	return
	// }

	// if err := sessions.LinkSessionWithUser(sessionID, user, mux.DB); err != nil {
	// 	http.Error(w, "Oops", http.StatusInternalServerError)
	// }
}

func LinkSessionWithUser(sessionID string, user *users.User, db *sql.DB) error {
	if _, err := db.Exec("INSERT INTO user_sessions (user_id, session_token) VALUES (?, ?)", user.ID, sessionID); err != nil {
		return err
	}

	return nil
}
