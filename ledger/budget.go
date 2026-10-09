package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
}

func (b Budget) Validate() error {
	if b.Limit <= 0 {
		return errors.New("Limit must be greater than zero")
	}
	if b.Category == "" {
		return errors.New("Category cannot be empty")
	}
	return nil
}

var budgets = make(map[string]Budget)

func SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return err
	}
	budgets[b.Category] = b
	return nil
}

func LoadBudgets(r io.Reader) error {
	decoder := json.NewDecoder(r)
	var loadedBudgets []Budget
	err := decoder.Decode(&loadedBudgets)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	for _, budget := range loadedBudgets {
		if err := SetBudget(budget); err != nil {
			return err
		}
	}
	return nil
}
