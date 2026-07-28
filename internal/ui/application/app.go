package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/JoaoLucasDaSilvaOliveira/learning-fyne/classes/utils"
)

type MainApplication struct {
	application  fyne.App
	masterWindow fyne.Window
}

// NewApplication creates the Fyne application shell without building pages.
func NewApplication() *MainApplication {
	application := app.NewWithID("labor.calculator.4companies.app")
	window := utils.NewWindowWithSize(application, "Calculadora Trabalhista")
	window.SetMaster()

	return &MainApplication{
		application:  application,
		masterWindow: window,
	}
}

// Run mounts the stable root view and starts the desktop event loop.
func (a *MainApplication) Run(rootView fyne.CanvasObject) {
	a.masterWindow.SetContent(rootView)
	a.masterWindow.ShowAndRun()
}
