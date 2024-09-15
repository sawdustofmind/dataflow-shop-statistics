package service

import (
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/entity"
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/storage"
)

type StatisticsService struct {
	storage storage.Storage
}

func NewStatisticsService(storage storage.Storage) *StatisticsService {
	return &StatisticsService{
		storage: storage,
	}
}

func (ds *StatisticsService) Calculate(sr entity.StatisticsRequest) (entity.StatisticsResult, error) {
	return ds.storage.SalesSum(sr.StoreId, sr.StartDate, sr.EndDate)
}
