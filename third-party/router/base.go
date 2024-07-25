package router

import (
	_ "third-party/docs"
	http_handler "third-party/handlers/http"
	"third-party/services/transaction"
	"third-party/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine, storage storage.Storage) {

	transactionService := transaction.NewTransactionService(storage)
	walletController := http_handler.NewWTransactionHandler(transactionService)

	r.POST("/third-party/payments", walletController.CreatePaymentTransaction)
	r.GET("/third-party/payments/:reference", walletController.GetPaymentTransaction)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
