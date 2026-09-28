package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDBService_LifecycleAndMigrations(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "lanshare_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbDir := filepath.Join(tempDir, "data")
	migrationsDir := filepath.Join(tempDir, "migrations")

	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		t.Fatalf("failed to create migrations dir: %v", err)
	}

	// Create test migration files
	v1Content := `
CREATE TABLE files (
    id TEXT PRIMARY KEY,
    filename TEXT,
    size INTEGER
);
`
	v2Content := `
-- Comment line
CREATE TABLE settings (
    id INTEGER PRIMARY KEY,
    storage_limit INTEGER NOT NULL
);

INSERT INTO settings (id, storage_limit) VALUES (1, 1048576);
`
	v10Content := `
CREATE TABLE admin_sessions (
    id TEXT PRIMARY KEY,
    expires_at DATETIME
);
`

	if err := os.WriteFile(filepath.Join(migrationsDir, "v1_start.sql"), []byte(v1Content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(migrationsDir, "V2_create_settings.sql"), []byte(v2Content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(migrationsDir, "V10_admin_sessions.sql"), []byte(v10Content), 0644); err != nil {
		t.Fatal(err)
	}

	// Initialize service with db directory (lanshare.db should be automatically appended and directory created)
	svc := NewService(Config{
		DBPath:        dbDir,
		MigrationsDir: migrationsDir,
	})

	// Start service (Connect + Setup)
	if err := svc.Start(); err != nil {
		t.Fatalf("svc.Start() failed: %v", err)
	}
	defer svc.Close()

	db := svc.DB()
	if db == nil {
		t.Fatal("expected db instance, got nil")
	}

	// Check if migrations table was created and has 3 records
	var migrations []Migration
	if err := db.Order("version asc").Find(&migrations).Error; err != nil {
		t.Fatalf("failed querying migrations table: %v", err)
	}

	if len(migrations) != 3 {
		t.Fatalf("expected 3 migrations, got %d", len(migrations))
	}

	if migrations[0].Version != 1 || migrations[0].Name != "start" || migrations[0].Filename != "v1_start.sql" {
		t.Errorf("unexpected migration 0: %+v", migrations[0])
	}
	if migrations[1].Version != 2 || migrations[1].Name != "create_settings" || migrations[1].Filename != "V2_create_settings.sql" {
		t.Errorf("unexpected migration 1: %+v", migrations[1])
	}
	if migrations[2].Version != 10 || migrations[2].Name != "admin_sessions" || migrations[2].Filename != "V10_admin_sessions.sql" {
		t.Errorf("unexpected migration 2: %+v", migrations[2])
	}

	// Verify migrated tables exist
	var count int64
	if err := db.Table("files").Count(&count).Error; err != nil {
		t.Errorf("files table does not exist or query failed: %v", err)
	}

	var settingCount int64
	if err := db.Table("settings").Count(&settingCount).Error; err != nil {
		t.Errorf("settings table does not exist or query failed: %v", err)
	}
	if settingCount != 1 {
		t.Errorf("expected 1 setting row, got %d", settingCount)
	}

	// Running Setup() again should be idempotent (no re-executions)
	if err := svc.Setup(); err != nil {
		t.Fatalf("second svc.Setup() failed: %v", err)
	}

	var migrationsAfter []Migration
	if err := db.Find(&migrationsAfter).Error; err != nil {
		t.Fatal(err)
	}
	if len(migrationsAfter) != 3 {
		t.Fatalf("expected still 3 migrations, got %d", len(migrationsAfter))
	}
}

func TestParseMigrationFilename(t *testing.T) {
	tests := []struct {
		filename    string
		wantOK      bool
		wantVersion int
		wantName    string
	}{
		{"v1_start.sql", true, 1, "start"},
		{"v2_create_users.sql", true, 2, "create_users"},
		{"V10_complex_name_with_underscores.sql", true, 10, "complex_name_with_underscores"},
		{"invalid.sql", false, 0, ""},
		{"V_start.sql", false, 0, ""},
		{"V1.sql", false, 0, ""},
		{"V1_test.txt", false, 0, ""},
	}

	for _, tt := range tests {
		mf, ok := parseMigrationFilename(tt.filename)
		if ok != tt.wantOK {
			t.Errorf("parseMigrationFilename(%q) ok = %v, want %v", tt.filename, ok, tt.wantOK)
			continue
		}
		if ok {
			if mf.Version != tt.wantVersion {
				t.Errorf("parseMigrationFilename(%q) version = %d, want %d", tt.filename, mf.Version, tt.wantVersion)
			}
			if mf.Name != tt.wantName {
				t.Errorf("parseMigrationFilename(%q) name = %q, want %q", tt.filename, mf.Name, tt.wantName)
			}
		}
	}
}
