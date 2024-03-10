package sessions

import (
	"context"
	"database/sql"
	"time"

	"faisonz.net/cms/internal/users"
	"faisonz.net/cms/web/mux"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

type SessionData struct {
	User     *users.User
	LoggedIn bool
}

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

func GetSessionData(session *scs.SessionManager, ctx context.Context) SessionData {
	var data SessionData

	data.User, data.LoggedIn = mux.GetUserFromContext(ctx)

	return data
}
