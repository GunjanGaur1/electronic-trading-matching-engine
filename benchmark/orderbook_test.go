package benchmark

import (
	"testing"

	"github.com/GunjanGaur1/electronic-trading-matching-engine/book"
	"github.com/GunjanGaur1/electronic-trading-matching-engine/order"
)

func BenchmarkOrderInsertion(b *testing.B) {

	for i := 0; i < b.N; i++ {

		ob := book.NewOrderBook()

		ob.AddOrder(&order.Order{
			ID:       i,
			Side:     order.Buy,
			Type:     order.Limit,
			Price:    100,
			Quantity: 10,
		})
	}
}
