package components

import (
	"strings"

	query "labor-calculador-4companies/internal/application/query/company"
	"labor-calculador-4companies/internal/domain/entity"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

const (
	companiesLoadErrorMessage   = "Não foi possível carregar as empresas."
	companiesEmptyMessage       = "Nenhuma empresa cadastrada."
	companiesSearchEmptyMessage = "Nenhuma empresa encontrada para essa busca."
)

// CompanyFinder is the query contract required by the company list.
// GetCompanyUsecase implements this interface without an adapter.
type CompanyFinder interface {
	Execute(query.GetCompanyWithFilter) ([]*entity.Company, error)
}

// NewCompaniesListComponent creates the complete company list.
func NewCompaniesListComponent(finder CompanyFinder, onCompanySelected func(companyID int)) fyne.CanvasObject {
	return newCompaniesList(
		finder,
		query.GetCompanyWithFilter{},
		companiesEmptyMessage,
		onCompanySelected,
	)
}

// NewCompaniesSearchByNameListComponent creates a company list filtered by name.
func NewCompaniesSearchByNameListComponent(finder CompanyFinder, companyName string, onCompanySelected func(companyID int)) fyne.CanvasObject {
	trimmedName := strings.TrimSpace(companyName)
	if trimmedName == "" {
		return newStatusLabel("Digite um nome para buscar.")
	}

	return newCompaniesList(
		finder,
		query.GetCompanyWithFilter{Name: trimmedName},
		companiesSearchEmptyMessage,
		onCompanySelected,
	)
}

func newCompaniesList(finder CompanyFinder, filter query.GetCompanyWithFilter, emptyMessage string, onCompanySelected func(companyID int)) fyne.CanvasObject {
	companies, err := finder.Execute(filter)
	if err != nil {
		return newStatusLabel(companiesLoadErrorMessage)
	}
	if len(companies) == 0 {
		return newStatusLabel(emptyMessage)
	}

	companiesList := widget.NewList(
		func() int {
			return len(companies)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			label.SetText(companies[id].Name())
		},
	)
	companiesList.OnSelected = func(id widget.ListItemID) {
		// The page is kept in router history. Clearing the selection allows the
		// same company to be opened again after the user navigates back.
		companiesList.Unselect(id)
		if onCompanySelected != nil {
			onCompanySelected(companies[id].GetId())
		}
	}

	return companiesList
}

func newStatusLabel(message string) *widget.Label {
	label := widget.NewLabel(message)
	label.Wrapping = fyne.TextWrapWord
	return label
}
