package config

import "time"

const (
	DomainName           = "lanshare.local"
	DefaultPort          = 80
	DefaultDBPath        = "data/lanshare.db"
	DefaultMigrationsDir = "migrations"
	DefaultStoragePath   = "data/storage"

	CookiesOnHTTPSOnly = false

	AdminSessionCookie = "admin_session"
	SessionCookie      = "session_id"

	// AdminSessionCookieAge In seconds
	AdminSessionCookieAge = 1 * 24 * 3600

	// SessionCookieAge In seconds
	SessionCookieAge = 365 * 24 * 3600

	FileExpirationCheckerTick = 5 * time.Minute
)
