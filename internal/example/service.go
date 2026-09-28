package example

import "errors"

// Service atbild tikai par biznesa loģiku.
// Šeit NAV Gin atkarību, HTTP statusu vai JSON apstrādes.
type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ListItems() ([]ExampleItem, error) {
	return s.repo.GetAll()
}

func (s *Service) GetItem(id uint) (*ExampleItem, error) {
	if id == 0 {
		return nil, errors.New("nederīgs ID")
	}
	return s.repo.GetByID(id)
}

func (s *Service) CreateItem(dto *CreateExampleDTO) (*ExampleItem, error) {
	if dto == nil || dto.Title == "" {
		return nil, errors.New("nosaukums nedrīkst būt tukšs")
	}

	item := &ExampleItem{
		Title:   dto.Title,
		Content: dto.Content,
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}
