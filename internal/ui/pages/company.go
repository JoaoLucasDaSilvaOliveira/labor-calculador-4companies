package pages

import (
	command "labor-calculador-4companies/internal/application/command/company"
	query "labor-calculador-4companies/internal/application/query/employee"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/ui/components"
	"labor-calculador-4companies/internal/ui/navigation"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CompanyFinderByID interface {
	Execute(command.GetCompanyById) (*entity.Company, error)
}

type EmployeeFinder interface {
	Execute(query.GetEmployeeWithFilter) ([]*entity.Employee, error)
}

type CompanyPageDeps struct {
	CompanyID int
	Company   CompanyFinderByID
	Employees EmployeeFinder
	Navigator navigation.Navigator
}

// NewCompanyPage composes screen 1.1 using the prototype's exact nesting,
// dimensions and bottom action area.
func NewCompanyPage(deps CompanyPageDeps) fyne.CanvasObject {
	return newAsyncPrototypeContent("Carregando informações da empresa...", func() fyne.CanvasObject {
		company, err := deps.Company.Execute(command.GetCompanyById{IDCompany: deps.CompanyID})
		if err != nil {
			return components.NewRecoverableErrorState("Não foi possível carregar as informações da empresa.", nil)
		}

		employees, err := deps.Employees.Execute(query.GetEmployeeWithFilter{CompanyID: deps.CompanyID})
		if err != nil {
			return components.NewRecoverableErrorState("Não foi possível carregar os funcionários desta empresa.", nil)
		}

		content := container.NewStack(components.NewCompanyDetailsComponent(company))
		employeeListComponent := components.NewEmployeesListComponent(employees, func(employeeID int) {
			if err := deps.Navigator.Push(
				navigation.RouteEmployeeDetails,
				navigation.EmployeeDetailsParams{EmployeeID: employeeID, CompanyName: company.Name()},
			); err != nil {
				fyne.LogError("Não foi possível abrir o funcionário selecionado.", err)
			}
		})

		btnGravar := widget.NewButton("GRAVAR", nil)
		btnGravar.Disable()
		btnVoltar := widget.NewButton("VOLTAR", func() { resetToQuickAccess(deps.Navigator) })

		employeeListComponentBorder := container.NewBorder(
			uiUtils.VPadding(50),
			container.NewVBox(
				uiUtils.VPadding(50),
				container.NewHBox(
					layout.NewSpacer(),
					btnGravar,
					uiUtils.HPadding(10),
					btnVoltar,
				),
			),
			nil,
			nil,
			employeeListComponent,
		)

		insideBorder := container.NewBorder(content, nil, nil, nil, employeeListComponentBorder)
		insideBorderWithPadding := container.NewBorder(
			uiUtils.VPadding(10),
			uiUtils.VPadding(10),
			uiUtils.HPadding(10),
			uiUtils.HPadding(10),
			insideBorder,
		)

		return insideBorderWithPadding
	})
}

func resetToQuickAccess(navigator navigation.Navigator) {
	if err := navigator.Reset(navigation.RouteQuickAccess, nil); err != nil {
		fyne.LogError("Não foi possível resetar o workspace.", err)
		return
	}
	if err := navigator.Replace(navigation.RouteQuickAccess, nil); err != nil {
		fyne.LogError("Não foi possível abrir o acesso rápido.", err)
	}
}

func NewCompanyRegistrationPage(navigator navigation.Navigator) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Cadastrar empresa", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	description := widget.NewLabel("Informe os dados da empresa para começar a organizar os cálculos.")
	description.Wrapping = fyne.TextWrapWord
	content := container.NewVBox(
		title,
		description,
		widget.NewSeparator(),
		widget.NewLabel("O formulário de cadastro será implementado nesta página."),
	)
	return newPageWithBackAction(content, navigator)
}

func pageWithBody(header, body fyne.CanvasObject) fyne.CanvasObject {
	return container.NewPadded(container.NewBorder(header, nil, nil, nil, body))
}

func newPageWithBackAction(content fyne.CanvasObject, navigator navigation.Navigator) fyne.CanvasObject {
	return pageWithBody(
		widget.NewButtonWithIcon("Voltar", theme.NavigateBackIcon(), func() { navigator.Back() }),
		container.NewPadded(content),
	)
}
