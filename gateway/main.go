package main

import (
	"fmt"
	"gateway/internal/api"
	"log"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pong")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", pingHandler)
	mux.HandleFunc("POST /api/transactions", api.CreateTransaction)
	mux.HandleFunc("GET /api/transactions", api.ListTransactions)
	mux.HandleFunc("POST /api/budgets", api.CreateBudget)
	mux.HandleFunc("GET /api/budgets", api.ListBudgets)

	log.Println("Gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", api.Logging(mux)))
}
