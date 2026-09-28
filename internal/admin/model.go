package admin

import (
	"time"
)

type AdminSession struct {
	ID        string    `gorm:"primaryKey;not null" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
}

func (AdminSession) TableName() string {
	return "admin_sessions"
}

type Session = AdminSession
