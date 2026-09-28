package db

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Migration struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Version   int       `gorm:"uniqueIndex;not null" json:"version"`
	Name      string    `gorm:"not null" json:"name"`
	Filename  string    `gorm:"not null" json:"filename"`
	AppliedAt time.Time `gorm:"not null" json:"applied_at"`
}

func (Migration) TableName() string {
	return "migrations"
}

type MigrationFile struct {
	Version  int
	Name     string
	Filename string
	FullPath string
}

var migrationRegex = regexp.MustCompile(`(?i)^v(\d+)_(.+)\.sql$`)

func parseMigrationFilename(filename string) (*MigrationFile, bool) {
	// atgriež 3 vērtības
	// 0 -> faila nosaukums
	// 1 -> versija
	// 2 -> migrācijas nosaukums
	matches := migrationRegex.FindStringSubmatch(filename)

	if len(matches) < 3 {
		return nil, false
	}

	version, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, false
	}

	return &MigrationFile{
		Version:  version,
		Name:     matches[2],
		Filename: filename,
	}, true
}

func loadMigrationFiles(dir string) ([]MigrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("migrations directory '%s' not found: %w", dir, err)
		}
		return nil, fmt.Errorf("reading migrations directory '%s': %w", dir, err)
	}

	var files []MigrationFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		mf, ok := parseMigrationFilename(entry.Name())
		if !ok {
			continue
		}
		mf.FullPath = filepath.Join(dir, entry.Name())
		files = append(files, *mf)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Version < files[j].Version
	})

	return files, nil
}

func (s *Service) runMigrations() error {
	if err := s.db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	resolvedDir := s.resolveMigrationsDir()
	files, err := loadMigrationFiles(resolvedDir)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	var applied []Migration
	if err := s.db.Order("version asc").Find(&applied).Error; err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}

	appliedVersions := make(map[int]bool, len(applied))
	for _, m := range applied {
		appliedVersions[m.Version] = true
	}

	for _, file := range files {
		if appliedVersions[file.Version] {
			continue
		}

		if err := s.executeMigrationFile(file); err != nil {
			return fmt.Errorf("failed applying migration %s: %w", file.Filename, err)
		}
	}

	return nil
}

func (s *Service) executeMigrationFile(file MigrationFile) error {
	content, err := os.ReadFile(file.FullPath)
	if err != nil {
		return fmt.Errorf("reading migration file %s: %w", file.FullPath, err)
	}

	sqlText := strings.TrimSpace(string(content))

	return s.db.Transaction(func(tx *gorm.DB) error {
		if sqlText != "" {
			statements := splitSQLStatements(sqlText)
			for _, stmt := range statements {
				if strings.TrimSpace(stmt) == "" {
					continue
				}
				if err := tx.Exec(stmt).Error; err != nil {
					return fmt.Errorf("executing statement in %s: %w", file.Filename, err)
				}
			}
		}

		record := Migration{
			Version:   file.Version,
			Name:      file.Name,
			Filename:  file.Filename,
			AppliedAt: time.Now(),
		}

		if err := tx.Create(&record).Error; err != nil {
			return fmt.Errorf("recording migration %s: %w", file.Filename, err)
		}

		return nil
	})
}

func splitSQLStatements(sql string) []string {
	var statements []string
	var current strings.Builder

	inSingleQuote := false
	inDoubleQuote := false
	inLineComment := false
	inBlockComment := false

	runes := []rune(sql)
	n := len(runes)

	for i := 0; i < n; i++ {
		r := runes[i]
		next := rune(0)
		if i+1 < n {
			next = runes[i+1]
		}

		if inLineComment {
			if r == '\n' {
				inLineComment = false
			}
			current.WriteRune(r)
			continue
		}

		if inBlockComment {
			if r == '*' && next == '/' {
				inBlockComment = false
				current.WriteString("*/")
				i++
				continue
			}
			current.WriteRune(r)
			continue
		}

		if !inSingleQuote && !inDoubleQuote {
			if r == '-' && next == '-' {
				inLineComment = true
				current.WriteString("--")
				i++
				continue
			}
			if r == '/' && next == '*' {
				inBlockComment = true
				current.WriteString("/*")
				i++
				continue
			}
		}

		if r == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			current.WriteRune(r)
			continue
		}

		if r == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			current.WriteRune(r)
			continue
		}

		if r == ';' && !inSingleQuote && !inDoubleQuote {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
			continue
		}

		current.WriteRune(r)
	}

	stmt := strings.TrimSpace(current.String())
	if stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
}
