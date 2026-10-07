package files

import (
	"mime/multipart"
	"time"
)

type UploadDTO struct {
	File      *multipart.FileHeader `form:"file"`
	Text      string                `form:"text"`
	URL       string                `form:"url"`
	ExpiresAt *time.Time            `form:"expires_at"`
}
type FileResponse struct {
	ID          string     `gorm:"primaryKey;not null" json:"id"`
	Filename    *string    `json:"filename,omitempty"`
	StoragePath *string    `gorm:"unique" json:"storage_path,omitempty"`
	Size        *int64     `json:"size,omitempty"`
	Text        *string    `json:"text,omitempty"`
	URL         *string    `gorm:"column:url" json:"url,omitempty"`
	UploadedAt  time.Time  `gorm:"not null" json:"uploaded_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

func FilesToFileResponsefunc(files []File) []FileResponse {
	response := make([]FileResponse, 0, len(files))

	for _, file := range files {
		response = append(response, FileResponse{
			ID:          file.ID,
			Filename:    file.Filename,
			StoragePath: file.StoragePath,
			Size:        file.Size,
			Text:        file.Text,
			URL:         file.URL,
			UploadedAt:  file.UploadedAt,
			ExpiresAt:   file.ExpiresAt,
		})
	}

	return response
}
