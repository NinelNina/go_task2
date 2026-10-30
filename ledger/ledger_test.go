package ledger

import (
	"errors"
	"testing"
	"time"
)

func TestTransactionValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tx      Transaction
		wantErr bool
	}{
		{name: "valid", tx: Transaction{Amount: 100, Category: "еда"}, wantErr: false},
		{name: "zero amount", tx: Transaction{Amount: 0, Category: "еда"}, wantErr: true},
		{name: "negative amount", tx: Transaction{Amount: -100, Category: "еда"}, wantErr: true},
		{name: "empty category", tx: Transaction{Amount: 100, Category: ""}, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.tx.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestBudgetValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		b       Budget
		wantErr bool
	}{
		{name: "valid", b: Budget{Limit: 100, Category: "еда"}, wantErr: false},
		{name: "zero limit", b: Budget{Limit: 0, Category: "еда"}, wantErr: true},
		{name: "negative limit", b: Budget{Limit: -100, Category: "еда"}, wantErr: true},
		{name: "empty category", b: Budget{Limit: 100, Category: ""}, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.b.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestBudgetLimit(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	wantErr := ErrBudgetExceeded
	if err := SetBudget(Budget{Category: "еда", Limit: 5000}); err != nil {
		t.Fatalf("SetBudget() error = %v", err)
	}

	if err := AddTransaction(Transaction{
		Amount:      3000,
		Category:    "еда",
		Description: "транзакция",
		Date:        time.Now(),
	}); err != nil {
		t.Fatalf("first AddTransaction() error = %v", err)
	}

	err := AddTransaction(Transaction{
		Amount:      2500,
		Category:    "еда",
		Description: "транзакция",
		Date:        time.Now(),
	})

	if !errors.Is(err, wantErr) {
		t.Errorf("AddTransaction() error = %v, wantErr %v", err, wantErr)
	}

	txLen := len(ListTransactions())
	if txLen != 1 {
		t.Errorf("AddTransaction() len(ListTransactions()) = %v, want %v", txLen, 1)
	}
}
