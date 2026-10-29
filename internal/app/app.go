package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"

	"golang-starter-kit/internal/config"
	"golang-starter-kit/pkg/cache"
	"golang-starter-kit/pkg/database"
)

type Container struct {
	Config  *config.Config
	DB      *database.DBManager
	Cache   cache.Cache
	Domains []DomainRouter
}

func (c *Container) RegisterDomain(d DomainRouter) {
	c.Domains = append(c.Domains, d)
}

type App struct {
	Config    *config.Config
	Server    *fiber.App
	Container *Container
}

func BuildApp(cfg *config.Config, db *database.DBManager, c ...cache.Cache) *App {
	server := NewServer(cfg)

	var activeCache cache.Cache
	if len(c) > 0 && c[0] != nil {
		activeCache = c[0]
	} else {
		activeCache = cache.NewMemoryCache()
	}

	container := &Container{
		Config:  cfg,
		DB:      db,
		Cache:   activeCache,
		Domains: make([]DomainRouter, 0),
	}

	// Daftarkan seluruh modul aplikasi dari katalog domains.go
	RegisterAppDomains(container)

	RegisterRoutes(server, container)

	return &App{
		Config:    cfg,
		Server:    server,
		Container: container,
	}
}

// Start menjalankan server dengan graceful shutdown
func (a *App) Start() error {
	addr := fmt.Sprintf(":%d", a.Config.Port)

	// Channel untuk menangkap signal sistem (SIGINT, SIGTERM)
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("Starting HTTP Server", "port", a.Config.Port, "env", a.Config.AppEnv)
		if err := a.Server.Listen(addr); err != nil {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdownChan:
		slog.Info("Received shutdown signal", "signal", sig.String())
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.Server.ShutdownWithContext(ctx); err != nil {
			slog.Error("Forced shutdown error", "error", err)
		}

		if a.Container.DB != nil {
			a.Container.DB.CloseAll(ctx)
		}
		slog.Info("Graceful shutdown completed successfully")
		return nil
	}
}
