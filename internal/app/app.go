package app

import (
	"lanshare/internal/admin"
	"lanshare/internal/api"
	"lanshare/internal/db"
	"lanshare/internal/files"
	"lanshare/internal/settings"
)

type App struct {
	DB              *db.Service
	SettingsService *settings.Service
	AdminService    *admin.Service
	FilesService    *files.Service
	Router          *api.Router
}
