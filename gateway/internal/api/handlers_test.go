package api

import (
	"encoding/json"
	"ledger"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBudgetsAPI(t *testing.T) {
	ledger.Reset()
	t.Cleanup(ledger.Reset)

	t.Run("create", func(t *testing.T) {
		rec := doRequest(http.MethodPost, "/api/budgets", `{"category":"еда","limit":5000}`)

		if rec.Code != http.StatusCreated {
			t.Fatalf("POST /api/budgets error = %v", rec.Code)
		}

		checkContentType(t, rec)

		var response BudgetResponse
		err := json.NewDecoder(rec.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Decode error = %v", err)
		}

		if response.Category != "еда" {
			t.Errorf("category = %s; want еда", response.Category)
		}
		if response.Limit != 5000 {
			t.Errorf("limit = %v; want %v", response.Limit, 5000)
		}
	})

	t.Run("list", func(t *testing.T) {
		rec := doRequest(http.MethodGet, "/api/budgets", "")

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/budgets error = %v", rec.Code)
		}

		checkContentType(t, rec)

		var response []BudgetResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("Decode error = %v", err)
		}

		isFood := false
		for _, budget := range response {
			if budget.Category == "еда" {
				isFood = true
			}
		}
		if !isFood {
			t.Errorf("category does not match")
		}
	})

	t.Run("bad_json", func(t *testing.T) {
		rec := doRequest(http.MethodPost, "/api/budgets", `{oops`)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %v, want %v", rec.Code, http.StatusBadRequest)
		}

		checkContentType(t, rec)

		var response map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if response["error"] == "" {
			t.Errorf("want non-empty \"error\" field, got %v", response)
		}
	})

	t.Run("invalid_limit", func(t *testing.T) {
		rec := doRequest(http.MethodPost, "/api/budgets", `{"category":"еда","limit":0}`)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("POST /api/budgets error = %v", rec.Code)
		}

		checkContentType(t, rec)

		var response map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if response["error"] == "" {
			t.Errorf("want non-empty \"error\" field, got %v", response)
		}
	})
}

func TestTransactionsAPI(t *testing.T) {
	ledger.Reset()
	t.Cleanup(ledger.Reset)

	rec := doRequest(http.MethodPost, "/api/budgets", `{"category":"еда","limit":5000}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /api/budgets error = %v", rec.Code)
	}

	t.Run("ok", func(t *testing.T) {
		rec := doRequest(http.MethodPost, "/api/transactions", `{"amount":3000,"category":"еда","date":"2026-10-29"}`)
		if rec.Code != http.StatusCreated {
			t.Fatalf("POST /api/transactions error = %v", rec.Code)
		}

		checkContentType(t, rec)

		var response TransactionResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("decode error: %v", err)
		}

		if response.Amount != 3000 {
			t.Errorf("amount = %v; want %v", response.Amount, 3000)
		}
		if response.Category != "еда" {
			t.Errorf("category = %v; want %v", response.Category, "еда")
		}
		if response.Date != "2026-10-29" {
			t.Errorf("date = %v; want %v", response.Date, "2026-10-29")
		}
	})

	t.Run("list", func(t *testing.T) {
		rec := doRequest(http.MethodGet, "/api/transactions", "")
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/transactions error = %v", rec.Code)
		}

		checkContentType(t, rec)

		var response []TransactionResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("decode error: %v", err)
		}

		if len(response) != 1 {
			t.Fatalf("len(response) = %v; want %v", len(response), 1)
		}
		if response[0].Amount != 3000 {
			t.Errorf("amount = %v; want %v", response[0].Amount, 3000)
		}
		if response[0].Category != "еда" {
			t.Errorf("category = %v; want %v", response[0].Category, "еда")
		}
		if response[0].Date != "2026-10-29" {
			t.Errorf("date = %v; want %v", response[0].Date, "2026-10-29")
		}
	})

	t.Run("exceeded", func(t *testing.T) {
		rec := doRequest(http.MethodPost, "/api/transactions", `{"amount":2500,"category":"еда","date":"2026-10-29"}`)
		if rec.Code != http.StatusConflict {
			t.Fatalf("POST /api/transactions error = %v", rec.Code)
		}

		checkContentType(t, rec)

		var response map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if response["error"] != "budget exceeded" {
			t.Errorf("want \"budget exceeded\" field, got %v", response)
		}
	})

	t.Run("bad_json", func(t *testing.T) {
		rec := doRequest(http.MethodPost, "/api/transactions", `{oops`)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %v, want %v", rec.Code, http.StatusBadRequest)
		}

		checkContentType(t, rec)

		var response map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		if response["error"] == "" {
			t.Errorf("want non-empty \"error\" field, got %v", response)
		}
	})
}

func doRequest(method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	NewRouter().ServeHTTP(rec, req)
	return rec
}

func checkContentType(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %s; want application/json; charset=utf-8", rec.Header().Get("Content-Type"))
	}
}
