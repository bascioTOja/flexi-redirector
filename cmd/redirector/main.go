package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"flexi-redirector/internal/config"
	"flexi-redirector/internal/db"
	httpapp "flexi-redirector/internal/http"
	"flexi-redirector/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	logger := log.New(os.Stdout, "[REDIRECTOR] ", log.LstdFlags)

	cfg, err := config.LoadFromEnv()
	if err != nil {
		logger.Fatalf("invalid config: %v", err)
	}

	gormDB, closeFn, err := db.Open(cfg.DB)
	if err != nil {
		logger.Fatalf("failed to open database: %v", err)
	}
	defer func() {
		if err := closeFn(); err != nil {
			logger.Printf("failed to close db: %v", err)
		}
	}()

	repos := repository.NewGormRepositories(gormDB)

	router := httpapp.NewRouter(httpapp.Deps{
		ShortURLs:  repos.ShortURLs,
		CountViews: cfg.CountViews,
		Logger:     logger,
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Printf("Starting server on :%s...", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("server forced to shutdown: %v", err)
	}

	logger.Println("Server exited properly")
}
