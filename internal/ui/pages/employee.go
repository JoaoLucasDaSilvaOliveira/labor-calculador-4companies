package pages

import (
	command "labor-calculador-4companies/internal/application/command/employee"
	receiptQuery "labor-calculador-4companies/internal/application/query/receipt"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/ui/components"
	"labor-calculador-4companies/internal/ui/navigation"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type EmployeeFinderByID interface {
	Execute(command.GetEmployeeById) (*entity.Employee, error)
}

type ReceiptFinder interface {
	Execute(receiptQuery.GetReceiptWithFilter) ([]*entity.Receipt, error)
}

type EmployeePageDeps struct {
	EmployeeID int
	Employee   EmployeeFinderByID
	Receipts   ReceiptFinder
	Navigator  navigation.Navigator
	CompanyName string
}

// NewEmployeePage composes screen 1.2 using the prototype's exact nesting,
// dimensions, receipt actions and bottom action area.
func NewEmployeePage(deps EmployeePageDeps) fyne.CanvasObject {
	return newAsyncPrototypeContent("Carregando informações do funcionário...", func() fyne.CanvasObject {
		employee, err := deps.Employee.Execute(command.GetEmployeeById{IDEmployee: deps.EmployeeID})
		if err != nil {
			return components.NewRecoverableErrorState("Não foi possível carregar as informações do funcionário.", nil)
		}

		receipts, err := deps.Receipts.Execute(receiptQuery.GetReceiptWithFilter{IDEmployee: deps.EmployeeID})
		if err != nil {
			return components.NewRecoverableErrorState("Não foi possível carregar os recibos deste funcionário.", nil)
		}

		content := container.NewStack(components.NewEmployeeDetailsComponent(employee, deps.CompanyName))
		receiptListComponent := components.NewReceiptsListComponent(receipts, nil, nil, nil)
		employeeListComponentBorder := container.NewBorder(
			uiUtils.VPadding(50),
			container.NewVBox(
				uiUtils.VPadding(50),
				container.NewHBox(
					layout.NewSpacer(),
					widget.NewButton("NOVO RECIBO", nil),
					uiUtils.HPadding(10),
					widget.NewButton("GRAVAR", nil),
					uiUtils.HPadding(10),
					widget.NewButton("VOLTAR", func() { deps.Navigator.Back() }),
				),
			),
			nil,
			nil,
			receiptListComponent,
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
