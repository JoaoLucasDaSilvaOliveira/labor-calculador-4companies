package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/JoaoLucasDaSilvaOliveira/learning-fyne/classes/utils"
	usecase "labor-calculador-4companies/internal/application/usecase/company"
)

type MainApplication struct {
	//fyne basic components
	Application  fyne.App
	MasterWindow fyne.Window
	//app usecases
	GetCompaniesUC *usecase.GetCompanyUsecase
}

// like a main func but for ui
func NewApplication() *MainApplication {
	application := app.NewWithID("labor.calculator.4companies.app")
	window := utils.NewWindowWithSize(application, "Calculadora Trabalhista")
	window.SetMaster()
	return &MainApplication{
		Application:  application,
		MasterWindow: window,
	}
}
