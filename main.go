package main

import (
	"github.com/julioisaac/daxxer-api/routers"
	"github.com/julioisaac/daxxer-api/routers/gin"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	mongodb2 "github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	"github.com/julioisaac/daxxer-api/src/helpers/ticker"
	"github.com/julioisaac/daxxer-api/src/wallet/account/controller"
	currencies "github.com/julioisaac/daxxer-api/src/wallet/currencies/controller"
	api "github.com/julioisaac/daxxer-api/src/wallet/prices/repository"
	"github.com/julioisaac/daxxer-api/src/wallet/prices/service"
	"github.com/julioisaac/daxxer-api/storage"
	"github.com/julioisaac/daxxer-api/storage/mongodb"
	"net/http"
	"time"
)

var (
	dbConfig   storage.DBConfig = mongodb.NewMongoConfig()
	httpRouter routers.Router   = gin.NewGinRouter()
	daxxerTicker      = ticker.NewDaxxerTicker()

	//db name and collections config or env
	cryptoRepo repository.DBRepository  = mongodb2.NewMongodbRepository("daxxer", "cryptoCurrencies")
	currencyRepo repository.DBRepository = mongodb2.NewMongodbRepository("daxxer", "currencies")
	pricesRepo  repository.DBRepository = mongodb2.NewMongodbRepository("daxxer", "prices")
	// url and timeout in config or env
	apiRepo  api.ApiRepository = api.NewCoinGeckoApiRepo("https://api.coingecko.com/api/v3/simple/price", &http.Client{Timeout: 5 * time.Second})
	pricesApiService  = service.NewApiService(cryptoRepo, currencyRepo, pricesRepo, apiRepo)

	accountController = controller.NewAccountController()
	currencyController       currencies.CurrencyHandler = currencies.NewCurrencyController()
	cryptoCurrencyController currencies.CurrencyHandler = currencies.NewCryptoCurrencyController()
)

func main() {
	dbConfig.Init()

	//update prices interval config or env
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

	//app port config or env
	httpRouter.SERVE(":8000")
}