package pages

import (
	"fmt"

	"labor-calculador-4companies/internal/ui/navigation"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CompanyPageDeps contains the route data and navigation required by the page.
type CompanyPageDeps struct {
	CompanyID int
	Navigator navigation.Navigator
}

// NewCompanyPage creates the destination opened when a company is selected.
// The details use case can be added to CompanyPageDeps when this screen evolves.
func NewCompanyPage(deps CompanyPageDeps) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(
		fmt.Sprintf("Empresa #%d", deps.CompanyID),
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	description := widget.NewLabel(
		"Consulte os dados da empresa e acesse seus cálculos trabalhistas.",
	)
	description.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		title,
		description,
		widget.NewSeparator(),
		widget.NewLabel("Os detalhes da empresa serão adicionados nesta página."),
	)

	return newPageWithBackAction(content, deps.Navigator)
}

// NewCompanyRegistrationPage creates the initial destination for company creation.
func NewCompanyRegistrationPage(navigator navigation.Navigator) fyne.CanvasObject {
	title := widget.NewLabelWithStyle(
		"Cadastrar empresa",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)
	description := widget.NewLabel(
		"Informe os dados da empresa para começar a organizar os cálculos.",
	)
	description.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		title,
		description,
		widget.NewSeparator(),
		widget.NewLabel("O formulário de cadastro será implementado nesta página."),
	)

	return newPageWithBackAction(content, navigator)
}

func newPageWithBackAction(
	content fyne.CanvasObject,
	navigator navigation.Navigator,
) fyne.CanvasObject {
	backButton := widget.NewButtonWithIcon("Voltar", theme.NavigateBackIcon(), func() {
		navigator.Back()
	})
	header := container.NewHBox(backButton)

	return container.NewPadded(container.NewBorder(
		header,
		nil,
		nil,
		nil,
		container.NewPadded(content),
	))
}
