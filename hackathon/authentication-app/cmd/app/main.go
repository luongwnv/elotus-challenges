package main

import (
	"authentication-app/config"
	_ "authentication-app/docs"
	"authentication-app/internal/controllers"
	database "authentication-app/pkg/database/postgresql"
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/joho/godotenv"
	golog "github.com/luongwnv/go-log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"gorm.io/gorm"
)

type Server struct {
	fiber  *fiber.App
	cfg    *config.Config
	rdbIns *gorm.DB
	logger golog.Logger
}

type Option func(*Server)

func Logger(logger golog.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func NewServer(cfg *config.Config, rdb *gorm.DB, opts ...Option) *Server {
	s := &Server{
		fiber: fiber.New(fiber.Config{
			JSONEncoder:       json.Marshal,
			JSONDecoder:       json.Unmarshal,
			StreamRequestBody: true,
		}),
		cfg:    cfg,
		rdbIns: rdb,
	}

	s.fiber.Use(etag.New(etag.Config{
		Weak: true,
	}))

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Run() error {
	if err := s.MapHandlers(); err != nil {
		return err
	}

	go func() {
		s.logger.Infof("server is listening on port: %s!", s.cfg.ServerPort)
		if err := s.fiber.Listen(":" + s.cfg.ServerPort); err != nil {
			s.logger.Panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	s.logger.Info("server stopped successfully")
	return s.fiber.Shutdown()
}

// @title SIMPLE AUTHENTICATION APP API
// @version 1.0
// @description API documentation for SIMPLE AUTHENTICATION APP services
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
// @security BearerAuth
func (s *Server) MapHandlers() error {
	app := s.fiber
	// validate := validator.New()

	app.Use(func(c *fiber.Ctx) error {
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(func(c *fiber.Ctx) error {
		defer HandlePanic("HTTP Service")
		return c.Next()
	})

	app.Get("/api/swagger/*", fiberSwagger.WrapHandler)

	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05.999",
		Format:     "[${ip}] ${time} ${locals:requestid} ${method} ${path} ${status} ${latency}\n",
	}))

	// Health check routes
	monitoringHandler := controllers.NewHandler(s.cfg, s.logger, s.rdbIns)
	app.Get("/api/readiness", timeout.New(monitoringHandler.Readiness, time.Duration(s.cfg.ServerCtxDefaultTimeout)*time.Second))
	app.Get("/api/liveness", monitoringHandler.Liveness)

	// Auth routes
	authController := controllers.NewAuthController(s.cfg, s.logger, s.rdbIns)
	authGroup := app.Group("/auth")
	authGroup.Post("/register", authController.Register)

	golog.Info("Loaded all route!")

	return nil
}

func HandlePanic(contextName string) {
	if r := recover(); r != nil {
		fmt.Println("An error occurred:", r)
	}
}

func main() {
	defer HandlePanic("MAIN")

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

	s := NewServer(cfg, db, Logger(appLogger))

	go func() {
		defer HandlePanic("HTTP Service")
		if err = s.Run(); err != nil {
			_, file, line, _ := runtime.Caller(1)
			appLogger.Errorf("Server failed: %v at %s:%d", err, path.Base(file), line)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
