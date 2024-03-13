package clients

import "database/sql"

type Client struct {
	ID     int
	Name   string
	Active bool
}

func SaveClient(client *Client, db *sql.DB) (int, error) {
	active := 0
	if client.Active {
		active = 1
	}

	result, err := db.Exec(
		"INSERT INTO clients (name, active) VALUES (?, ?)",
		client.Name,
		active,
	)

	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}
