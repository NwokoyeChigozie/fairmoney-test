package mongo_store

import (
	"context"
	"fmt"
	"log"
	"strings"
	"third-party/storage"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	DB *mongo.Database
}

func NewMongoStore() storage.Storage {
	return &MongoStore{}
}

var connection *MongoStore

func (*MongoStore) Connect(conectionString, dbName string) {
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	clientOptions := options.Client().ApplyURI(conectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Panic(fmt.Sprintf("error connecting to db %v", err.Error()))
	}
	connection = &MongoStore{DB: client.Database(dbName)}
	log.Println("connected to mongo db")
}

func (m *MongoStore) GetConnection() {
	m.DB = connection.DB
	return
}

func (db *MongoStore) GetCollection(collectionName string) *mongo.Collection {
	return db.DB.Collection(collectionName)
}

func (db *MongoStore) CreateUniqueIndex(collName, field string, order int) error {
	collection := db.GetCollection(collName)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{field: order},
		Options: options.Index().SetUnique(true),
	}

	timeOutFactor := 3
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOutFactor)*time.Second)

	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("failed to create unique index on field %s in %s", field, collName)
		return fmt.Errorf("failed to create unique index on field %s in %s", field, collName)
	}

	return nil
}

func (db *MongoStore) CollectionExists(ctx context.Context, name string) bool {
	filter := bson.M{"name": name}
	collections, err := db.DB.ListCollectionNames(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, coll := range collections {
		if strings.EqualFold(coll, name) {
			return true
		}
	}
	return false
}

func (db *MongoStore) GetCollectionNameForModel(model interface{}) (string, error) {
	ctx := context.Background()
	name := CollectionName(model)
	if !db.CollectionExists(ctx, name) {
		return "", fmt.Errorf("collection for model %v does not exist, add to migrations and apply", name)
	}
	return name, nil
}

func (db *MongoStore) GetCollectionForModel(model interface{}) (*mongo.Collection, error) {
	name, err := db.GetCollectionNameForModel(model)
	if err != nil {
		return nil, err
	}
	return db.GetCollection(name), nil
}
