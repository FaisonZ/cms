package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type DBManager struct {
	Main    *sql.DB
	clients map[int]*sql.DB
}

func NewDBManager() (*DBManager, error) {
	main, err := loadDB()
	if err != nil {
		return nil, err
	}

	return &DBManager{
		Main:    main,
		clients: make(map[int]*sql.DB, 5),
	}, nil
}

// Checks clients for existing open db
// returns if found
// Opens new *sql.DB if not, stores in clients
// func (dbm *DBManager) ClientDB(int) (*sql.DB, error)

func loadDB() (*sql.DB, error) {
	db_source := os.Getenv("DB")

	switch db_source {
	case "local":
		db, err := loadSQLite()
		if err != nil {
			return nil, err
		}
		return db, nil
	case "remote":
		db, err := loadLibSQL()
		if err != nil {
			return nil, err
		}
		return db, nil
	default:
		return nil, fmt.Errorf("load_db: Invalid value for db: %s", db_source)
	}
}

func loadSQLite() (*sql.DB, error) {
	err := os.Mkdir("local-db", 0750)
	if err != nil && !os.IsExist(err) {
		return nil, fmt.Errorf("cannot create local-db dir: %w", err)
	}

	fn := filepath.Join("local-db", "local.db")

	db, err := sql.Open("sqlite", fn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func loadLibSQL() (*sql.DB, error) {
	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", os.Getenv("DB_NAME"), os.Getenv("DB_TOKEN"))

	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}

	return db, nil
}
