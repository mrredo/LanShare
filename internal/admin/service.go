package admin

import (
	"crypto/rand"
	"encoding/hex"
	"lanshare/internal/settings"
	"time"
)

type Service struct {
	repo            *Repo
	settingsService *settings.Service
}

func NewService(repo *Repo, settingsService *settings.Service) *Service {
	return &Service{
		repo:            repo,
		settingsService: settingsService,
	}
}
func (s *Service) CreateSession() (*AdminSession, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	session := &AdminSession{
		ID:        hex.EncodeToString(bytes),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}
func (s *Service) IsSessionValid(sessionId string) bool {
	session, err := s.GetById(sessionId)
	// Neeksistē sesija
	if err != nil {
		return false
	}
	// Vai ir beidzies termiņš
	if time.Now().After(session.ExpiresAt) {
		return false
	}
	return true
}
func (s *Service) GetById(id string) (*AdminSession, error) {
	return s.repo.GetById(id)
}

func (s *Service) Create(session *AdminSession) error {
	return s.repo.Create(session)
}

func (s *Service) DeleteById(id string) error {
	return s.repo.DeleteById(id)
}

func (s *Service) DeleteExpired() error {
	return s.repo.DeleteExpired()
}
