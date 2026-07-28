package main

import (
	"fmt"

	companyUC "labor-calculador-4companies/internal/application/usecase/company"
	"labor-calculador-4companies/internal/ui/navigation"
	"labor-calculador-4companies/internal/ui/pages"
	error_factory "labor-calculador-4companies/internal/domain/error"

	"fyne.io/fyne/v2"
)

var (
	ErrOnOpenHome = error_factory.NewError("erro ao abrir a página inicial")
)

// uiRouteProvider adapts application routes to page constructors.
// It receives ready dependencies and never creates repositories or use cases.
type uiRouteProvider struct {
	getCompanies *companyUC.GetCompanyUsecase
	navigator    navigation.Navigator
}

func newUIRouteProvider(
	getCompanies *companyUC.GetCompanyUsecase,
	navigator navigation.Navigator,
) *uiRouteProvider {
	return &uiRouteProvider{
		getCompanies: getCompanies,
		navigator:    navigator,
	}
}

func (p *uiRouteProvider) Register(router *navigation.Router) error {
	if err := router.Register(navigation.RouteHome, p.homePage); err != nil {
		return err
	}
	if err := router.Register(navigation.RouteCompanyDetails, p.companyPage); err != nil {
		return err
	}
	if err := router.Register(navigation.RouteCompanyCreate, p.companyRegistrationPage); err != nil {
		return err
	}
	if err := router.Register(navigation.RouteCalculationCreate, p.calculationCreationPage); err != nil {
		return err
	}
	if err := router.Replace(navigation.RouteHome, nil); err != nil {
		return fmt.Errorf("%w: %w", ErrOnOpenHome, err)
	}

	return nil
}

func (p *uiRouteProvider) homePage(_ any) (fyne.CanvasObject, error) {
	return pages.NewHomePage(pages.HomePageDeps{
		Companies: p.getCompanies,
		Navigator: p.navigator,
	}), nil
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
		Navigator: p.navigator,
	}), nil
}

func (p *uiRouteProvider) companyRegistrationPage(_ any) (fyne.CanvasObject, error) {
	return pages.NewCompanyRegistrationPage(p.navigator), nil
}

func (p *uiRouteProvider) calculationCreationPage(_ any) (fyne.CanvasObject, error) {
	return pages.NewCalculationCreationPage(p.navigator), nil
}
