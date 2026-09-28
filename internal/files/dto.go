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
