package mongo_store

import "context"

func (db *MongoStore) SelectOneFromDb(receiver interface{}, query map[string]interface{}) (interface{}, error) {
	coll, err := db.GetCollectionForModel(receiver)
	if err != nil {
		return receiver, err
	}
	err = coll.FindOne(context.Background(), query).Decode(receiver)
	return receiver, err
}
