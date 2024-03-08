package db

import (
	"database/sql"
	"fmt"
	"log"
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

func dbMain() {
	db, err := loadDB()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if closeError := db.Close(); closeError != nil {
			fmt.Println("Error closing database", closeError)
		}
	}()

	// insertUsers(db)
	queryUsers(db)

	defer db.Close()
}

func insertUsers(db *sql.DB) {
	for i := 0; i < 10; i++ {
		_, err := db.Exec("INSERT INTO users (name) VALUES (?)", fmt.Sprintf("test-%d", i))
		if err != nil {
			log.Fatal("Insert User:", err)
		}
	}
}

type User struct {
	ID   int
	Name string
}

func queryUsers(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		users = append(users, user)
		fmt.Println(user.ID, user.Name)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}
}

func insertUser(db *sql.DB) {
	name := "Bubba"
	result, err := db.Exec("INSERT INTO users (name) VALUES (?);", name)
	if err != nil {
		fmt.Println("Error inserting:", err)
		os.Exit(1)
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Did not insert:", err)
	}

	fmt.Printf("Inserted with ID: %d", id)
}
