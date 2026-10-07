package files

import (
	"errors"
	"fmt"
	"io"
	"lanshare/internal/auth"
	"lanshare/internal/settings"
	"lanshare/pkg/str"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	repo            *Repo
	settingsService *settings.Service
	authService     *auth.Service
}

func NewService(repo *Repo, settingsService *settings.Service, authService *auth.Service) *Service {
	return &Service{
		repo:            repo,
		settingsService: settingsService,
		authService:     authService,
	}
}
func (s *Service) CreateFile(uploadDto *UploadDTO, ownerCookie string) (*File, error) {
	var storagePath *string
	var fileSize *int64
	var originalFilename *string
	id, _ := str.GenerateRandomString(8)

	if uploadDto.File != nil {
		if err := s.CanUploadFile(uploadDto.File.Size); err != nil {
			return nil, err
		}

		ext := filepath.Ext(uploadDto.File.Filename)
		path := filepath.Join(s.settingsService.Get().StoragePath, id+ext)

		// izveido mapi, kur glabāt failus
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, fmt.Errorf("neizdevās izveidot mapi: %w", err)
		}

		// Atveram ienākošo failu
		src, err := uploadDto.File.Open()
		if err != nil {
			return nil, fmt.Errorf("neizdevās atvērt augšupielādēto failu: %w", err)
		}
		defer src.Close()

		// Izveidojam failu mapē
		dst, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("neizdevās izveidot failu diskā: %w", err)
		}
		defer dst.Close()

		// Kopējam saturu uz disku
		if _, err := io.Copy(dst, src); err != nil {
			os.Remove(path)
			return nil, fmt.Errorf("neizdevās saglabāt failu: %w", err)
		}

		storagePath = &path
		fileSize = &uploadDto.File.Size
		originalFilename = &uploadDto.File.Filename
	}

	defaultExpiry := int64(86400) // 24 stundas
	if s.settingsService.Get() != nil && s.settingsService.Get().DefaultExpiry != nil {
		defaultExpiry = *s.settingsService.Get().DefaultExpiry
	}
	expiresAt := time.Now().Add(time.Duration(defaultExpiry) * time.Second)
	if uploadDto.ExpiresAt != nil {
		expiresAt = *uploadDto.ExpiresAt
	}

	fileRecord := &File{
		ID:          id,
		Filename:    originalFilename,
		StoragePath: storagePath,
		Size:        fileSize,
		OwnerCookie: ownerCookie,
		Text:        &uploadDto.Text,
		URL:         &uploadDto.URL,
		UploadedAt:  time.Now(),
		ExpiresAt:   &expiresAt,
	}

	if err := s.repo.Create(fileRecord); err != nil {
		if storagePath != nil {
			os.Remove(*storagePath)
		}
		return nil, err
	}

	return fileRecord, nil
}

func (s *Service) CanUploadFile(fileSize int64) error {
	if !s.settingsService.Get().UploadsEnabled {
		return errors.New("Augšupielādēšana ir izslēgta.")
	}
	// FUnckcionālā prasība: Vai ir pietiekama vieta failiem
	if fileSize > s.getAvailableSpace() {
		return errors.New("Fails ir pārāk liels, lai tas tiktu saglabāts mapē.")
	}
	if fileSize > s.settingsService.Get().MaxFileSize {
		return errors.New("Fails ir pārāk liels.")
	}
	return nil
}
func (s *Service) getAvailableSpace() int64 {
	if s.settingsService != nil && s.settingsService.Get() != nil {
		return s.settingsService.Get().StorageLimit
	}
	return 0
}
func (s *Service) GetById(id string) (File, error) {
	file, err := s.repo.GetById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return File{}, errors.New("Fails netika atrasts")
	}
	return file, err
}

// CountFilesByOwnerCookie returns the total count of files associated with the given owner cookie.
func (s *Service) CountFilesByOwnerCookie(ownerCookie string) (int64, error) {
	return s.repo.CountFilesByOwnerCookie(ownerCookie)
}

func (s *Service) CountAllFiles() (int64, error) {
	return s.repo.CountAllFiles()
}

// DeleteAll removes all records from the underlying repository and returns an error if the operation fails.
func (s *Service) DeleteAll() error {
	return s.repo.DeleteAll()
}
func (s *Service) DeleteById(id string) error {
	return s.repo.DeleteById(id)
}
func (s *Service) DeleteByOwnerCookie(ownerCookie string) error {
	return s.repo.DeleteByOwnerCookie(ownerCookie)
}

func (s *Service) FindExpired() ([]File, error) {
	return s.repo.FindExpired()
}

func (s *Service) DeleteExpired() (count int64, err error) {
	return s.repo.DeleteExpired()
}
func (s *Service) FindFileList(limit int) ([]File, error) {
	return s.repo.GetAll(limit)
}
func (s *Service) FindFileListForUser(ownerCookie string, limit int) ([]File, error) {
	return s.repo.GetAllByOwnerCookie(ownerCookie, limit)
}
