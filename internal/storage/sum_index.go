package storage

import (
	"github.com/igrmk/treemap/v2"
	"github.com/shopspring/decimal"
)

type sumIndex struct {
	m *treemap.TreeMap[int64, decimal.Decimal]
}

func newSumIndex() *sumIndex {
	return &sumIndex{
		m: treemap.New[int64, decimal.Decimal](),
	}
}

func (idx *sumIndex) AddSum(date int64, value decimal.Decimal) {
	v, ok := idx.m.Get(date)
	if ok {
		idx.m.Set(date, v.Add(value))
		return
	}
	idx.m.Set(date, value)
}

func (idx *sumIndex) CalcSum(since int64, to int64) (int64, int64, decimal.Decimal) {
	rangeSince := int64(0)
	rangeTo := int64(0)

	sum := decimal.Zero
	for it := idx.m.LowerBound(since); it.Valid(); it.Next() {
		if it.Key() > to {
			break
		}

		if rangeSince == 0 {
			rangeSince = it.Key()
		}
		rangeTo = it.Key()

		sum = sum.Add(it.Value())
	}

	return rangeSince, rangeTo, sum
}
