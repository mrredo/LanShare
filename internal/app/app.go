package app

import (
	"lanshare/internal/api"
	"lanshare/internal/auth"
	"lanshare/internal/db"
	"lanshare/internal/files"
	"lanshare/internal/settings"
)

type App struct {
	DB              *db.Service
	SettingsService *settings.Service
	AdminService    *auth.Service
	FilesService    *files.Service
	Router          *api.Router
}
