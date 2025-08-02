package main

import (
	"authentication-app/config"
	_ "authentication-app/docs"
	server "authentication-app/internal/server"
	database "authentication-app/pkg/database/postgresql"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"

	"github.com/joho/godotenv"
	golog "github.com/luongwnv/go-log"
)

func main() {
	defer server.HandlePanic("MAIN")

	var (
		cfg *config.Config
		err error
	)
	err = godotenv.Load(".env")
	if err != nil {
		golog.Info("No .env file found, using environment variables from Docker Compose")
	}

	cfg, err = config.LoadConfig()
	if err != nil {
		golog.Panicf("Load config: %v", err)
	}

	appLogger := golog.NewLogger(
		golog.WithFormat(cfg.LoggerEncoding),
		golog.WithLevel(cfg.LoggerLevel),
		golog.WithCallerPathType(cfg.LoggerIsFullPathCaller),
	)

	// Initialize database
	db, err := database.New(cfg)
	if err != nil {
		appLogger.Errorf("Failed to connect to database: %v", err)
		golog.Panicf("Database connection failed: %v", err)
	}

	s := server.NewServer(cfg, db, server.Logger(appLogger))

	go func() {
		defer server.HandlePanic("HTTP Service")
		if err = s.Run(); err != nil {
			_, file, line, _ := runtime.Caller(1)
			appLogger.Errorf("Server failed: %v at %s:%d", err, path.Base(file), line)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
