package controllers

import (
	"authentication-app/config"
	dto "authentication-app/internal/DTOs"
	"authentication-app/internal/models"
	"authentication-app/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	golog "github.com/luongwnv/go-log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	cfg    *config.Config
	logger golog.Logger
	db     *gorm.DB
}

func NewAuthController(cfg *config.Config, logger golog.Logger, db *gorm.DB) *AuthController {
	return &AuthController{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /auth/register [post]
func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username and password are required",
		})
	}

	if len(req.Username) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username must be at least 3 characters",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters",
		})
	}

	// Check if user already exists
	var existingUser models.User
	if err := ac.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username already exists",
		})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ac.logger.Errorf("Failed to hash password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Create user
	user := models.User{
		ID:           uuid.New(),
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := ac.db.Create(&user).Error; err != nil {
		ac.logger.Errorf("Failed to create user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	// Generate JWT token
	token, expiresAt, err := utils.GenerateJWTToken(user.ID, user.Username, ac.cfg.JWTSecret, ac.cfg.JWTExpireMinutes)
	if err != nil {
		ac.logger.Errorf("Failed to generate token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: dto.UserInfo{
			ID:       user.ID,
			Username: user.Username,
		},
	})
}
