package service

import (
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/entity"
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/storage"
)

type DataService struct {
	storage storage.Storage
}

func NewDataService(storage storage.Storage) *DataService {
	return &DataService{
		storage: storage,
	}
}

func (ds *DataService) PutData(data entity.SalesData) error {
	return ds.storage.PutData(data)
}

func (ds *DataService) GetData() ([]entity.SalesData, error) {
	return ds.storage.GetData()
}
