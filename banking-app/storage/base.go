package storage

type Storage interface {
	Connect(connectionString, dbName string)
	GetConnection()
	Migrate(models []interface{}) error
	CreateOneRecord(model interface{}) (interface{}, error)
	UpdateRecord(model interface{}) (interface{}, error)
	SelectOneFromDb(receiver interface{}, query map[string]interface{}) (interface{}, error)
	SeedData()
}
