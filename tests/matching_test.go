package tests

import (
	"testing"

	"github.com/GunjanGaur1/electronic-trading-matching-engine/book"
	"github.com/GunjanGaur1/electronic-trading-matching-engine/order"
)

func TestLimitOrderMatching(t *testing.T) {
	ob := book.NewOrderBook()

	ob.AddOrder(&order.Order{
		ID:       1,
		Side:     order.Sell,
		Type:     order.Limit,
		Price:    100,
		Quantity: 50,
	})

	trades := ob.AddOrder(&order.Order{
		ID:       2,
		Side:     order.Buy,
		Type:     order.Limit,
		Price:    100,
		Quantity: 30,
	})

	if len(trades) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(trades))
	}

	if trades[0].Quantity != 30 {
		t.Fatalf(
			"expected quantity 30, got %d",
			trades[0].Quantity,
		)
	}

	if len(ob.Asks) != 1 {
		t.Fatalf("expected remaining sell order")
	}

	if ob.Asks[0].Quantity != 20 {
		t.Fatalf(
			"expected remaining quantity 20, got %d",
			ob.Asks[0].Quantity,
		)
	}
}

func TestFIFO(t *testing.T) {
	ob := book.NewOrderBook()

	ob.AddOrder(&order.Order{
		ID:       1,
		Side:     order.Buy,
		Type:     order.Limit,
		Price:    100,
		Quantity: 20,
	})

	ob.AddOrder(&order.Order{
		ID:       2,
		Side:     order.Buy,
		Type:     order.Limit,
		Price:    100,
		Quantity: 30,
	})

	trades := ob.AddOrder(&order.Order{
		ID:       3,
		Side:     order.Sell,
		Type:     order.Limit,
		Price:    100,
		Quantity: 40,
	})

	if len(trades) != 2 {
		t.Fatalf("expected 2 trades, got %d", len(trades))
	}

	if trades[0].BuyOrderID != 1 {
		t.Fatalf("expected Order 1 to execute first")
	}

	if trades[1].BuyOrderID != 2 {
		t.Fatalf("expected Order 2 to execute second")
	}
}

func TestCancel(t *testing.T) {
	ob := book.NewOrderBook()

	ob.AddOrder(&order.Order{
		ID:       1,
		Side:     order.Buy,
		Type:     order.Limit,
		Price:    100,
		Quantity: 10,
	})

	if !ob.CancelOrder(1) {
		t.Fatal("expected cancellation to succeed")
	}

	if ob.CancelOrder(1) {
		t.Fatal("expected second cancellation to fail")
	}
}
