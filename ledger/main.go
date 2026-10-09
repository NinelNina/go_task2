package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Ledger service started")

	file, err := os.Open("budgets.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	if err := LoadBudgets(reader); err != nil {
		fmt.Println("Error loading budgets:", err)
		return
	}
	fmt.Println("Budgets loaded")

	tx1 := Transaction{
		Amount:      4600.50,
		Category:    "еда",
		Description: "Покупка в супермаркете",
		Date:        time.Now(),
	}

	tx2 := Transaction{
		Amount:      450.00,
		Category:    "еда",
		Description: "Кофе и круассан",
		Date:        time.Now(),
	}

	tx3 := Transaction{
		Amount:      3000.00,
		Category:    "Транспорт",
		Description: "Проездной на месяц",
		Date:        time.Now(),
	}

	err = AddTransaction(tx1)
	if err != nil {
		fmt.Println("Error adding transaction:", err)
	}

	err = AddTransaction(tx2)
	if err != nil {
		fmt.Println("Error adding transaction:", err)
	}

	err = AddTransaction(tx3)
	if err != nil {
		fmt.Println("Error adding transaction:", err)
	}

	allTxs := ListTransactions()
	for _, tx := range allTxs {
		fmt.Println(tx)
	}

	fmt.Println("------------------------\nPolymorphism test")
	testTx := Transaction{Amount: -10, Category: "еда"}
	testBudget := Budget{Category: "", Limit: 1000}

	if err := CheckValid(testTx); err != nil {
		fmt.Println("CheckValid(testTx) found error:", err)
	}

	if err := CheckValid(testBudget); err != nil {
		fmt.Println("CheckValid(testBudget) found error:", err)
	}
}
