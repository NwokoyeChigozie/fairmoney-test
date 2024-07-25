package main

import (
	"banking-app/config"
	"banking-app/models"
	"banking-app/router"
	"banking-app/storage/mongo_store"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// config
	config.Setup()

	//db connection and migration
	store := mongo_store.NewMongoStore()
	store.Connect(config.GetConfig().MongoDb.ConnectionString, config.GetConfig().MongoDb.DbName)
	store.GetConnection()
	store.Migrate(models.ModelsForMigration())
	store.SeedData()

	//router
	router.Setup(r, store)

	serverPort := config.GetConfig().ServerPort
	log.Printf("starting server on 127.0.0.1:%v", serverPort)
	log.Fatal(r.Run(fmt.Sprintf(":%v", serverPort)))
}
