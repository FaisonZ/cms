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
