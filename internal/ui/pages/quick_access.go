package pages

import (
	"labor-calculador-4companies/internal/ui/components"
	"labor-calculador-4companies/internal/ui/navigation"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// NewQuickAccessPage creates the root view of the workspace router.
func NewQuickAccessPage(navigator navigation.Navigator) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("ACESSO RÁPIDO", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	defaultColor := theme.Color(theme.ColorNameForeground)
	hoverColor := theme.Color(theme.ColorNamePrimary)

	addCompany := components.NewClickableText(
		"+ CADASTRAR EMPRESA",
		defaultColor,
		hoverColor,
		func() { _ = navigator.Push(navigation.RouteCompanyCreate, nil) },
	)
	newCalculation := components.NewClickableText(
		"+ NOVO CÁLCULO",
		defaultColor,
		hoverColor,
		func() { _ = navigator.Push(navigation.RouteCalculationCreate, nil) },
	)

	return container.NewCenter(container.NewVBox(title, addCompany, newCalculation))
}
