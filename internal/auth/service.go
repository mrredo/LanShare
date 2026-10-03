package auth

import (
	"crypto/rand"
	"encoding/hex"
	"lanshare/config"
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
		ExpiresAt: time.Now().Add(config.AdminSessionCookieAge * time.Second),
	}

	if err := s.repo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) IsAdminSessionValid(adminCookie string) bool {
	cookie, err := s.GetByCookie(adminCookie)
	if err != nil {
		return false
	}
	if cookie.IsExpired() {
		return false
	}
	return true
}

func (s *Service) GetById(id string) (*AdminSession, error) {
	return s.repo.GetById(id)
}
func (s *Service) GetByCookie(cookie string) (*AdminSession, error) {
	return s.GetById(cookie)
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
