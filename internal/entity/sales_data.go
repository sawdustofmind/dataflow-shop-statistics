package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// Data defines model for Data.
type SalesData struct {
	StoreId      string          `json:"store_id"`
	ProductId    string          `json:"product_id"`
	QuantitySold uint64          `json:"quantity_sold"`
	SalePrice    decimal.Decimal `json:"sale_price"`
	SaleDate     time.Time       `json:"sale_date"`
}
