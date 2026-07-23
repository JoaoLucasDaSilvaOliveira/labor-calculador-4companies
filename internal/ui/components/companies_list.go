package components

import (
	"fmt"
	query "labor-calculador-4companies/internal/application/query/company"
	"labor-calculador-4companies/internal/application/usecase/company"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewCompaniesListComponent(uc *company.GetCompanyUsecase) *widget.List {
	//usecase call
	companiesSlice, err := uc.Execute(query.GetCompanyWithFilter{}) //no args because there is no filter yet

	if err != nil {
		//todo: handle error
		fmt.Println(err)
	}

	//list creation
	companiesList := widget.NewList(
		func() int { return len(companiesSlice) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, co fyne.CanvasObject) {
			label := co.(*widget.Label) //type assertion
			labelText := companiesSlice[id].Name()
			label.SetText(labelText)
		},
	)

	return companiesList
}

func NewCompaniesSearchByNameListComponent(uc *company.GetCompanyUsecase, companyName string) *widget.List {
	//clean spaces
	companyNameTrimmed := strings.TrimSpace(companyName)

	if companyNameTrimmed == "" {
		return nil
	}

	//query command build
	qry := query.GetCompanyWithFilter{
		Name: companyName,
	}

	companiesSlice, err := uc.Execute(qry)

	if err != nil {
		//todo: handle error
		fmt.Println(err)
	}

	//list creation
	companiesList := widget.NewList(
		func() int { return len(companiesSlice) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, co fyne.CanvasObject) {
			label := co.(*widget.Label) //type assertion
			labelText := companiesSlice[id].Name()
			label.SetText(labelText)
		},
	)

	return companiesList
}
