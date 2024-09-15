package entity

import "time"

type StatisticsRequest struct {
	// EndDate starting date for statistics calculation
	EndDate time.Time `json:"end_date"`

	// Operation operation name. possible values: `total_sales`
	Operation string `json:"operation"`

	// StartDate starting date for statistics calculation
	StartDate time.Time `json:"start_date"`

	// StoreId store id
	StoreId string `json:"store_id"`
}
