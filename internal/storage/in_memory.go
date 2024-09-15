package storage

import (
	"math"
	"sync"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/sawdustofmind/dataflow-shop-statistics/internal/entity"
	"github.com/sawdustofmind/dataflow-shop-statistics/internal/log"
)

type inMemoryStorage struct {
	mu sync.RWMutex

	Data []entity.SalesData

	StoreIndex map[string][]entity.SalesData
}

func NewInMemoryStorage(capacity int) Storage {
	return &inMemoryStorage{
		Data:       make([]entity.SalesData, 0, capacity),
		StoreIndex: make(map[string][]entity.SalesData),
	}
}

func (s *inMemoryStorage) PutData(data entity.SalesData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Data = append(s.Data, data)
	s.StoreIndex[data.StoreId] = append(s.StoreIndex[data.StoreId], data)

	return nil
}

func (s *inMemoryStorage) GetData() ([]entity.SalesData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cpy := make([]entity.SalesData, len(s.Data))
	copy(cpy, s.Data)

	return cpy, nil
}

func (s *inMemoryStorage) SalesSum(storeId string, since int64, to int64) (entity.StatisticsResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storeIndex := s.StoreIndex[storeId]

	log.Info("ss", zap.Int64("since", since), zap.Int64("to", to), zap.Int64("cur", s.Data[0].SaleDate))

	sum := decimal.New(0, 8) // TODO: come up with something better

	minDate := int64(math.MaxInt64)
	maxDate := int64(math.MinInt64)

	for _, data := range storeIndex {
		if data.SaleDate >= since && data.SaleDate <= to {
			if data.SaleDate <= minDate {
				minDate = data.SaleDate
			}
			if data.SaleDate >= maxDate {
				maxDate = data.SaleDate
			}

			sum = sum.Add(data.SalePrice)
		}
	}

	result := entity.StatisticsResult{
		Sum: sum,
	}

	if !sum.IsZero() {
		result.From = minDate
		result.To = maxDate
	}

	return result, nil
}
