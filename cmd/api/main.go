package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	polyhhttp "poly.app/api/internal/adapters/http"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := polyhhttp.NewRouter()

	addr := fmt.Sprintf(":%s", port)
	log.Printf("poly-api listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
