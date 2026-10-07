package book

import (
	"fmt"
	"sort"
	"sync"

	"github.com/GunjanGaur1/electronic-trading-matching-engine/order"
)

type Trade struct {
	BuyOrderID  int
	SellOrderID int
	Price       float64
	Quantity    int
}

type OrderBook struct {
	Bids []*order.Order
	Asks []*order.Order

	// Fast order-ID lookup.
	Orders map[int]*order.Order

	// Protects shared order-book state.
	mu sync.RWMutex

	nextTimestamp int64
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids:   make([]*order.Order, 0),
		Asks:   make([]*order.Order, 0),
		Orders: make(map[int]*order.Order),
	}
}

// AddOrder inserts an order and attempts to match it.
func (ob *OrderBook) AddOrder(o *order.Order) []Trade {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	ob.nextTimestamp++
	o.Timestamp = ob.nextTimestamp

	ob.Orders[o.ID] = o

	trades := ob.matchOrder(o)

	// Only remaining LIMIT orders rest on the book.
	if o.Quantity > 0 && o.Type == order.Limit {
		ob.addToBook(o)
	}

	return trades
}

func (ob *OrderBook) addToBook(o *order.Order) {
	if o.Side == order.Buy {
		ob.Bids = append(ob.Bids, o)

		// Highest bid first.
		sort.SliceStable(ob.Bids, func(i, j int) bool {
			return ob.Bids[i].Price > ob.Bids[j].Price
		})
	} else {
		ob.Asks = append(ob.Asks, o)

		// Lowest ask first.
		sort.SliceStable(ob.Asks, func(i, j int) bool {
			return ob.Asks[i].Price < ob.Asks[j].Price
		})
	}
}

func (ob *OrderBook) matchOrder(incoming *order.Order) []Trade {
	var trades []Trade

	for incoming.Quantity > 0 {

		var resting *order.Order

		if incoming.Side == order.Buy {

			if len(ob.Asks) == 0 {
				break
			}

			resting = ob.Asks[0]

			// LIMIT BUY must be >= best ask.
			if incoming.Type == order.Limit &&
				incoming.Price < resting.Price {
				break
			}

		} else {

			if len(ob.Bids) == 0 {
				break
			}

			resting = ob.Bids[0]

			// LIMIT SELL must be <= best bid.
			if incoming.Type == order.Limit &&
				incoming.Price > resting.Price {
				break
			}
		}

		tradeQuantity := incoming.Quantity

		if resting.Quantity < tradeQuantity {
			tradeQuantity = resting.Quantity
		}

		trade := Trade{
			Price:    resting.Price,
			Quantity: tradeQuantity,
		}

		if incoming.Side == order.Buy {
			trade.BuyOrderID = incoming.ID
			trade.SellOrderID = resting.ID
		} else {
			trade.BuyOrderID = resting.ID
			trade.SellOrderID = incoming.ID
		}

		trades = append(trades, trade)

		incoming.Quantity -= tradeQuantity
		resting.Quantity -= tradeQuantity

		if resting.Quantity == 0 {
			ob.removeRestingOrder(resting)
		}
	}

	return trades
}

func (ob *OrderBook) removeRestingOrder(o *order.Order) {
	if o.Side == order.Buy {
		for i, current := range ob.Bids {
			if current.ID == o.ID {
				ob.Bids = append(ob.Bids[:i], ob.Bids[i+1:]...)
				break
			}
		}
	} else {
		for i, current := range ob.Asks {
			if current.ID == o.ID {
				ob.Asks = append(ob.Asks[:i], ob.Asks[i+1:]...)
				break
			}
		}
	}

	delete(ob.Orders, o.ID)
}

func (ob *OrderBook) CancelOrder(orderID int) bool {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	o, exists := ob.Orders[orderID]

	if !exists {
		return false
	}

	if o.Side == order.Buy {
		for i, current := range ob.Bids {
			if current.ID == orderID {
				ob.Bids = append(ob.Bids[:i], ob.Bids[i+1:]...)
				delete(ob.Orders, orderID)
				return true
			}
		}
	} else {
		for i, current := range ob.Asks {
			if current.ID == orderID {
				ob.Asks = append(ob.Asks[:i], ob.Asks[i+1:]...)
				delete(ob.Orders, orderID)
				return true
			}
		}
	}

	return false
}

func (ob *OrderBook) PrintBook() {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	fmt.Println("\n========== ORDER BOOK ==========")

	fmt.Println("\nASKS:")
	for _, o := range ob.Asks {
		fmt.Printf(
			"Order %d | SELL | $%.2f | Qty: %d\n",
			o.ID,
			o.Price,
			o.Quantity,
		)
	}

	fmt.Println("\nBIDS:")
	for _, o := range ob.Bids {
		fmt.Printf(
			"Order %d | BUY | $%.2f | Qty: %d\n",
			o.ID,
			o.Price,
			o.Quantity,
		)
	}

	fmt.Println("================================")
}
