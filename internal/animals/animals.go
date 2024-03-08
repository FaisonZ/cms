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
	defer rows.Close()

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

func Update(animal Animal, db *sql.DB) error {
	if animal.ID < 1 {
		return fmt.Errorf("id invalid")
	}

	updateStmt := "Update animals\n" +
		"SET name=?\n" +
		"WHERE id=?"

	result, err := db.Exec(updateStmt, animal.Name, animal.ID)
	if err != nil {
		return err
	}

	if num, err := result.RowsAffected(); err == nil && num != 1 {
		return fmt.Errorf("animal not updated")
	}

	return nil
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
