package userhttp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/sawdustofmind/dataflow-shop-statistics/internal/entity"
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/http/dto"
	api_types "github.com/sawdustofmind/dataflow-shop-statistics/internal/openapi"
)

func (s ServerImpl) GetV1Data(c *gin.Context) {
	data, err := s.dataService.GetData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewAPIError(err))
		return
	}

	response := make([]api_types.Data, 0, len(data))
	for _, d := range data {
		response = append(response, api_types.Data{
			ProductId:    d.ProductId,
			QuantitySold: d.QuantitySold,
			SaleDate:     time.Unix(0, d.SaleDate).In(time.UTC),
			SalePrice:    json.Number(d.SalePrice.String()),
			StoreId:      d.StoreId,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (s ServerImpl) PostV1Data(c *gin.Context) {
	req := &api_types.PostV1DataJSONRequestBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewAPIError(err))
		return
	}

	price, err := decimal.NewFromString(req.SalePrice.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewAPIError(err))
		return
	}

	sd := entity.SalesData{
		StoreId:      req.StoreId,
		ProductId:    req.ProductId,
		QuantitySold: req.QuantitySold,
		SalePrice:    price,
		SaleDate:     req.SaleDate.UnixNano(),
	}

	err = s.dataService.PutData(sd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewAPIError(err))
		return
	}

	c.JSON(http.StatusOK, dto.APISuccessStatus)
}
