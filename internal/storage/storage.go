package storage

import (
	"time"

	"github.com/sawdustofmind/dataflow-shop-statistics/internal/entity"
)

type Storage interface {
	PutData(data entity.SalesData) error
	GetData() ([]entity.SalesData, error)

	SalesSum(storeId string, since time.Time, to time.Time) (entity.StatisticsResult, error)
}
