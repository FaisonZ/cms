package sessions

import (
	"database/sql"
	"time"

	"faisonz.net/cms/internal/users"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

func New(db *sql.DB) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = time.Hour * 4

	return sessionManager
}

func LinkSessionWithUser(sessionID string, user *users.User, db *sql.DB) error {
	if _, err := db.Exec("INSERT INTO user_sessions (user_id, session_token) VALUES (?, ?)", user.ID, sessionID); err != nil {
		return err
	}

	return nil
}
