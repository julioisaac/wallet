package controller

import (
	"encoding/json"
	"github.com/julioisaac/daxxer-api/src/helpers/repository"
	"github.com/julioisaac/daxxer-api/src/helpers/repository/mongodb"
	service2 "github.com/julioisaac/daxxer-api/src/wallet/prices/service"
	"net/http"
)

var (
	priceRepo repository.DBRepository = mongodb.NewMongodbRepository("daxxer", "prices")
	priceService = service2.NewPricesService(priceRepo)
)

type pricesController struct {}

func NewPricesController() *pricesController {
	return &pricesController{}
}
func (c *pricesController) GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	prices, _ := priceService.GetAll()
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(prices)
}