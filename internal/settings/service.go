package settings

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Service struct {
	repo    *Repo
	current *Settings
}

func NewService(repo *Repo, initial ...*Settings) *Service {
	s := &Service{
		repo: repo,
	}
	if len(initial) > 0 && initial[0] != nil {
		s.current = initial[0]
	} else if repo != nil && repo.db != nil {
		if latest, err := repo.GetLatestSettings(); err == nil && latest != nil {
			s.current = latest
		}
	}
	return s
}

func (s *Service) Get() *Settings {
	if s.current == nil && s.repo != nil && s.repo.db != nil {
		if latest, err := s.repo.GetLatestSettings(); err == nil {
			s.current = latest
		}
	}
	return s.current
}

func (s *Service) SetCurrent(settings *Settings) {
	s.current = settings
}

func (s *Service) GetSettings(id string) (*Settings, error) {
	return s.repo.GetSettings(id)
}

func (s *Service) GetLatestSettings() (*Settings, error) {
	settings, err := s.repo.GetLatestSettings()
	if err == nil && settings != nil {
		s.current = settings
	}
	return settings, err
}

func (s *Service) CreateSettings(settings *Settings) error {
	if err := s.repo.Create(settings); err != nil {
		return err
	}
	s.current = settings
	return nil
}

func (s *Service) UpdateSettings(settings *Settings) error {
	if err := s.repo.Update(settings); err != nil {
		return err
	}
	s.current = settings
	return nil
}

func (s *Service) UpdateAdminPasswordHash(id string, adminPasswordHash string) error {
	if err := s.repo.UpdateAdminPasswordHash(id, adminPasswordHash); err != nil {
		return err
	}
	if s.current != nil && fmt.Sprint(s.current.ID) == id {
		s.current.AdminPasswordHash = adminPasswordHash
	}
	return nil
}

func (s *Service) SetAdminPassword(settingsID uint, newPassword string) error {
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.UpdateAdminPasswordHash(fmt.Sprint(settingsID), hash)
}
func (s *Service) IsPasswordValid(password string) bool {
	return CheckPasswordHash(password, s.current.AdminPasswordHash)
}
func (s *Service) DeleteSettings(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	if s.current != nil && fmt.Sprint(s.current.ID) == id {
		s.current = nil
	}
	return nil
}
