package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// IsImageContentType checks if the given content type is an image type
func IsImageContentType(contentType string) bool {
	imageTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	contentType = strings.ToLower(contentType)
	for _, imageType := range imageTypes {
		if contentType == imageType {
			return true
		}
	}
	return false
}

// GenerateUniqueFilename generates a unique filename based on the original name and current timestamp
func GenerateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d_%s%s", timestamp, GenerateRandomString(8), ext)
}
