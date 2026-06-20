package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/lazzerex/gitrpg/internal/config"
	"github.com/lazzerex/gitrpg/internal/github"
	"github.com/lazzerex/gitrpg/internal/server"
	"github.com/lazzerex/gitrpg/internal/worker"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", "error", err)
		os.Exit(1)
	}

	if cfg.Server.Env == "development" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	db, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		logger.Error("db connect failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		logger.Error("db ping failed", "error", err)
		os.Exit(1)
	}
	logger.Info("db connected")

	opts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		logger.Error("redis url parse failed", "error", err)
		os.Exit(1)
	}
	rdb := redis.NewClient(opts)
	defer func() { _ = rdb.Close() }()

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Error("redis ping failed", "error", err)
		os.Exit(1)
	}
	logger.Info("redis connected")

	githubSvc := github.NewService(db, logger)
	w := worker.New(githubSvc, logger)

	srv := server.New(cfg, db, rdb, logger, w)

	if err := srv.LoadTemplates("web/templates"); err != nil {
		logger.Error("template load failed", "error", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	go w.Start(workerCtx)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	logger.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
}
