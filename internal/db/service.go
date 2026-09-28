package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Service struct {
	db     *gorm.DB
	config Config
}

func NewService(configs ...Config) *Service {
	var cfg Config
	if len(configs) > 0 {
		cfg = configs[0]
	}

	return &Service{
		config: cfg,
	}
}

func (s *Service) Start() error {
	if err := s.Connect(); err != nil {
		return fmt.Errorf("starting db service: %w", err)
	}
	if err := s.Setup(); err != nil {
		return fmt.Errorf("running db setup: %w", err)
	}
	return nil
}

func (s *Service) Connect() error {
	targetPath := s.resolveDBPath()

	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating database directory '%s': %w", dir, err)
	}

	db, err := gorm.Open(sqlite.Open(targetPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return fmt.Errorf("opening sqlite connection at '%s': %w", targetPath, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("retrieving underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(1)

	// Ļauj lasīt datubāzi pat tad, kad notiek rakstīšana
	if err := db.Exec("PRAGMA journal_mode = WAL;").Error; err != nil {
		return fmt.Errorf("setting PRAGMA journal_mode: %w", err)
	}

	// SQlite atļauj eksistēt tikai pareizām saitēm starp ierakstiem
	if err := db.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		return fmt.Errorf("setting PRAGMA foreign_keys: %w", err)
	}

	// 5000 ms (5s) līdz kāda darbība tiek pārtraukta
	if err := db.Exec("PRAGMA busy_timeout = 5000;").Error; err != nil {
		return fmt.Errorf("setting PRAGMA busy_timeout: %w", err)
	}
	s.db = db
	return nil
}

func (s *Service) Setup() error {
	if s.db == nil {
		return fmt.Errorf("database connection not initialized: call Connect() first")
	}

	return s.runMigrations()
}

func (s *Service) DB() *gorm.DB {
	return s.db
}

func (s *Service) SQLDB() (*sql.DB, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return s.db.DB()
}

func (s *Service) Close() error {
	if s.db == nil {
		return nil
	}
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *Service) resolveDBPath() string {
	p := strings.TrimSpace(s.config.DBPath)
	if p == "" {
		if fi, err := os.Stat("/data"); err == nil && fi.IsDir() {
			p = "/data/lanshare.db"
		} else {
			p = DefaultDBPath
		}
	}

	if strings.HasSuffix(p, "/") || strings.HasSuffix(p, "\\") || filepath.Ext(p) == "" {
		p = filepath.Join(p, "lanshare.db")
	}

	return p
}

func (s *Service) resolveMigrationsDir() string {
	dir := strings.TrimSpace(s.config.MigrationsDir)
	if dir == "" {
		dir = DefaultMigrationsDir
	}

	if filepath.IsAbs(dir) {
		return dir
	}
	if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
		return dir
	}

	candidates := []string{
		filepath.Join("..", dir),
		filepath.Join("..", "..", dir),
	}
	for _, candidate := range candidates {
		if fi, err := os.Stat(candidate); err == nil && fi.IsDir() {
			return candidate
		}
	}

	return dir
}
