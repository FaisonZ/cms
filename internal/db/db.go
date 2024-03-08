package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func loadDB() (*sql.DB, error) {
	db_source := os.Getenv("DB")
	fmt.Println(db_source)

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
	fn := "./local.db"

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
