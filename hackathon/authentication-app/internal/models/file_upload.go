package models

import (
	"time"

	"github.com/google/uuid"
)

type FileUpload struct {
	ID           uuid.UUID  `gorm:"column:id;primaryKey" json:"id"`
	UserID       uuid.UUID  `gorm:"column:user_id;not null;index" json:"user_id"`
	User         User       `gorm:"foreignKey:UserID" json:"user"`
	Filename     string     `gorm:"column:filename;not null" json:"filename"`
	OriginalName string     `gorm:"column:original_name;not null" json:"original_name"`
	ContentType  string     `gorm:"column:content_type;not null" json:"content_type"`
	Size         int64      `gorm:"column:size;not null" json:"size"`
	FilePath     string     `gorm:"column:file_path;not null" json:"file_path"`
	UserAgent    string     `gorm:"column:user_agent" json:"user_agent"`
	IPAddress    string     `gorm:"column:ip_address" json:"ip_address"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (FileUpload) TableName() string {
	return "authentication-app.file_uploads"
}
