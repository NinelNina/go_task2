package api

import (
	"ledger"
	"time"
)

type CreateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type TransactionResponse struct {
	ID          int       `json:"id"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type CreateBudgetRequest struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

type BudgetResponse struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

func ToTransaction(tReq CreateTransactionRequest) (ledger.Transaction, error) {
	parsedDate, err := time.Parse("2006-01-02", tReq.Date)
	if err != nil {
		return ledger.Transaction{}, err
	}

	return ledger.Transaction{
		Category:    tReq.Category,
		Amount:      tReq.Amount,
		Description: tReq.Description,
		Date:        parsedDate,
	}, nil
}

func FromTransaction(t ledger.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:          t.ID,
		Amount:      t.Amount,
		Category: t.Category,
		Description: t.Description,
		Date:        t.Date,
	}
}

func ToBudget(bReq CreateBudgetRequest) ledger.Budget {
	return ledger.Budget{
		Category: bReq.Category,
		Limit:    bReq.Limit,
	}
}

func FromBudget(b ledger.Budget) BudgetResponse {
	return BudgetResponse{
		Category: b.Category,
		Limit: b.Limit,
	}
}



