package supabase

import (
	"context"
	"labor-calculador-4companies/internal/infra/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB(cfg *config.SupabaseConfig) (*gorm.DB, error) {
	gormDbPrt, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDbPrt.DB()

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return gormDbPrt, nil
}
