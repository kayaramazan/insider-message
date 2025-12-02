package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/kayaramazan/insider-message/config"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db  *sql.DB
	cfg *config.DatabaseConfig
}

func NewPostgresDB(cfg *config.DatabaseConfig) *PostgresDB {
	return &PostgresDB{cfg: cfg}
}

func (p *PostgresDB) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.cfg.Host,
		p.cfg.Port,
		p.cfg.User,
		p.cfg.Password,
		p.cfg.DBName,
		p.cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	p.db = db
	log.Printf("Connection successful")
	return nil
}

func (p *PostgresDB) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresDB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.QueryContext(ctx, query, args...)
}

func (p *PostgresDB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.db.QueryRowContext(ctx, query, args...)
}

func (p *PostgresDB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.db.ExecContext(ctx, query, args...)
}
