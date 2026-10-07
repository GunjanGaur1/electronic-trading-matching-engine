package main

import (
	"fmt"

	"github.com/GunjanGaur1/electronic-trading-matching-engine/book"
	"github.com/GunjanGaur1/electronic-trading-matching-engine/concurrency"
	"github.com/GunjanGaur1/electronic-trading-matching-engine/order"
)

func main() {

	ob := book.NewOrderBook()

	// ---------------------------------------------
	// LIMIT ORDERS
	// ---------------------------------------------

	ob.AddOrder(&order.Order{
		ID:       1,
		Side:     order.Sell,
		Type:     order.Limit,
		Price:    101,
		Quantity: 50,
	})

	ob.AddOrder(&order.Order{
		ID:       2,
		Side:     order.Sell,
		Type:     order.Limit,
		Price:    102,
		Quantity: 100,
	})

	ob.AddOrder(&order.Order{
		ID:       3,
		Side:     order.Buy,
		Type:     order.Limit,
		Price:    99,
		Quantity: 100,
	})

	ob.PrintBook()

	// ---------------------------------------------
	// MARKET ORDER
	// ---------------------------------------------

	fmt.Println("\n--- MARKET BUY ---")

	ob.AddOrder(&order.Order{
		ID:       4,
		Side:     order.Buy,
		Type:     order.Market,
		Quantity: 70,
	})

	ob.PrintBook()

	// ---------------------------------------------
	// PARTIAL FILL
	// ---------------------------------------------

	fmt.Println("\n--- PARTIAL FILL ---")

	ob.AddOrder(&order.Order{
		ID:       5,
		Side:     order.Buy,
		Type:     order.Limit,
		Price:    102,
		Quantity: 150,
	})

	ob.PrintBook()

	// ---------------------------------------------
	// CANCEL
	// ---------------------------------------------

	fmt.Println("\n--- CANCEL ---")

	fmt.Println(
		"Cancelled:",
		ob.CancelOrder(3),
	)

	ob.PrintBook()

	// ---------------------------------------------
	// FIFO
	// ---------------------------------------------

	fmt.Println("\n--- FIFO ---")

	ob.AddOrder(&order.Order{
		ID:       6,
		Side:     order.Buy,
		Type:     order.Limit,
		Price:    100,
		Quantity: 20,
	})

	ob.AddOrder(&order.Order{
		ID:       7,
		Side:     order.Buy,
		Type:     order.Limit,
		Price:    100,
		Quantity: 30,
	})

	ob.AddOrder(&order.Order{
		ID:       8,
		Side:     order.Sell,
		Type:     order.Limit,
		Price:    100,
		Quantity: 40,
	})

	ob.PrintBook()

	// ---------------------------------------------
	// CONCURRENT PROCESSING
	// ---------------------------------------------

	fmt.Println("\n--- CONCURRENT PROCESSING ---")

	engine := concurrency.NewEngine()

	orders := []*order.Order{
		{
			ID:       100,
			Side:     order.Buy,
			Type:     order.Limit,
			Price:    99,
			Quantity: 10,
		},
		{
			ID:       101,
			Side:     order.Buy,
			Type:     order.Limit,
			Price:    98,
			Quantity: 20,
		},
		{
			ID:       102,
			Side:     order.Sell,
			Type:     order.Limit,
			Price:    105,
			Quantity: 15,
		},
	}

	engine.SubmitOrders(orders)

	engine.Book.PrintBook()
}
