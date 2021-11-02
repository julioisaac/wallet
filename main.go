package main

import (
	"context"
	"github.com/julioisaac/daxxer-api/internal/database"
	"github.com/julioisaac/daxxer-api/internal/database/mongodb"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/julioisaac/daxxer-api/metrics"
	"github.com/julioisaac/daxxer-api/routers"
	"github.com/julioisaac/daxxer-api/routers/gin"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	mongodb2 "github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	"github.com/julioisaac/daxxer-api/src/helpers/ticker"
	pkg "github.com/julioisaac/daxxer-api/src/wallet"
	"github.com/julioisaac/daxxer-api/src/wallet/account/controller"
	currencies "github.com/julioisaac/daxxer-api/src/wallet/currencies/controller"
	controller2 "github.com/julioisaac/daxxer-api/src/wallet/prices/controller"
	api "github.com/julioisaac/daxxer-api/src/wallet/prices/repository"
	"github.com/julioisaac/daxxer-api/src/wallet/prices/service"
)

var (
	logConfig    logs.ILog         = logs.NewZapLogger()
	dbConfig     database.DBConfig = mongodb.NewMongoConfig()
	httpRouter   routers.Router    = gin.NewGinRouter()
	daxxerTicker                   = ticker.NewDaxxerTicker()

	cryptoRepo   repository.DBRepository = mongodb2.NewMongodbRepository("daxxer", "crypto_currencies")
	currencyRepo repository.DBRepository = mongodb2.NewMongodbRepository("daxxer", "currencies")
	pricesRepo   repository.DBRepository = mongodb2.NewMongodbRepository("daxxer", "prices")
	apiRepo            api.ApiRepository = api.NewCoinGeckoApiRepo(metrics.Metric().GetClient())
	pricesApiService                   	 = service.NewApiService(cryptoRepo, currencyRepo, pricesRepo, apiRepo)

	healthCheck                                         = pkg.NewHealthCheck()
	pricesController                                    = controller2.NewPricesController()
	accountController                                   = controller.NewAccountController()
	currencyController       currencies.CurrencyHandler = currencies.NewCurrencyController()
	cryptoCurrencyController currencies.CurrencyHandler = currencies.NewCryptoCurrencyController()
)

func main() {
	logConfig.Init()
	httpRouter.SetupTracer()
	dbConfig.Init()

	daxxerTicker.Run(context.TODO(), 1, pricesApiService.Update)

	httpRouter.GET("/health-check", healthCheck.IsAlive)

	httpRouter.POST("/create", accountController.Create)
	httpRouter.POST("/deposit", accountController.Deposit)
	httpRouter.POST("/withdraw", accountController.Withdraw)
	httpRouter.GET("/history", accountController.History)
	httpRouter.GET("/balance", accountController.Balance)

	httpRouter.GET("/prices", pricesController.GetAll)

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