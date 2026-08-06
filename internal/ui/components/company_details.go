package components

import (
	command "labor-calculador-4companies/internal/application/command/company"
	"labor-calculador-4companies/internal/domain/entity"
	error_factory "labor-calculador-4companies/internal/domain/error"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	ErrCompanyNotFound = error_factory.NewError("Empresa não encontrada")
)

type CompanyDetails interface {
	Execute(cmd command.GetCompanyById) (*entity.Company, error)
}

func CompanyDetailsComponent() {}

func newCompanyDetails(finder CompanyDetails, cmd command.GetCompanyById) fyne.CanvasObject {
	company, err := finder.Execute(cmd)

	if err != nil {
		return widget.NewLabel(ErrCompanyNotFound.Error())
	}

	return companyDetailsContainer(company)
}

func companyDetailsContainer(company *entity.Company) fyne.CanvasObject {
	codLabel := widget.NewLabel("CÓDIGO")
	codEntry := widget.NewEntry()
	codEntry.SetText(strconv.Itoa(company.GetId()))

	nameLabel := widget.NewLabel("NOME")
	nameEntry := widget.NewEntry()
	nameEntry.SetText(company.Name())

	cnpjLabel := widget.NewLabel("CNPJ")
	cnpjEntry := widget.NewEntry()
	cnpjEntry.SetText(company.CNPJ())

	hbox := container.NewHBox(
		codLabel, codEntry,
		nameLabel, nameEntry,
		cnpjLabel, cnpjEntry,
	)

	return hbox
}
