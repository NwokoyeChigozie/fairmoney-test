package mongo_store

import (
	"context"
	"reflect"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db *MongoStore) Migrate(models []interface{}) error {
	ctx := context.Background()
	for _, model := range models {
		err := db.MigrateModel(ctx, model)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

// MigrateModel creates or updates the MongoDB collection associated with the given model.
func (db *MongoStore) MigrateModel(ctx context.Context, model interface{}) error {
	// Get the collection name from the model.
	name := CollectionName(model)

	// Get the collection from the database.
	coll := db.GetCollection(name)

	// Create the indexes and options from the model.
	indexes := make([]mongo.IndexModel, 0)
	collation := &options.Collation{
		Locale: "en_US",
	}
	structType := reflect.TypeOf(model).Elem()
	// fields := make([]string, 0)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		bsonTag, ok := field.Tag.Lookup("bson")
		if !ok {
			jsonTag, ok := field.Tag.Lookup("json")
			if !ok {
				bsonTag = field.Name
			} else {
				bsonTag = strings.Split(jsonTag, ",")[0]
			}
		} else if bsonTag == "-" {
			continue
		}
		// fields = append(fields, bsonTag)
		indexTag, ok := field.Tag.Lookup("index")
		if ok {
			indexModel := mongo.IndexModel{
				Keys:    bson.M{bsonTag: 1},
				Options: options.Index(),
			}
			if indexTag != "" {
				indexModel.Options.SetName(bsonTag)
			}

			if strings.Contains(indexTag, "unique") {
				indexModel.Options.SetName(bsonTag).SetUnique(true)
			}
			indexes = append(indexes, indexModel)
		}
	}

	if !db.CollectionExists(ctx, coll.Name()) {
		// Create the collection if it doesn't exist.
		if err := db.DB.CreateCollection(ctx, name, &options.CreateCollectionOptions{
			Collation: collation,
		}); err != nil {
			return err
		}
	}

	// Create the indexes if there are any.
	if len(indexes) > 0 {
		_, err := coll.Indexes().CreateMany(ctx, indexes)
		if err != nil {
			return err
		}
	}

	return nil
}

// CollectionName returns the name of the MongoDB collection associated with the given model.
func CollectionName(model interface{}) string {
	var (
		name                string
		collectionNameValue reflect.Value
		modelValue          = reflect.ValueOf(model)
		modelType           = reflect.TypeOf(model)
	)

	if modelType.Kind() == reflect.Ptr && modelType.Elem().Kind() == reflect.Pointer {
		name = modelType.Elem().Elem().Name()
		collectionNameValue = modelValue.Elem().Elem().MethodByName("CollectionName")
	} else if modelType.Kind() == reflect.Ptr && modelType.Elem().Kind() == reflect.Struct {
		name = modelType.Elem().Name()
		collectionNameValue = modelValue.Elem().MethodByName("CollectionName")
	} else if modelValue.Kind() == reflect.Pointer {
		collectionNameValue = modelValue.Elem().MethodByName("CollectionName")
		name = modelValue.Elem().Type().Name()
	} else {
		collectionNameValue = modelValue.Elem().MethodByName("CollectionName")
		name = reflect.TypeOf(model).Elem().Name()
		if name == "" {
			collectionNameValue = modelValue.MethodByName("CollectionName")
			name = reflect.TypeOf(model).Name()
		}
	}

	if collectionNameValue.IsValid() {
		result := collectionNameValue.Call([]reflect.Value{})
		if len(result) > 0 {
			collectionName := result[0].String()
			return collectionName
		}
	}

	pluralize := pluralize.NewClient()
	snake := strcase.ToSnake(name)
	words := strings.Split(snake, "_")
	plural := pluralize.Plural(words[len(words)-1])
	words[len(words)-1] = plural
	return strings.Join(words, "_")

}

// GetCollection returns the MongoDB collection with the given name from the context.
func GetCollection(ctx context.Context, name string) *mongo.Collection {
	if coll, ok := ctx.Value(name).(*mongo.Collection); ok {
		return coll
	} else {
		return nil
	}
}
