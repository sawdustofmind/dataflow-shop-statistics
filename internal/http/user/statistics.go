package userhttp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sawdustofmind/dataflow-shop-statistics/internal/entity"
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/http/dto"
	api_types "github.com/sawdustofmind/dataflow-shop-statistics/internal/openapi"
)

func (s ServerImpl) PostV1Calculate(c *gin.Context) {
	req := &api_types.PostV1CalculateJSONRequestBody{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewAPIError(err))
		return
	}

	result, err := s.statisticsService.Calculate(entity.StatisticsRequest{
		EndDate:   req.EndDate,
		Operation: req.Operation,
		StartDate: req.StartDate,
		StoreId:   req.StoreId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewAPIError(err))
		return
	}

	response := api_types.StatisticsResponse{
		EndDate:    time.Unix(0, result.From).In(time.UTC),
		StartDate:  time.Unix(0, result.To).In(time.UTC),
		StoreId:    req.StoreId,
		TotalSales: json.Number(result.Sum.String()),
	}

	c.JSON(http.StatusOK, response)
}
