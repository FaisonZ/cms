package animals

import (
	"database/sql"
	"fmt"
	"slices"
)

type AnimalType struct {
	ID   int
	Name string
}

func FillTypeForAnimals(anmls []Animal, aTypes []AnimalType) error {
	for i := 0; i < len(anmls); i++ {
		anml := &anmls[i]
		if err := FillTypeForAnimal(anml, aTypes); err != nil {
			return err
		}
	}

	return nil
}

func FillTypeForAnimal(anml *Animal, aTypes []AnimalType) error {
	typeIndex := slices.IndexFunc(aTypes, func(aType AnimalType) bool {
		return aType.ID == anml.AnimalType.ID
	})

	if typeIndex == -1 {
		return fmt.Errorf("animal_type invalid")
	}

	anml.AnimalType.Name = aTypes[typeIndex].Name

	return nil
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
