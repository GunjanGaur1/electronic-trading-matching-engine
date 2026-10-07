package order

type Side string
type OrderType string

const (
	Buy  Side = "BUY"
	Sell Side = "SELL"

	Limit  OrderType = "LIMIT"
	Market OrderType = "MARKET"
)

type Order struct {
	ID        int
	Side      Side
	Type      OrderType
	Price     float64
	Quantity  int
	Timestamp int64
}
