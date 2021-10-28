package wallet

import (
	"github.com/julioisaac/daxxer-api/routers"
	"github.com/julioisaac/daxxer-api/routers/gin"
	"github.com/julioisaac/daxxer-api/src/wallet/account/controller"
	"github.com/julioisaac/daxxer-api/storage"
	"github.com/julioisaac/daxxer-api/storage/mongodb"
)

var (
	dbConfig   storage.DBConfig = mongodb.NewMongoConfig()
	httpRouter routers.Router   = gin.NewGinRouter()

	accountController = controller.NewAccountController()
)

func main() {
	dbConfig.Init()

	httpRouter.POST("/create", accountController.Create)
	httpRouter.POST("/deposit", accountController.Deposit)
	httpRouter.POST("/withdraw", accountController.Withdraw)

	httpRouter.SERVE(":8000")
}