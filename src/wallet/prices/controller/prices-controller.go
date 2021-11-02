package controller

import (
	"encoding/json"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	pricesService "github.com/julioisaac/daxxer-api/src/wallet/prices/service"
	"go.uber.org/zap"
	"net/http"
)

var (
	priceRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "prices")
	priceService = pricesService.NewPricesService(priceRepo)
)

type pricesController struct {}

func NewPricesController() *pricesController {
	return &pricesController{}
}
func (c *pricesController) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	prices, err := priceService.GetAll(request.Context())
	if err != nil {
		
	}
	logs.Instance.Log.Debug(request.Context(), "prices request success")
	response.WriteHeader(http.StatusOK)
	err = json.NewEncoder(response).Encode(prices)
	if err != nil {
		logs.Instance.Log.Error(request.Context(), "error when trying to encode prices", zap.Error(err))
	}
}