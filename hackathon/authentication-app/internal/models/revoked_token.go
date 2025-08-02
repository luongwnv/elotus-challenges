package models

import (
	"time"

	"github.com/google/uuid"
)

type RevokedToken struct {
	ID        uuid.UUID  `gorm:"column:id;primaryKey" json:"id"`
	TokenID   string     `gorm:"column:token_id;not null;index" json:"token_id"`
	UserID    uuid.UUID  `gorm:"column:user_id;not null;index" json:"user_id"`
	RevokedAt time.Time  `gorm:"column:revoked_at;not null" json:"revoked_at"`
	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;" json:"updated_at"`
}

func (RevokedToken) TableName() string {
	return "authentication-app.revoked_tokens"
}
