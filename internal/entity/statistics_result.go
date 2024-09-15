package entity

import (
	"github.com/shopspring/decimal"
)

type StatisticsResult struct {
	From int64           `json:"from"`
	To   int64           `json:"to"`
	Sum  decimal.Decimal `json:"sum"`
}
