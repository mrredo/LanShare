package app

import (
	"lanshare/config"
	"lanshare/internal/api"
	"lanshare/internal/auth"
	"lanshare/internal/db"
	"lanshare/internal/files"
	"lanshare/internal/settings"
)

func NewApp(args ...any) *App {
	var dbService *db.Service
	var active *settings.Settings

	for _, arg := range args {
		switch v := arg.(type) {
		case *db.Service:
			dbService = v
		case *settings.Settings:
			active = v
		}
	}

	if dbService == nil {
		dbService = db.NewService()
		if err := dbService.Start(); err != nil {
			panic(err)
		}
	}

	settingsRepo := settings.NewRepo(dbService.DB())
	settingsService := settings.NewService(settingsRepo, active)

	adminRepo := auth.NewRepo(dbService.DB())
	adminService := auth.NewService(adminRepo, settingsService)

	filesRepo := files.NewRepo(dbService.DB())
	filesService := files.NewService(filesRepo, settingsService, adminService)

	router := api.NewRouter(config.DefaultPort)
	router.InitializeMiddlewares()
	apiGroup := router.Router().Group("/api")
	authGroup := router.Router().Group("/auth")

	adminHandler := auth.NewHandler(adminService)
	adminHandler.RegisterRoutes(authGroup)

	filesHandler := files.NewHandler(filesService)
	filesHandler.RegisterRoutes(apiGroup)

	return &App{
		DB:              dbService,
		AdminService:    adminService,
		FilesService:    filesService,
		SettingsService: settingsService,
		Router:          router,
	}
}
