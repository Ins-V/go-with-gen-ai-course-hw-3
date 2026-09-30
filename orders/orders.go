// Package orders implements Section 2 of the homework: OrderService is
// refactored so that its only database dependency is the minimal
// OrderStore interface, injected through the constructor.
//
// The interface is declared here, in the consumer package, not next to a
// concrete database implementation.
package orders

// OrderStore is the minimal database dependency of OrderService.
//
// It holds exactly one method because PlaceOrder — the only behaviour of
// the service — issues a single INSERT and needs nothing but the error it
// returns; every other capability of *sql.DB (Query, QueryRow, Begin,
// Prepare, Close, PingContext, ...) is unused by this package and would
// only weaken the abstraction.
type OrderStore interface {
	// Exec runs a write query with the given arguments and reports failure.
	Exec(query string, args ...any) error
}

// OrderService places orders through an injected OrderStore.
type OrderService struct {
	store OrderStore
}

// NewOrderService returns an OrderService backed by store.
func NewOrderService(store OrderStore) *OrderService {
	return &OrderService{store: store}
}

// PlaceOrder inserts the order identified by orderID with the given
// amount. It returns the store error unchanged, or nil on success.
func (s *OrderService) PlaceOrder(orderID string, amount float64) error {
	return s.store.Exec(
		"INSERT INTO orders (id, amount) VALUES (?, ?)",
		orderID, amount,
	)
}
