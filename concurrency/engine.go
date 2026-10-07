package concurrency

import (
	"sync"

	"github.com/GunjanGaur1/electronic-trading-matching-engine/book"
	"github.com/GunjanGaur1/electronic-trading-matching-engine/order"
)

type Engine struct {
	Book *book.OrderBook
}

func NewEngine() *Engine {
	return &Engine{
		Book: book.NewOrderBook(),
	}
}

func (e *Engine) SubmitOrders(orders []*order.Order) {
	var wg sync.WaitGroup

	for _, o := range orders {
		wg.Add(1)

		go func(ord *order.Order) {
			defer wg.Done()

			e.Book.AddOrder(ord)
		}(o)
	}

	wg.Wait()
}
