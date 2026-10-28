package api

import (
	"encoding/json"
	"errors"
	"ledger"
	"log"
	"net/http"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	tx, err := ToTransaction(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = tx.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = ledger.AddTransaction(tx)
	if err != nil {
		if errors.Is(err, ledger.ErrBudgetExceeded) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	txs := ledger.ListTransactions()
	created := txs[len(txs)-1]
	writeJSON(w, http.StatusCreated, FromTransaction(created))
}

func ListTransactions(w http.ResponseWriter, r *http.Request)  {
	txs := ledger.ListTransactions()
	resp := make([]TransactionResponse, 0, len(txs))

	for _, tx := range txs {
		resp = append(resp, FromTransaction(tx))
	}

	writeJSON(w, http.StatusOK, resp)
}

func CreateBudget(w http.ResponseWriter, r *http.Request) {
	var req CreateBudgetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	b := ToBudget(req)
	err = b.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = ledger.SetBudget(b)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, FromBudget(b))
}

func ListBudgets(w http.ResponseWriter, r *http.Request)  {
	budgets := ledger.ListBudgets()
	resp := make([]BudgetResponse, 0, len(budgets))

	for _, budget := range budgets {
		resp = append(resp, FromBudget(budget))
	}
	
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("encode error:", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
