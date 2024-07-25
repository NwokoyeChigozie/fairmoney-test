package mongo_store

import (
	"banking-app/models"

	"go.mongodb.org/mongo-driver/bson"
)

// var seedUsersData = []
func (db *MongoStore) SeedData() {
	email := "test@gmail.com"

	user := models.User{Email: email}

	_, err := db.SelectOneFromDb(&user, bson.M{"email": email})
	if err != nil {
		_, err = db.CreateOneRecord(&user)
		if err != nil {
			panic(err)
		}
	}

	wallet := models.Wallet{
		AccountId: user.ID,
		Balance:   0,
	}
	_, err = db.SelectOneFromDb(&wallet, bson.M{"_id": user.ID})
	if err != nil {
		_, err = db.CreateOneRecord(&wallet)
		if err != nil {
			panic(err)
		}
	}

}
