package storage

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/sawdustofmind/dataflow-shop-statistics/internal/entity"

	"github.com/stretchr/testify/require"
)

func TestNewInMemoryStorage(t *testing.T) {
	s := NewInMemoryStorage(10)

	// put some data
	err := s.PutData(entity.SalesData{
		StoreId:      "1",
		ProductId:    "1",
		QuantitySold: 10,
		SalePrice:    decimal.NewFromFloat(3.3),
		SaleDate:     tm("2024-06-15T14:30:00Z"),
	})
	require.NoError(t, err)

	// same time
	err = s.PutData(entity.SalesData{
		StoreId:      "1",
		ProductId:    "2",
		QuantitySold: 5,
		SalePrice:    decimal.NewFromFloat(4),
		SaleDate:     tm("2024-06-15T14:30:00Z"),
	})
	require.NoError(t, err)

	// different time
	err = s.PutData(entity.SalesData{
		StoreId:      "1",
		ProductId:    "1",
		QuantitySold: 7,
		SalePrice:    decimal.NewFromFloat(5),
		SaleDate:     tm("2024-06-16T14:30:00Z"),
	})
	require.NoError(t, err)

	data, err := s.GetData()
	require.NoError(t, err)

	assert.Equal(t, []entity.SalesData{
		{
			StoreId:      "1",
			ProductId:    "1",
			QuantitySold: 10,
			SalePrice:    decimal.NewFromFloat(3.3),
			SaleDate:     tm("2024-06-15T14:30:00Z"),
		},
		{
			StoreId:      "1",
			ProductId:    "2",
			QuantitySold: 5,
			SalePrice:    decimal.NewFromFloat(4),
			SaleDate:     tm("2024-06-15T14:30:00Z"),
		},
		{
			StoreId:      "1",
			ProductId:    "1",
			QuantitySold: 7,
			SalePrice:    decimal.NewFromFloat(5),
			SaleDate:     tm("2024-06-16T14:30:00Z"),
		},
	}, data)

	// test different ranges
	sum, err := s.SalesSum("1", tm("2024-06-15T14:30:00Z"), tm("2024-06-15T14:30:00Z"))
	require.NoError(t, err)
	assert.Equal(t, tm("2024-06-15T14:30:00Z").UnixNano(), sum.From)
	assert.Equal(t, tm("2024-06-15T14:30:00Z").UnixNano(), sum.To)
	assert.Equal(t, decimal.NewFromFloat(53).String(), sum.Sum.String())

	sum, err = s.SalesSum("1", tm("2024-06-14T14:30:00Z"), tm("2024-06-15T15:30:00Z"))
	require.NoError(t, err)
	assert.Equal(t, tm("2024-06-15T14:30:00Z").UnixNano(), sum.From)
	assert.Equal(t, tm("2024-06-15T14:30:00Z").UnixNano(), sum.To)
	assert.Equal(t, decimal.NewFromFloat(53).String(), sum.Sum.String())

	sum, err = s.SalesSum("1", tm("2024-06-14T14:30:00Z"), tm("2024-06-16T15:30:00Z"))
	require.NoError(t, err)
	assert.Equal(t, tm("2024-06-15T14:30:00Z").UnixNano(), sum.From)
	assert.Equal(t, tm("2024-06-16T14:30:00Z").UnixNano(), sum.To)
	assert.Equal(t, decimal.NewFromFloat(88).String(), sum.Sum.String())

	sum, err = s.SalesSum("1", tm("2024-06-15T14:30:01Z"), tm("2024-06-16T15:30:00Z"))
	require.NoError(t, err)
	assert.Equal(t, tm("2024-06-16T14:30:00Z").UnixNano(), sum.From)
	assert.Equal(t, tm("2024-06-16T14:30:00Z").UnixNano(), sum.To)
	assert.Equal(t, decimal.NewFromFloat(35).String(), sum.Sum.String())
}

func BenchmarkInCalcSum(b *testing.B) {
	s := NewInMemoryStorage(10)

	initDate := tm("2024-06-15T14:30:00Z")

	for i := 0; i < 100000; i++ {
		err := s.PutData(entity.SalesData{
			StoreId:      "1",
			ProductId:    "1",
			QuantitySold: 10,
			SalePrice:    decimal.NewFromFloat(3.3),
			SaleDate:     initDate.Add(time.Duration(i) * time.Second),
		})
		require.NoError(b, err)
	}

	from := initDate.Add(time.Duration(20000) * time.Second)
	to := initDate.Add(time.Duration(21100) * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.SalesSum("1", from, to)
	}
}

func tm(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
