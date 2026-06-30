package sqlite

import (
	"database/sql"
	"fmt"
	"labor-calculador-4companies/internal/infra/config"
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDB(cfg *config.SQLiteConfig) (*gorm.DB, error) {

	if err := os.MkdirAll(os.Getenv("DB_DIR"), os.ModePerm); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório do banco: %v", err)
	}

	dbPtr, err := gorm.Open(sqlite.Open(cfg.DataSourceName), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db, err := dbPtr.DB()

	if err != nil {
		return nil, err
	}

	err = initDBMinimalSchema(db, cfg.MigrationsDir)

	if err != nil {
		return nil, err
	}

	return dbPtr, nil
}

func initDBMinimalSchema(db *sql.DB, schemaDirPath string) error {

	files, err := os.ReadDir(schemaDirPath)

	if err != nil {
		return fmt.Errorf("Erro ao ler o diretorio de schema: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && path.Ext(file.Name()) == ".sql" {
			fullPath := path.Join(schemaDirPath, file.Name())

			fmt.Printf("Executando migration: %s\n", file.Name())

			migrationCommand, err := os.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("Erro ao ler arquivo %s: %w", file.Name(), err)
			}

			if _, err := db.Exec(string(migrationCommand)); err != nil {
				return fmt.Errorf("erro ao executar %s: %w", file.Name(), err)
			}
		}
	}

	return nil
}
