package files

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

func (r *Repo) CountFilesByOwnerCookie(ownerCookie string) (int64, error) {
	var count int64
	tx := r.db.Model(&File{}).Where("owner_cookie = ?", ownerCookie).Count(&count)
	return count, tx.Error
}
func (r *Repo) CountAllFiles() (int64, error) {
	var count int64
	tx := r.db.Model(&File{}).Count(&count)
	return count, tx.Error
}
func (r *Repo) GetById(id string) (File, error) {
	var file File
	tx := r.db.Where("id = ?", id).Find(&file)
	return file, tx.Error
}
func (r *Repo) GetAllByOwnerCookie(ownerCookie string, limit int) ([]File, error) {
	var files []File
	expr := r.db.Where("owner_cookie = ?", ownerCookie)
	if limit == -1 {
		expr = expr.Limit(limit)
	}
	tx := expr.Order("uploaded_at DESC").Find(&files)
	return files, tx.Error
}
func (r *Repo) GetAll(limit int) ([]File, error) {
	var files []File
	expr := r.db
	if limit == -1 {
		expr = expr.Limit(limit)
	}
	tx := expr.Order("uploaded_at DESC").Find(&files)
	return files, tx.Error
}
func (r *Repo) Create(entity *File) error {
	return r.db.Create(entity).Error
}
func (r *Repo) DeleteAll() error {
	return r.db.Delete(&File{}).Error
}
func (r *Repo) DeleteById(id string) error {
	return r.db.Where("id = ?", id).Delete(&File{}).Error
}
func (r *Repo) DeleteByOwnerCookie(ownerCookie string) error {
	return r.db.Where("owner_cookie = ?", ownerCookie).Delete(&File{}).Error
}
func (r *Repo) FindExpired() ([]File, error) {
	var expiredFiles []File
	tx := r.db.Where("expires_at < ?", time.Now()).Find(&expiredFiles)
	return expiredFiles, tx.Error
}

func (r *Repo) DeleteExpired() (count int64, err error) {
	result := r.db.
		Where("expires_at < ?", time.Now()).
		Delete(&File{})

	return result.RowsAffected, result.Error
}
