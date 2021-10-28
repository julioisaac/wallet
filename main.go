package main

import (
	"github.com/julioisaac/daxxer-api/routers"
	"github.com/julioisaac/daxxer-api/routers/gin"
	"github.com/julioisaac/daxxer-api/src/helpers/ticker"
	"github.com/julioisaac/daxxer-api/src/wallet/account/controller"
	currencies "github.com/julioisaac/daxxer-api/src/wallet/currencies/controller"
	"github.com/julioisaac/daxxer-api/src/wallet/prices/service"
	"github.com/julioisaac/daxxer-api/storage"
	"github.com/julioisaac/daxxer-api/storage/mongodb"
)

var (
	dbConfig   storage.DBConfig = mongodb.NewMongoConfig()
	httpRouter routers.Router   = gin.NewGinRouter()
	daxxerTicker      = ticker.NewDaxxerTicker()
	pricesApiService  = service.NewApiService()
	accountController = controller.NewAccountController()
	currencyController       currencies.CurrencyHandler = currencies.NewCurrencyController()
	cryptoCurrencyController currencies.CurrencyHandler = currencies.NewCryptoCurrencyController()
)

func main() {
	dbConfig.Init()

	daxxerTicker.Run(1, pricesApiService.Update)

	httpRouter.POST("/create", accountController.Create)
	httpRouter.POST("/deposit", accountController.Deposit)
	httpRouter.POST("/withdraw", accountController.Withdraw)

	httpRouter.POST("/currency", currencyController.Upsert)
	httpRouter.PUT("/currency", currencyController.Upsert)
	httpRouter.DELETE("/currency", currencyController.Delete)
	httpRouter.GET("/currency", currencyController.GetById)
	httpRouter.GET("/currencies", currencyController.GetAll)

	httpRouter.POST("/crypto-currency", cryptoCurrencyController.Upsert)
	httpRouter.PUT("/crypto-currency", cryptoCurrencyController.Upsert)
	httpRouter.DELETE("/crypto-currency", cryptoCurrencyController.Delete)
	httpRouter.GET("/crypto-currency", cryptoCurrencyController.GetById)
	httpRouter.GET("/crypto-currencies", cryptoCurrencyController.GetAll)

	httpRouter.SERVE(":8000")
}