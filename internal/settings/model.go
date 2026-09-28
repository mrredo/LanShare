package settings

import "time"

type Settings struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	StoragePath string `gorm:"unique;not null" json:"storage_path"`
	// bytes
	StorageLimit int64 `gorm:"not null" json:"storage_limit"`
	// bytes
	MaxFileSize    int64 `gorm:"not null" json:"max_file_size"`
	UploadsEnabled bool  `gorm:"not null" json:"uploads_enabled"`
	// seconds
	DefaultExpiry     *int64    `json:"default_expiry"`
	AdminPasswordHash string    `gorm:"not null" json:"-"`
	CreatedAt         time.Time `gorm:"not null" json:"created_at"`
}

func (Settings) TableName() string {
	return "settings"
}

type Setting = Settings

func DefaultSettings(adminPasswordHash string) *Settings {
	defaultExpiry := int64(24 * 60 * 60) // 24 hours
	return &Settings{
		StoragePath:       "./uploads",
		StorageLimit:      10 * 1024 * 1024 * 1024, // 10 GB
		MaxFileSize:       1024 * 1024 * 1024,      // 1 GB
		UploadsEnabled:    true,
		DefaultExpiry:     &defaultExpiry,
		AdminPasswordHash: adminPasswordHash,
	}
}
