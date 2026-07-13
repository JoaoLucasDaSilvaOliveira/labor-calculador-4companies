package sqlite

import (
	"database/sql"
	"fmt"
	"labor-calculador-4companies/internal/infra/config"
	"os"
	"path"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDB(cfg *config.SQLiteConfig) (*gorm.DB, error) {
	dbDir := os.Getenv("DB_DIR")
	if dbDir == "" {
		dbDir = filepath.Dir(cfg.DataSourceName)
	}

	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
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

	if err := ensureEmployeeCompanyIDColumn(db); err != nil {
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

func ensureEmployeeCompanyIDColumn(db *sql.DB) error {
	rows, err := db.Query("PRAGMA table_info(employee)")
	if err != nil {
		return fmt.Errorf("erro ao inspecionar tabela employee: %w", err)
	}
	defer rows.Close()

	hasCompanyID := false
	for rows.Next() {
		var cid int
		var name, columnType string
		var notNull int
		var defaultValue any
		var pk int

		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &pk); err != nil {
			return fmt.Errorf("erro ao ler colunas de employee: %w", err)
		}

		if name == "company_id" {
			hasCompanyID = true
			break
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("erro ao inspecionar colunas de employee: %w", err)
	}

	if !hasCompanyID {
		if _, err := db.Exec("ALTER TABLE employee ADD COLUMN company_id INTEGER NOT NULL DEFAULT 0"); err != nil {
			return fmt.Errorf("erro ao adicionar company_id em employee: %w", err)
		}
	}

	if _, err := db.Exec("CREATE INDEX IF NOT EXISTS idx_employee_company_id ON employee(company_id)"); err != nil {
		return fmt.Errorf("erro ao criar índice de employee(company_id): %w", err)
	}

	return nil
}
