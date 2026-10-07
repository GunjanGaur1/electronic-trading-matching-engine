package tests

import (
	"sync"
	"testing"

	"github.com/GunjanGaur1/electronic-trading-matching-engine/book"
	"github.com/GunjanGaur1/electronic-trading-matching-engine/order"
)

func TestConcurrentOrders(t *testing.T) {
	ob := book.NewOrderBook()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			ob.AddOrder(&order.Order{
				ID:       id,
				Side:     order.Buy,
				Type:     order.Limit,
				Price:    float64(100 + (id % 5)),
				Quantity: 10,
			})
		}(i)
	}

	wg.Wait()

	if len(ob.Bids) != 100 {
		t.Fatalf(
			"expected 100 orders, got %d",
			len(ob.Bids),
		)
	}
}
