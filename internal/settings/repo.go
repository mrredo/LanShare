package settings

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

func (r *Repo) GetSettings(id string) (*Settings, error) {
	var settings Settings
	tx := r.db.Where("id = ?", id).Find(&settings)
	return &settings, tx.Error
}
func (r *Repo) UpdateAdminPasswordHash(id string, adminPasswordHash string) error {
	return r.db.Model(&Settings{}).Where("id = ?", id).Update("admin_password_hash", adminPasswordHash).Error
}
func (r *Repo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Settings{}).Error
}
func (r *Repo) Create(entity *Settings) error {
	entity.CreatedAt = time.Now()
	return r.db.Create(entity).Error
}
func (r *Repo) Update(entity *Settings) error {
	return r.db.Save(entity).Error
}
func (r *Repo) GetLatestSettings() (*Settings, error) {
	var settings Settings
	tx := r.db.Order("created_at DESC").Limit(1).Find(&settings)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 || settings.ID == 0 {
		return nil, nil
	}
	return &settings, nil
}
