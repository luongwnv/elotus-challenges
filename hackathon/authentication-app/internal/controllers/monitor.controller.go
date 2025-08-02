package controllers

import (
	"authentication-app/config"
	"fmt"

	golog "github.com/luongwnv/go-log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type serverHandlers struct {
	cfg    *config.Config
	logger golog.Logger
	rdbIns *gorm.DB
}

type Readiness struct {
	Redis        string `json:"redis,omitempty" validate:"omitempty,oneof=true"`
	RedisCluster string `json:"redis_cluster,omitempty" validate:"omitempty,oneof=true"`
	RDB          string `json:"relational_database,omitempty" validate:"omitempty,oneof=true"`
}

type Liveness struct {
	Message string `json:"message"`
}

type SentryCheck struct {
	Message string `json:"message"`
}

// Readiness handler
func (h *serverHandlers) Readiness(c *fiber.Ctx) (err error) {
	var (
		resp       Readiness
		httpStatus int
	)

	return c.Status(httpStatus).JSON(resp)
}

// Liveness handler
func (h *serverHandlers) PublicLiveness(c *fiber.Ctx) (err error) {
	fmt.Println("PublicLiveness")
	var (
		resp Liveness
	)
	resp.Message = "Alive"

	return c.Status(fiber.StatusOK).JSON(resp)
}

// Liveness handler
func (h *serverHandlers) Liveness(c *fiber.Ctx) (err error) {
	var (
		resp Liveness
	)
	resp.Message = "Alive"

	return c.Status(fiber.StatusOK).JSON(resp)
}

// HTTP handlers interface
type Handlers interface {
	AddMonitoringRoutes(router fiber.Router)
	Readiness(c *fiber.Ctx) error
	Liveness(c *fiber.Ctx) error
	PublicLiveness(c *fiber.Ctx) error
	SentryCheck(c *fiber.Ctx) error
}

// Create new handler instance
func NewHandler(cfg *config.Config, logger golog.Logger, rdb *gorm.DB) (h *serverHandlers) {
	return &serverHandlers{
		cfg:    cfg,
		logger: logger,
		rdbIns: rdb,
	}
}
