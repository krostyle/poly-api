package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	polyhhttp "poly.app/api/internal/adapters/http"
	"poly.app/api/migrations"
)

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	// Track whether schema_migrations is brand-new so we can seed it with
	// migrations that were applied manually before this runner existed.
	var migrationsTableExisted bool
	pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM pg_tables
			WHERE schemaname = 'public' AND tablename = 'schema_migrations'
		)
	`).Scan(&migrationsTableExisted)

	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	// First-run bootstrap: if the table is new but the schema already exists,
	// mark 001 as applied so we don't try to re-run it.
	if !migrationsTableExisted {
		var schemaExists bool
		pool.QueryRow(ctx, `
			SELECT EXISTS(
				SELECT 1 FROM pg_tables
				WHERE schemaname = 'public' AND tablename = 'estudios'
			)
		`).Scan(&schemaExists)
		if schemaExists {
			if _, err := pool.Exec(ctx,
				`INSERT INTO schema_migrations (version) VALUES ('001_initial_schema') ON CONFLICT DO NOTHING`,
			); err != nil {
				return fmt.Errorf("seed schema_migrations: %w", err)
			}
			log.Println("schema_migrations seeded: 001_initial_schema already present")
		}
	}

	rows, err := pool.Query(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("query applied migrations: %w", err)
	}
	applied := map[string]bool{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return err
		}
		applied[v] = true
	}
	rows.Close()

	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	var upFiles []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	for _, name := range upFiles {
		version := strings.TrimSuffix(name, ".up.sql")
		if applied[version] {
			continue
		}
		sql, err := migrations.FS.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("apply %s: %w", name, err)
		}
		if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback(ctx)
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
		log.Printf("migration applied: %s", version)
	}
	return nil
}

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

	if err := runMigrations(context.Background(), pool); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

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
