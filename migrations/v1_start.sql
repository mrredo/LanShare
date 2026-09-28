CREATE TABLE IF NOT EXISTS files (
    id TEXT PRIMARY KEY NOT NULL,
    filename TEXT,
    storage_path TEXT UNIQUE,
    size INTEGER CHECK (size >= 0),
    owner_cookie TEXT NOT NULL,
    text TEXT,
    url TEXT,
    uploaded_at DATETIME NOT NULL,
    expires_at DATETIME
    );

CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY,
    storage_path TEXT NOT NULL UNIQUE,
    storage_limit INTEGER NOT NULL CHECK (storage_limit > 0),
    max_file_size INTEGER NOT NULL CHECK (max_file_size > 0),
    uploads_enabled BOOLEAN NOT NULL,
    created_at DATETIME NOT NULL,
    default_expiry INTEGER CHECK (default_expiry > 0),
    admin_password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS admin_sessions (
    id TEXT PRIMARY KEY NOT NULL,
    created_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL
);