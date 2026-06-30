package main

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	companyUsecase "labor-calculador-4companies/internal/application/usecase/company"
	"labor-calculador-4companies/internal/domain/repository"
	"labor-calculador-4companies/internal/infra/config"
	"labor-calculador-4companies/internal/infra/persistence/sqlite"
	"labor-calculador-4companies/internal/ui"
)

type repositories struct {
	company repository.CompanyRepository
}

type usecases struct {
	company ui.CompanyWindowUsecases
}

func main() {
	loadEnv()

	cfg := config.LoadSQLiteConfig()
	db, err := sqlite.OpenDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer closeDatabase(db)

	repositories := buildRepositories(db)
	usecases := buildUsecases(repositories)
	window := buildWindow(usecases)
	window.ShowAndRun()
}

func buildRepositories(db *gorm.DB) repositories {
	return repositories{
		company: sqlite.NewCompanyRepository(db),
	}
}

func buildUsecases(repositories repositories) usecases {
	return usecases{
		company: ui.CompanyWindowUsecases{
			Create: companyUsecase.NewCreateCompanyUsecase(repositories.company),
			Update: companyUsecase.NewUpdateCompanyUsecase(repositories.company),
			Delete: companyUsecase.NewDeleteCompanyUsecase(repositories.company),
			Get:    companyUsecase.NewGetCompanyUsecase(repositories.company),
		},
	}
}

func buildWindow(usecases usecases) fyne.Window {
	desktopApp := app.NewWithID("labor-calculador-4companies.company")
	return ui.NewCompanyWindow(desktopApp, usecases.company)
}

func loadEnv() {
	envPath, err := findEnvFile()
	if err != nil {
		log.Fatal(err)
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Fatal(err)
	}
}

func findEnvFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func closeDatabase(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("erro ao recuperar conexão do banco: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		log.Printf("erro ao fechar banco: %v", err)
	}
}
