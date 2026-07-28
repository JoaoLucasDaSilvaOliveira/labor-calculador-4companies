package pages

import (
	"labor-calculador-4companies/internal/ui/navigation"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewCalculationCreationPage creates the initial destination for a new calculation.
func NewCalculationCreationPage(navigator navigation.Navigator) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(
		"Novo cálculo",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	description := widget.NewLabel(
		"Selecione uma empresa e informe os dados necessários para o cálculo.",
	)
	description.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		title,
		description,
		widget.NewSeparator(),
		widget.NewLabel("O fluxo do cálculo será implementado nesta página."),
	)

	return newPageWithBackAction(content, navigator)
}
