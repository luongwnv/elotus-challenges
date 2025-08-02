package controllers

import (
	"authentication-app/internal/models"
	"authentication-app/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	golog "github.com/luongwnv/go-log"
	"gorm.io/gorm"
)

type FileController struct {
	logger golog.Logger
	db     *gorm.DB
}

func NewFileController(logger golog.Logger, db *gorm.DB) *FileController {
	return &FileController{
		logger: logger,
		db:     db,
	}
}

// @Summary Upload file
// @Description Upload an image file (max 8MB)
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file to upload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 413 {object} map[string]string
// @Security BearerAuth
// @Router /files/upload [post]
func (ac *FileController) UploadFile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file provided or invalid form data",
		})
	}

	// Check file size (8MB = 8 * 1024 * 1024 bytes)
	maxSize := int64(8 * 1024 * 1024)
	if file.Size > maxSize {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
			"error": "File size exceeds 8MB limit",
		})
	}

	// Check if file is an image
	if !utils.IsImageContentType(file.Header.Get("Content-Type")) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File must be an image (JPEG, PNG, GIF, WebP)",
		})
	}

	// Open file to verify content type
	src, err := file.Open()
	if err != nil {
		ac.logger.Errorf("Failed to open uploaded file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process file",
		})
	}
	defer src.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		ac.logger.Errorf("Failed to read file header: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process file",
		})
	}

	// Reset file pointer
	src.Close()
	src, err = file.Open()
	if err != nil {
		ac.logger.Errorf("Failed to reopen uploaded file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process file",
		})
	}
	defer src.Close()

	// Generate unique filename
	filename := utils.GenerateUniqueFilename(file.Filename)
	filePath := filepath.Join("/tmp", filename)

	// Save file to /tmp
	dst, err := os.Create(filePath)
	if err != nil {
		ac.logger.Errorf("Failed to create temp file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		ac.logger.Errorf("Failed to copy file data: %v", err)
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Save file metadata to database
	fileUpload := models.FileUpload{
		ID:           uuid.New(),
		UserID:       userID,
		Filename:     filename,
		OriginalName: file.Filename,
		ContentType:  file.Header.Get("Content-Type"),
		Size:         file.Size,
		FilePath:     filePath,
		UserAgent:    c.Get("User-Agent"),
		IPAddress:    c.IP(),
		CreatedAt:    time.Now(),
	}

	if err := ac.db.Create(&fileUpload).Error; err != nil {
		ac.logger.Errorf("Failed to save file metadata: %v", err)
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file metadata",
		})
	}

	return c.JSON(fiber.Map{
		"message":       "File uploaded successfully",
		"file_id":       fileUpload.ID,
		"filename":      filename,
		"original_name": file.Filename,
		"content_type":  file.Header.Get("Content-Type"),
		"size":          file.Size,
		"uploaded_at":   fileUpload.CreatedAt,
	})
}
