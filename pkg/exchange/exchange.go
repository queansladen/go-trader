// Package exchange defines the core interfaces and types for interacting
// with trading exchanges. All exchange implementations must satisfy the
// Exchange interface defined here.
package exchange

import (
	"context"
	"time"
)

// Side represents the side of an order (buy or sell).
type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)

// OrderType represents the type of an order.
type OrderType string

const (
	OrderTypeMarket OrderType = "market"
	OrderTypeLimit  OrderType = "limit"
)

// OrderStatus represents the current status of an order.
type OrderStatus string

const (
	OrderStatusOpen      OrderStatus = "open"
	OrderStatusFilled    OrderStatus = "filled"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusPartial   OrderStatus = "partial"
)

// Order represents a single trade order.
type Order struct {
	ID        string
	Symbol    string
	Side      Side
	Type      OrderType
	Status    OrderStatus
	Price     float64
	Quantity  float64
	Filled    float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Ticker holds the latest price information for a symbol.
type Ticker struct {
	Symbol    string
	Bid       float64
	Ask       float64
	Last      float64
	Volume    float64
	Timestamp time.Time
}

// Balance represents the available and total balance for a currency.
type Balance struct {
	Currency  string
	Available float64
	Total     float64
}

// PlaceOrderRequest contains the parameters needed to place a new order.
type PlaceOrderRequest struct {
	Symbol   string
	Side     Side
	Type     OrderType
	Price    float64 // ignored for market orders
	Quantity float64
}

// Exchange defines the interface that all exchange implementations must satisfy.
type Exchange interface {
	// Name returns the human-readable name of the exchange.
	Name() string

	// GetTicker returns the latest ticker data for the given symbol.
	GetTicker(ctx context.Context, symbol string) (*Ticker, error)

	// GetBalances returns all non-zero balances for the authenticated account.
	GetBalances(ctx context.Context) ([]Balance, error)

	// PlaceOrder submits a new order to the exchange.
	PlaceOrder(ctx context.Context, req PlaceOrderRequest) (*Order, error)

	// CancelOrder cancels an open order by its ID.
	CancelOrder(ctx context.Context, orderID string) error

	// GetOrder retrieves the current state of an order by its ID.
	GetOrder(ctx context.Context, orderID string) (*Order, error)

	// GetOpenOrders returns all currently open orders, optionally filtered by symbol.
	GetOpenOrders(ctx context.Context, symbol string) ([]Order, error)
}
