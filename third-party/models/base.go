package models

func ModelsForMigration() []interface{} {
	return []interface{}{
		&Transaction{},
	}
}
