package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5"
)

func RunMigrations(connString string) error {
	// Simple migration runner
	ctx := context.Background()
	config, err := pgx.ParseConfig(connString)
	if err != nil {
		return err
	}
	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		var version int
		fmt.Sscanf(file.Name(), "%d_", &version)
		if version == 0 {
			continue
		}

		var exists bool
		err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)`, version).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		sql, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file.Name(), err)
		}

		if _, err := conn.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("execute migration %s: %w", file.Name(), err)
		}
		if _, err := conn.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
			return err
		}
		fmt.Printf("Applied migration: %s\n", file.Name())
	}
	return nil
}
