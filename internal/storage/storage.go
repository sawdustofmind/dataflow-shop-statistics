package storage

import (
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/entity"
)

type Storage interface {
	PutData(data entity.SalesData) error
	GetData() ([]entity.SalesData, error)

	SalesSum(storeId string, since int64, to int64) (entity.StatisticsResult, error)
}
