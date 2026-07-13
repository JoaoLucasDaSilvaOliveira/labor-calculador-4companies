package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type SupabaseConfig struct {
	URL string
}

func LoadSupabaseConfig() (*SupabaseConfig, error) {
	dsn := os.Getenv("SUPABASE_URL")

	if dsn == "" {
		return nil, fmt.Errorf("a variável de ambiente url do supabase não foi definida")
	}

	return &SupabaseConfig{URL: dsn}, nil
}

//-----------------------------------------------------

type SQLiteConfig struct {
	DatabaseDriver, DataSourceName, MigrationsDir string
}

func LoadSQLiteConfig() *SQLiteConfig {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = filepath.Join("data", "app.db") + "?_foreign_keys=on"
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = filepath.Join("internal", "infra", "persistence", "sqlite", "migrations")
	}

	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite3"
	}

	return &SQLiteConfig{
		DatabaseDriver: driver,
		DataSourceName: dsn,
		MigrationsDir:  migrationsDir,
	}
}
