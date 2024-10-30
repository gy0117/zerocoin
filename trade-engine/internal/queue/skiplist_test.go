package queue

import (
	"github.com/shopspring/decimal"
	"testing"
)

type order struct {
	Id       string
	Quantity decimal.Decimal
}

func (o *order) GetId() string {
	return o.Id
}

var _ NodeValue = (*order)(nil)

func TestSkipList_Insert(t *testing.T) {
	skipList, _ := NewSkipList()

	for i := 0; i < 100000; i++ {
		o := &order{
			Id:       "123",
			Quantity: decimal.Decimal{},
		}
		skipList.Insert(decimal.NewFromInt(int64(i)), o)
	}
	t.Logf("size: %d, level: %d\n", skipList.size, skipList.level)
}
