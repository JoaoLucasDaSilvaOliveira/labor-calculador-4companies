package main

import (
	"errors"
	"log"
	"os"

	companyUC "labor-calculador-4companies/internal/application/usecase/company"
	employeeUC "labor-calculador-4companies/internal/application/usecase/employee"
	receiptUC "labor-calculador-4companies/internal/application/usecase/receipt"
	"labor-calculador-4companies/internal/infra/config"
	"labor-calculador-4companies/internal/infra/persistence/sqlite"
	uiApplication "labor-calculador-4companies/internal/ui/application"
	"labor-calculador-4companies/internal/ui/components"
	"labor-calculador-4companies/internal/ui/navigation"

	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load("/home/dev_jao/personal-projects/labor-calculador-4companies/.env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("erro ao carregar .env: %v", err)
	}
}

func main() {
	loadEnv()

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
	employeeRepository := sqlite.NewEmployeeRepository(database)
	receiptRepository := sqlite.NewReceiptRepository(database)

	getCompaniesUseCase := companyUC.NewGetCompanyUsecase(companyRepository)
	getCompanyByIDUseCase := companyUC.NewGetCompanyByIdUsecase(companyRepository)
	updateCompanyUseCase := companyUC.NewUpdateCompanyUsecase(companyRepository)
	getEmployeesUseCase := employeeUC.NewGetEmployeeUsecase(employeeRepository)
	getEmployeeByIDUseCase := employeeUC.NewGetEmployeeByIdUsecase(employeeRepository)
	updateEmployeeUseCase := employeeUC.NewUpdateEmployeeUsecase(employeeRepository)
	getReceiptsUseCase := receiptUC.NewGetReceiptUsecase(receiptRepository)

	// UI APPLICATION, ROUTERS AND PAGE PROVIDER
	application := uiApplication.NewApplication()
	workspaceRouter := navigation.NewRouter(navigation.RouteQuickAccess)
	editSession := &navigation.EditSession{}
	companyListReload := &components.CompanyListReloadHandle{}
	routeProvider := newUIRouteProvider(uiRouteProviderDeps{
		GetCompanies:      getCompaniesUseCase,
		GetCompanyByID:    getCompanyByIDUseCase,
		UpdateCompany:     updateCompanyUseCase,
		GetEmployees:      getEmployeesUseCase,
		GetEmployeeByID:   getEmployeeByIDUseCase,
		UpdateEmployee:    updateEmployeeUseCase,
		GetReceipts:       getReceiptsUseCase,
		Window:            application.Window(),
		EditSession:       editSession,
		CompanyListReload: companyListReload,
	})

	if err := routeProvider.RegisterWorkspace(workspaceRouter); err != nil {
		log.Fatalf("erro ao registrar as rotas do workspace: %v", err)
	}
	if err := workspaceRouter.Replace(navigation.RouteQuickAccess, nil); err != nil {
		log.Fatalf("erro ao abrir o acesso rápido: %v", err)
	}

	rootRouter := navigation.NewRouter(navigation.RouteMain)
	if err := routeProvider.RegisterRoot(rootRouter, workspaceRouter); err != nil {
		log.Fatalf("erro ao registrar as rotas externas: %v", err)
	}

	// DESKTOP EVENT LOOP
	application.Run(rootRouter.View())
}
