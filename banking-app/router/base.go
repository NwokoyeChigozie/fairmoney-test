package router

import (
	"banking-app/config"
	_ "banking-app/docs"
	http_handler "banking-app/handlers/http"
	third_party "banking-app/pkg/third-party"
	"banking-app/services/wallet"
	"banking-app/storage"

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine, storage storage.Storage) {

	thirdParty := third_party.NewThirdPartyPkg(config.GetConfig().ThirdPartyBaseUrl)
	walletService := wallet.NewWalletService(storage, thirdParty)
	walletController := http_handler.NewWalletHandler(walletService)

	r.POST("/bank-api/credit-wallet", walletController.CreditWallet)
	r.POST("/bank-api/debit-wallet", walletController.DebitWallet)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
