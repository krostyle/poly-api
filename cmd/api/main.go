package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	polyhhttp "poly.app/api/internal/adapters/http"
)

func main() {
	_ = godotenv.Load()

	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("database unreachable: %v", err)
	}
	log.Println("database connected")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := polyhhttp.NewRouter(pool)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("poly-api listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
