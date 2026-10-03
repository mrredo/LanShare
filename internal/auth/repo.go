package auth

import (
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetById(id string) (*AdminSession, error) {
	var session AdminSession
	tx := r.db.Where("id = ?", id).First(&session)
	return &session, tx.Error
}

func (r *Repo) Create(entity *AdminSession) error {
	return r.db.Create(entity).Error
}

func (r *Repo) DeleteById(id string) error {
	return r.db.Where("id = ?", id).Delete(&AdminSession{}).Error
}

func (r *Repo) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&AdminSession{}).Error
}
