package files

import (
	"fmt"
	"time"
)

type FileExpirationTicker struct {
	service  *Service
	interval time.Duration
}

func NewFileExpirationTicker(service *Service, interval time.Duration) *FileExpirationTicker {
	return &FileExpirationTicker{
		service:  service,
		interval: interval,
	}
}
func (fe *FileExpirationTicker) Start() {
	ticker := time.NewTicker(fe.interval)
	defer ticker.Stop()
	for range ticker.C {
		fe.Run()
	}
}

func (fe *FileExpirationTicker) Run() {
	filesDeleted, err := fe.service.DeleteExpired()
	fmt.Printf("Deleted %d expired files.", filesDeleted)
	if err != nil {
		return
	}
}
