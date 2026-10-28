package ledger

import (
	"errors"
	"time"
)

var ErrBudgetExceeded = errors.New("budget exceeded")

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        time.Time
}

func (tx Transaction) Validate() error {
	if tx.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if tx.Category == "" {
		return errors.New("category cannot be empty")
	}
	return nil
}

var transactions []Transaction

func AddTransaction(tx Transaction) error {
	if err := tx.Validate(); err != nil {
		return err
	}
	budget, exists := budgets[tx.Category]
	if exists {
		expensesAmount := tx.Amount
		for _, transaction := range transactions {
			if transaction.Category == tx.Category {
				expensesAmount += transaction.Amount
			}
		}
		if expensesAmount > budget.Limit {
			return ErrBudgetExceeded
		}
	}

	tx.ID = len(transactions) + 1
	transactions = append(transactions, tx)
	return nil
}

func ListTransactions() []Transaction {
	return transactions
}
