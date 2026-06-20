package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	"github.com/lazzerex/gitrpg/internal/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", "error", err)
		os.Exit(1)
	}

	db, err := sql.Open("pgx", cfg.Database.URL)
	if err != nil {
		logger.Error("db open failed", "error", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Error("set dialect failed", "error", err)
		os.Exit(1)
	}

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	migrationsDir := "migrations"
	if err := goose.RunContext(context.Background(), cmd, db, migrationsDir); err != nil {
		logger.Error("migration failed", "cmd", cmd, "error", err)
		os.Exit(1)
	}

	logger.Info("migration done", "cmd", cmd)
}
