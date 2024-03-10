package users

import (
	"database/sql"
	"log"
)

type User struct {
	ID       int
	Username string
	Password string
}

func SaveNew(user *User, db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		user.Username,
		user.Password,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetUserByID(id int, db *sql.DB) (*User, error) {
	rows, err := db.Query("SELECT id, username, password FROM users WHERE id=?", id)
	if err != nil {
		log.Println("query error")
		return nil, err
	}
	defer rows.Close()

	var user User

	if !rows.Next() {
		return nil, nil
	}

	if err := rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(username string, db *sql.DB) (*User, error) {
	rows, err := db.Query("SELECT id, username, password FROM users WHERE username=?", username)
	if err != nil {
		log.Println("query error")
		return nil, err
	}
	defer rows.Close()

	var user User

	if !rows.Next() {
		return nil, nil
	}

	if err := rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}
