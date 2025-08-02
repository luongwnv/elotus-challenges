package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `gorm:"column:id;primaryKey;type:uuid" json:"id"`
	Username     string     `gorm:"column:username;uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"column:password_hash;not null" json:"password_hash"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (User) TableName() string {
	return "authentication-app.users"
}
