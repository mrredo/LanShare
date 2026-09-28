package files

import (
	"time"
)

type File struct {
	ID          string     `gorm:"primaryKey;not null" json:"id"`
	Filename    *string    `json:"filename,omitempty"`
	StoragePath *string    `gorm:"unique" json:"storage_path,omitempty"`
	Size        *int64     `json:"size,omitempty"`
	OwnerCookie string     `gorm:"not null;index" json:"owner_cookie"`
	Text        *string    `json:"text,omitempty"`
	URL         *string    `gorm:"column:url" json:"url,omitempty"`
	UploadedAt  time.Time  `gorm:"not null" json:"uploaded_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

func (File) TableName() string {
	return "files"
}
