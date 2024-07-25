package mongo_store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (db *MongoStore) CreateOneRecord(model interface{}) (interface{}, error) {
	coll, err := db.GetCollectionForModel(model)
	if err != nil {
		return model, err
	}

	result, err := coll.InsertOne(context.Background(), model)
	if err != nil {
		return model, err
	}

	if result.InsertedID == nil {
		return model, fmt.Errorf("record creation for %v failed", coll.Name())
	}

	insertedID := result.InsertedID.(primitive.ObjectID)

	filter := bson.M{"_id": insertedID}
	update := bson.M{"$set": bson.M{"updated_at": time.Now(), "created_at": time.Now()}}
	re, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return model, err
	}

	if re.ModifiedCount > 0 {
		_, err := db.SelectOneFromDb(model, bson.M{"_id": insertedID})
		if err != nil {
			return model, err
		}
	}
	return model, nil
}
