package main

import (
	"context"
	"log/slog"
	"qore-be/internal/config"
	"qore-be/internal/person"
	"qore-be/internal/server"

	"log"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))

	cfg := loadConfig()
	start(cfg)
}

func start(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())

	db := setupDB(cfg.DBUrl)
	srv := server.New(
		server.WithConfig(cfg),
		server.WithPersonController(newPersonCtrl(db)),
	)

	go srv.Start(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancel()
}

func loadConfig() *config.Config {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return cfg
}

func setupDB(url string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func newPersonCtrl(db *gorm.DB) *person.Controller {
	svc, err := person.NewService(person.WithRepository(person.NewRepository(db)))
	if err != nil {
		log.Fatalf("failed to create person service: %v", err)
	}

	ctrl, err := person.NewController(person.WithService(svc))
	if err != nil {
		log.Fatalf("failed to create store controller: %v", err)
	}

	return ctrl
}
