package animals

import (
	"database/sql"
	"fmt"
	"strings"
)

type Animal struct {
	ID         int
	AnimalType AnimalType
	Name       string
}

func Get(id int, db *sql.DB) (*Animal, error) {
	if id < 1 {
		return nil, fmt.Errorf("id invalid")
	}

	rows, err := db.Query("SELECT id, type_id, name FROM animals WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	var anml Animal

	if err := rows.Scan(&anml.ID, &anml.AnimalType.ID, &anml.Name); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &anml, nil
}

func GetAll(db *sql.DB) ([]Animal, error) {
	rows, err := db.Query("SELECT id, type_id, name FROM animals")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	anmls := make([]Animal, 0, 10)

	for rows.Next() {
		var anml Animal

		if err := rows.Scan(&anml.ID, &anml.AnimalType.ID, &anml.Name); err != nil {
			return nil, err
		}

		anmls = append(anmls, anml)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return anmls, nil
}

func SaveMany(animal Animal, total int, db *sql.DB) error {
	if total < 1 || total > 20 {
		return fmt.Errorf("total out of rage, must be between 1 and 20")
	}

	insertStmt := "INSERT INTO animals (type_id) Values\n" +
		"(?)" + strings.Repeat(",\n(?)", total-1)

	args := make([]any, total)
	for i := 0; i < total; i++ {
		args[i] = animal.AnimalType.ID
	}

	result, err := db.Exec(insertStmt, args...)
	if err != nil {
		return err
	}

	if num, err := result.RowsAffected(); err != nil {
		fmt.Println("Can't get rows affected")
	} else {
		fmt.Println("Animals inserted:", num)
	}

	return nil
}
