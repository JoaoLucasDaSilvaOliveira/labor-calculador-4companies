package main

import (
	"fmt"
	"strings"

	companyUC "labor-calculador-4companies/internal/application/usecase/company"
	employeeUC "labor-calculador-4companies/internal/application/usecase/employee"
	receiptUC "labor-calculador-4companies/internal/application/usecase/receipt"
	error_factory "labor-calculador-4companies/internal/domain/error"
	"labor-calculador-4companies/internal/ui/navigation"
	"labor-calculador-4companies/internal/ui/pages"

	"fyne.io/fyne/v2"
)

var ErrOnOpenMain = error_factory.NewError("erro ao abrir a página inicial")

type uiRouteProviderDeps struct {
	GetCompanies    *companyUC.GetCompanyUsecase
	GetCompanyByID  *companyUC.GetCompanyByIdUsecase
	GetEmployees    *employeeUC.GetEmployeeUsecase
	GetEmployeeByID *employeeUC.GetEmployeeByIdUsecase
	GetReceipts     *receiptUC.GetReceiptUsecase
}

// uiRouteProvider adapts application routes to page constructors. It receives
// ready use cases and never creates repositories or owns the application shell.
type uiRouteProvider struct {
	deps               uiRouteProviderDeps
	workspaceNavigator navigation.Navigator
	workspaceView      fyne.CanvasObject
}

func newUIRouteProvider(deps uiRouteProviderDeps) *uiRouteProvider {
	return &uiRouteProvider{deps: deps}
}

func (p *uiRouteProvider) RegisterWorkspace(router *navigation.Router) error {
	p.workspaceNavigator = router

	registrations := []struct {
		route   navigation.RouteID
		factory navigation.PageFactory
	}{
		{navigation.RouteQuickAccess, p.quickAccessPage},
		{navigation.RouteCompanyDetails, p.companyPage},
		{navigation.RouteEmployeeDetails, p.employeePage},
		{navigation.RouteCompanyCreate, p.companyRegistrationPage},
		{navigation.RouteCalculationCreate, p.calculationCreationPage},
	}

	for _, registration := range registrations {
		if err := router.Register(registration.route, registration.factory); err != nil {
			return err
		}
	}

	return nil
}

func (p *uiRouteProvider) RegisterRoot(router *navigation.Router, workspaceRouter *navigation.Router) error {
	p.workspaceNavigator = workspaceRouter
	p.workspaceView = workspaceRouter.View()

	if err := router.Register(navigation.RouteMain, p.mainPage); err != nil {
		return err
	}
	if err := router.Replace(navigation.RouteMain, nil); err != nil {
		return fmt.Errorf("%w: %w", ErrOnOpenMain, err)
	}

	return nil
}

func (p *uiRouteProvider) mainPage(_ any) (fyne.CanvasObject, error) {
	return pages.NewMainPage(pages.MainPageDeps{
		Companies:     p.deps.GetCompanies,
		Navigator:     p.workspaceNavigator,
		WorkspaceView: p.workspaceView,
	}), nil
}

func (p *uiRouteProvider) quickAccessPage(_ any) (fyne.CanvasObject, error) {
	return pages.NewQuickAccessPage(p.workspaceNavigator), nil
}

func (p *uiRouteProvider) companyPage(params any) (fyne.CanvasObject, error) {
	companyParams, ok := params.(navigation.CompanyDetailsParams)
	if !ok {
		return nil, fmt.Errorf("parâmetros inválidos para a página da empresa")
	}
	if companyParams.CompanyID <= 0 {
		return nil, fmt.Errorf("ID da empresa inválido: %d", companyParams.CompanyID)
	}

	return pages.NewCompanyPage(pages.CompanyPageDeps{
		CompanyID: companyParams.CompanyID,
		Company:   p.deps.GetCompanyByID,
		Employees: p.deps.GetEmployees,
		Navigator: p.workspaceNavigator,
	}), nil
}

func (p *uiRouteProvider) employeePage(params any) (fyne.CanvasObject, error) {
	employeeParams, ok := params.(navigation.EmployeeDetailsParams)
	if !ok {
		return nil, fmt.Errorf("parâmetros inválidos para a página do funcionário")
	}
	if employeeParams.EmployeeID <= 0 {
		return nil, fmt.Errorf("ID do funcionário inválido: %d", employeeParams.EmployeeID)
	}

	companyName := strings.ToUpper(employeeParams.CompanyName)

	return pages.NewEmployeePage(pages.EmployeePageDeps{
		EmployeeID: employeeParams.EmployeeID,
		Employee:   p.deps.GetEmployeeByID,
		Receipts:   p.deps.GetReceipts,
		Navigator:  p.workspaceNavigator,
		CompanyName: companyName,
	}), nil
}

func (p *uiRouteProvider) companyRegistrationPage(_ any) (fyne.CanvasObject, error) {
	return pages.NewCompanyRegistrationPage(p.workspaceNavigator), nil
}

func (p *uiRouteProvider) calculationCreationPage(_ any) (fyne.CanvasObject, error) {
	return pages.NewCalculationCreationPage(p.workspaceNavigator), nil
}
