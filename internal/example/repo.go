package example

import (
	"time"

	"gorm.io/gorm"
)

// Repo atbild tikai par tiešo mijiedarbību ar datubāzi (GORM).
type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetAll() ([]ExampleItem, error) {
	var items []ExampleItem
	err := r.db.Find(&items).Error
	return items, err
}

func (r *Repo) GetByID(id uint) (*ExampleItem, error) {
	var item ExampleItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repo) Create(item *ExampleItem) error {
	item.CreatedAt = time.Now()
	return r.db.Create(item).Error
}
