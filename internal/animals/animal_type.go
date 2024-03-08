package animals

import "database/sql"

type AnimalType struct {
	ID   int
	Name string
}

func GetAnimalTypes(db *sql.DB) ([]AnimalType, error) {
	rows, err := db.Query("SELECT id, name FROM animal_types")
	if err != nil {
		return nil, err
	}

	aTypes := make([]AnimalType, 0, 5)

	for rows.Next() {
		aType := AnimalType{}

		if err := rows.Scan(&aType.ID, &aType.Name); err != nil {
			return nil, err
		}

		aTypes = append(aTypes, aType)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return aTypes, nil
}
