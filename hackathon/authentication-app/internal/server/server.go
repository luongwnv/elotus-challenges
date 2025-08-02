package server

import (
	"authentication-app/config"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	golog "github.com/luongwnv/go-log"
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
