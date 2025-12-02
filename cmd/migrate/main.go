package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kayaramazan/insider-message/config"
	_ "github.com/lib/pq"
)

func main() {
	direction := flag.String("direction", "up", "Migration direction: up or down")
	flag.Parse()

	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create migrations tracking table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Get migration files
	migrationsDir := "migrations"
	suffix := "." + *direction + ".sql"

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var migrationFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), suffix) {
			migrationFiles = append(migrationFiles, f.Name())
		}
	}

	// Sort files
	sort.Strings(migrationFiles)
	if *direction == "down" {
		// Reverse for down migrations
		for i, j := 0, len(migrationFiles)-1; i < j; i, j = i+1, j-1 {
			migrationFiles[i], migrationFiles[j] = migrationFiles[j], migrationFiles[i]
		}
	}

	for _, file := range migrationFiles {
		version := strings.Split(file, "_")[0]

		// Check if already applied (for up) or not applied (for down)
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to check migration status: %v", err)
		}

		if *direction == "up" && exists {
			log.Printf("Skipping %s (already applied)", file)
			continue
		}
		if *direction == "down" && !exists {
			log.Printf("Skipping %s (not applied)", file)
			continue
		}

		// Read and execute migration
		content, err := os.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			log.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		log.Printf("Applying %s...", file)

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Failed to begin transaction: %v", err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			log.Fatalf("Failed to execute migration %s: %v", file, err)
		}

		// Update migrations table
		if *direction == "up" {
			_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
		} else {
			_, err = tx.Exec("DELETE FROM schema_migrations WHERE version = $1", version)
		}
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to update migrations table: %v", err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalf("Failed to commit transaction: %v", err)
		}

		log.Printf("Applied %s successfully", file)
	}

	log.Println("Migration completed!")
}
