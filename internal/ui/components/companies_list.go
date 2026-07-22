package components

import (
	"fmt"
	query "labor-calculador-4companies/internal/application/query/company"
	"labor-calculador-4companies/internal/application/usecase/company"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewCompaniesListComponent(uc *company.GetCompanyUsecase) *widget.List {
	//usecase call
	companySlice, err := uc.Execute(query.GetCompanyWithFilter{}) //no args because there is no filter yet

	if err != nil {
		//todo: handle error
		fmt.Println(err)
	}

	//list creation
	companiesList := widget.NewList(
		func() int { return len(companySlice) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, co fyne.CanvasObject) {
			label := co.(*widget.Label) //type assertion
			labelText := companySlice[id].Name()
			label.SetText(labelText)
		},
	)

	return companiesList
}
