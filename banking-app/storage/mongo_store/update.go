package mongo_store

import (
	"context"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *MongoStore) UpdateRecord(model interface{}) (interface{}, error) {
	coll, err := db.GetCollectionForModel(model)
	if err != nil {
		return model, err
	}
	filter := bson.M{"_id": getDocumentID(model)}
	update := bson.M{"$set": model}

	result, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model, err
	}

	update = bson.M{"$set": bson.M{"updated_at": time.Now()}}
	_, err = coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model, err
	}

	if result.ModifiedCount > 0 {
		_, err := db.SelectOneFromDb(model, bson.M{"_id": getDocumentID(model)})
		if err != nil {
			return model, err
		}
	}

	return model, nil
}

func getDocumentID(model interface{}) primitive.ObjectID {
	var (
		modelValue = reflect.ValueOf(model)
		modelType  = reflect.TypeOf(model)
		value      reflect.Value
	)

	if modelType.Kind() == reflect.Ptr && modelType.Elem().Kind() == reflect.Pointer {
		value = modelValue.Elem().Elem().FieldByName("ID")
	} else if modelType.Kind() == reflect.Ptr && modelType.Elem().Kind() == reflect.Struct {
		value = modelValue.Elem().FieldByName("ID")
	} else {
		value = modelValue.Elem().FieldByName("ID")
	}

	if !value.IsValid() {
		return primitive.NilObjectID
	}
	return value.Interface().(primitive.ObjectID)
}
