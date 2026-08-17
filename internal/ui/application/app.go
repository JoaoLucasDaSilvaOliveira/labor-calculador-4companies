package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/JoaoLucasDaSilvaOliveira/learning-fyne/classes/utils"
)

type MainApplication struct {
	application  fyne.App
	masterWindow fyne.Window
}

var mainWindowSize = fyne.NewSize(1280, 800)

// NewApplication creates the Fyne application shell without building pages.
func NewApplication() *MainApplication {
	application := app.NewWithID("labor.calculator.4companies.app")
	application.Settings().SetTheme(theme.LightTheme())
	window := utils.NewWindowWithSize(application, "Calculadora Trabalhista")
	// The prototype uses fixed-width fields and rows. Give the persistent
	// sidebar and workspace enough room at startup so HSplit does not have to
	// shrink the sidebar just because a details screen was opened.
	window.Resize(mainWindowSize)
	window.CenterOnScreen()
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
