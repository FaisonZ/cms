package db

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/joho/godotenv"
)

func SetupDatabase() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("godotenv.Load(): ", err)
	}

	db, err := loadDB()
	if err != nil {
		return err
	}
	defer db.Close()

	tables, err := getMissingTables(db)
	if err != nil {
		return err
	}

	if err := createTables(tables, db); err != nil {
		return err
	}

	return nil
}

func getMissingTables(db *sql.DB) ([]string, error) {
	tables := []string{
		"animal_types",
		"animals",
	}

	tNames := make([]any, len(tables))
	for i, v := range tables {
		tNames[i] = v
	}

	queryString := `SELECT NAME FROM sqlite_master WHERE type="table" AND NAME IN (?` + strings.Repeat(", ?", len(tNames)-1) + `)`

	rows, err := db.Query(queryString, tNames...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tName string
		if err := rows.Scan(&tName); err != nil {
			return nil, err
		}

		i := slices.Index(tables, tName)

		if i == -1 {
			continue
		}

		tables = slices.Delete(tables, i, i+1)
	}

	return tables, nil
}

func createTables(tNames []string, db *sql.DB) error {
	for _, tName := range tNames {
		var createStmt string
		switch tName {
		case "animal_types":
			createStmt = "CREATE TABLE animal_types (\n" +
				"id INTEGER PRIMARY KEY AUTOINCREMENT,\n" +
				"name TEXT NOT NULL\n" +
				")"
		case "animals":
			createStmt = "CREATE TABLE animals (\n" +
				"id INTEGER PRIMARY KEY AUTOINCREMENT,\n" +
				"type_id INTEGER NOT NULL,\n" +
				"name TEXT NOT NULL DEFAULT '',\n" +
				"FOREIGN KEY (type_id) REFERENCES animal_types(id)\n" +
				")"
		default:
			continue
		}

		fmt.Println(createStmt)
		if _, err := db.Exec(createStmt); err != nil {
			return err
		}

		fmt.Println("Created table:", tName)
	}

	return nil
}
