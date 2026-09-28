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
func (r *Repo) GetAllByOwnerCookie(ownerCookie string) ([]File, error) {
	var files []File
	tx := r.db.Where("owner_cookie = ?", ownerCookie).Find(&files)
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

func (r *Repo) DeleteExpired() error {

	return r.db.Where("expires_at < ?", time.Now()).Delete(&File{}).Error
}
