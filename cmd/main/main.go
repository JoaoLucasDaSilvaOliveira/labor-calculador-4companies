package main

import (
	"log"

	companyUC "labor-calculador-4companies/internal/application/usecase/company"
	"labor-calculador-4companies/internal/infra/config"
	"labor-calculador-4companies/internal/infra/persistence/sqlite"
	uiApplication "labor-calculador-4companies/internal/ui/application"
	"labor-calculador-4companies/internal/ui/pages"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("/home/dev_jao/personal_projects/labor-calculator-for-companies/.env")
}

func main() {
	// CONFIG AND DATABASE
	database, err := sqlite.OpenDB(config.LoadSQLiteConfig())
	if err != nil {
		log.Fatalf("erro ao abrir banco de dados: %v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("erro ao obter conexão com banco de dados: %v", err)
	}
	defer sqlDB.Close()

	// REPOSITORIES AND USE CASES
	companyRepository := sqlite.NewCompanyRepository(database)
	getCompaniesUseCase := companyUC.NewGetCompanyUsecase(companyRepository)

	// UI
	application := uiApplication.NewApplication()
	application.GetCompaniesUC = getCompaniesUseCase
	pages.NewHomePage(application)

	application.MasterWindow.ShowAndRun()
}
