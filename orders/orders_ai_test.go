package orders

import (
	"errors"
	"testing"
)

// mockOrderStore is a hand-written fake OrderStore that records how it was
// called, so tests can run without a real database.
type mockOrderStore struct {
	execCalls int
	lastQuery string
	lastArgs  []any
	execErr   error
}

// compile-time check that the mock really satisfies the interface.
var _ OrderStore = (*mockOrderStore)(nil)

// Exec records the call and returns the preconfigured error.
func (m *mockOrderStore) Exec(query string, args ...any) error {
	m.execCalls++
	m.lastQuery = query
	m.lastArgs = args
	return m.execErr
}

var errStoreDown = errors.New("store is down")

func TestOrderServicePlaceOrder(t *testing.T) {
	tests := []struct {
		name    string
		execErr error
		wantErr error
	}{
		{name: "successful insert"},
		{name: "store error is propagated", execErr: errStoreDown, wantErr: errStoreDown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockOrderStore{execErr: tt.execErr}
			svc := NewOrderService(store)

			err := svc.PlaceOrder("order-42", 99.95)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("PlaceOrder() error = %v, want %v", err, tt.wantErr)
			}

			if store.execCalls != 1 {
				t.Fatalf("Exec called %d times, want 1", store.execCalls)
			}
			if store.lastQuery == "" {
				t.Error("Exec got an empty query")
			}
			if len(store.lastArgs) != 2 {
				t.Fatalf("Exec got %d args, want 2", len(store.lastArgs))
			}
			if got, want := store.lastArgs[0], any("order-42"); got != want {
				t.Errorf("Exec args[0] = %v, want %v", got, want)
			}
			if got, want := store.lastArgs[1], any(99.95); got != want {
				t.Errorf("Exec args[1] = %v, want %v", got, want)
			}
		})
	}
}
