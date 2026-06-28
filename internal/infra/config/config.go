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
