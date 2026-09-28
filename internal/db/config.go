package db

import "lanshare/config"

const (
	DefaultDBPath = config.DefaultDBPath

	DefaultMigrationsDir = config.DefaultMigrationsDir
)

type Config struct {
	DBPath string

	MigrationsDir string
}
