package db

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strings"

	"faisonz.net/cms/internal/animals"
	"github.com/joho/godotenv"
)

func SetupDatabase() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("godotenv.Load(): ", err)
	}

	db, err := LoadDB()
	if err != nil {
		return err
	}
	defer func() {
		if closeError := db.Close(); closeError != nil {
			fmt.Println("Error closing database", closeError)
		}
	}()

	tables, err := getMissingTables(db)
	if err != nil {
		return err
	}

	if err := createTables(tables, db); err != nil {
		return err
	}

	if err := insertDefaultData(tables, db); err != nil {
		return err
	}

	fmt.Println("Database created!")

	return nil
}

func getMissingTables(db *sql.DB) ([]string, error) {
	tables := []string{
		"users",
		"sessions",
		"user_sessions",
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
		case "sessions":
			createStmt = "CREATE TABLE sessions (\n" +
				"token TEXT PRIMARY KEY,\n" +
				"data BLOB NOT NULL,\n" +
				"expiry REAL NOT NULL" +
				");\n" +
				"CREATE INDEX sessions_expiry_idx ON sessions(expiry);"
		case "users":
			createStmt = "CREATE TABLE users (\n" +
				"id INTEGER PRIMARY KEY AUTOINCREMENT,\n" +
				"username TEXT UNIQUE NOT NULL,\n" +
				"password TEXT NOT NULL" +
				")"
		case "user_sessions":
			// Cascase delete not working
			createStmt = "CREATE TABLE user_sessions (\n" +
				"user_id INTEGER NOT NULL,\n" +
				"session_token TEXT NOT NULL,\n" +
				"FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,\n" +
				"FOREIGN KEY (session_token) REFERENCES sessions(token) ON DELETE CASCADE\n" +
				")"
		default:
			continue
		}

		if _, err := db.Exec(createStmt); err != nil {
			return err
		}

		fmt.Println("Created table:", tName)
	}

	return nil
}

func insertDefaultData(tNames []string, db *sql.DB) error {
	for _, tName := range tNames {
		switch tName {
		case "animal_types":
			if err := insertAnimalTypes(db); err != nil {
				return err
			}
		}
	}

	return nil
}

func insertAnimalTypes(db *sql.DB) error {
	animalTypes := []animals.AnimalType{
		{ID: 1, Name: "Cow"},
		{ID: 2, Name: "Chicken"},
	}

	insertStmt := "INSERT INTO animal_types (id, name) VALUES\n" +
		"(?, ?)" + strings.Repeat(",\n(?, ?)", len(animalTypes)-1)

	args := make([]any, len(animalTypes)*2)
	for i, aType := range animalTypes {
		args[2*i] = aType.ID
		args[2*i+1] = aType.Name
	}

	result, err := db.Exec(insertStmt, args...)
	if err != nil {
		return err
	}

	if num, err := result.RowsAffected(); err != nil {
		fmt.Println("Can't get rows affected")
	} else {
		fmt.Println("Animal Types inserted:", num)
	}

	return nil
}
