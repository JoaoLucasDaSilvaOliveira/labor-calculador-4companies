package config

import (
	"fmt"
	"os"
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
	return &SQLiteConfig{
		DatabaseDriver: os.Getenv("DB_DRIVER"),
		DataSourceName: os.Getenv("DB_DSN"),
		MigrationsDir: os.Getenv("MIGRATIONS_DIR"),
	}
}
