package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	log.Println("poly-worker starting (plazo recalculation cron)")
	// TODO: inicializar DB, instanciar RecalcularPlazosUseCase, registrar cron job diario.
	select {} // bloquea hasta señal de sistema
}
