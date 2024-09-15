package storage

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/sawdustofmind/dataflow-shop-statistics/internal/entity"
)

type inMemoryStorage struct {
	mu sync.RWMutex

	data       []entity.SalesData
	sumIndexes map[string]*sumIndex
}

func NewInMemoryStorage(capacity int) Storage {
	return &inMemoryStorage{
		data:       make([]entity.SalesData, 0, capacity),
		sumIndexes: make(map[string]*sumIndex),
	}
}

func (s *inMemoryStorage) PutData(data entity.SalesData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, data)

	// add to sum index
	si, ok := s.sumIndexes[data.StoreId]
	if !ok {
		si = newSumIndex()
		s.sumIndexes[data.StoreId] = si
	}
	total := data.SalePrice.Mul(decimal.NewFromUint64(data.QuantitySold))
	si.AddSum(data.SaleDate.UnixNano(), total)

	return nil
}

func (s *inMemoryStorage) GetData() ([]entity.SalesData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cpy := make([]entity.SalesData, len(s.data))
	copy(cpy, s.data)

	return cpy, nil
}

func (s *inMemoryStorage) SalesSum(storeId string, since time.Time, to time.Time) (entity.StatisticsResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	si, ok := s.sumIndexes[storeId]
	if !ok {
		return entity.StatisticsResult{}, nil
	}

	rangeFrom, rangeTo, sum := si.CalcSum(since.UnixNano(), to.UnixNano())

	result := entity.StatisticsResult{
		Sum: sum,
	}

	if !sum.IsZero() {
		result.From = rangeFrom
		result.To = rangeTo
	}

	return result, nil
}
