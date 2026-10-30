package api

import (
	"fmt"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", pingHandler)
	mux.HandleFunc("POST /api/transactions", CreateTransaction)
	mux.HandleFunc("GET /api/transactions", ListTransactions)
	mux.HandleFunc("POST /api/budgets", CreateBudget)
	mux.HandleFunc("GET /api/budgets", ListBudgets)

	return Logging(mux)
}
