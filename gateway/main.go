package main

import (
	"gateway/internal/api"
	"log"
	"net/http"
)

func main() {
	log.Println("Gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", api.NewRouter()))
}
