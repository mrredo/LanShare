package app

import (
	"path/filepath"
	"testing"

	"lanshare/internal/db"
	"lanshare/internal/settings"
)

func TestAppWiring(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_app.db")

	database := db.NewService(db.Config{
		DBPath:        dbPath,
		MigrationsDir: "../../migrations",
	})
	if err := database.Start(); err != nil {
		t.Fatalf("failed to start test db: %v", err)
	}
	defer func() { _ = database.Close() }()

	hash, _ := settings.HashPassword("adminPassword")
	activeSettings := settings.DefaultSettings(hash)
	activeSettings.StoragePath = "./test_storage"

	// Save into DB
	settingsRepo := settings.NewRepo(database.DB())
	if err := settingsRepo.Create(activeSettings); err != nil {
		t.Fatalf("failed to create settings: %v", err)
	}

	// Wire app
	application := NewApp(database, activeSettings)
	if application == nil {
		t.Fatal("expected application, got nil")
	}
	if application.DB == nil {
		t.Error("expected DB initialized")
	}
	if application.SettingsService == nil {
		t.Error("expected SettingsService initialized")
	}
	if application.AdminService == nil {
		t.Error("expected AdminService initialized")
	}
	if application.FilesService == nil {
		t.Error("expected FilesService initialized")
	}

	// Verify settings are accessible from SettingsService throughout the app
	current := application.SettingsService.Get()
	if current == nil {
		t.Fatal("expected settings from SettingsService, got nil")
	}
	if current.StoragePath != "./test_storage" {
		t.Errorf("expected StoragePath './test_storage', got %s", current.StoragePath)
	}

	// Verify files service can query available space using settingsService
	if !application.FilesService.CanUploadFile(1024) {
		t.Error("expected CanUploadFile to return true for 1024 bytes")
	}
}
