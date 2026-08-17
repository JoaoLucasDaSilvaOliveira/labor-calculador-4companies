package components

import (
	"fmt"
	"image/color"

	"labor-calculador-4companies/internal/domain/entity"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// NewEmployeesListComponent is the LEGO block for the employee panel in
// screen 1.1. Its internal dimensions and frame follow the prototype.
func NewEmployeesListComponent(employees []*entity.Employee, onEmployeeSelected func(employeeID int)) fyne.CanvasObject {
	if len(employees) == 0 {
		return NewEmptyState("Nenhum funcionário", "Esta empresa ainda não possui funcionários associados.")
	}

	employeeList := newEmployeesList(employees, onEmployeeSelected)
	// List.MinSize is the default viewport minimum, not the height of one
	// employee row. The prototype sizes the panel from the rows themselves.
	rowHeight := newEmployeeRow().MinSize().Height
	employeeCount := employeeList.Length()
	listHeight := rowHeight * float32(employeeCount)
	if employeeCount > 1 {
		listHeight += theme.Padding() * float32(employeeCount-1)
	}
	const employeeListMaxHeight float32 = 300
	if listHeight > employeeListMaxHeight {
		listHeight = employeeListMaxHeight
	}
	employeesPanel := newEmployeesPanel(employeeList, listHeight)
	panelHeight := listHeight + 70

	// Border gives its center the whole available height. Keep the LEGO panel
	// at the prototype height and let the invisible spacer consume the rest.
	return container.NewVBox(
		newFixedHeightObject(employeesPanel, panelHeight),
		layout.NewSpacer(),
	)
}

type employeeRow struct {
	widget.BaseWidget

	codStack  *fyne.Container
	nameStack *fyne.Container
	cpfStack  *fyne.Container

	cod  *widget.Label
	name *widget.Label
	cpf  *widget.Label
}

func newEmployeeRow() *employeeRow {
	row := &employeeRow{
		cod:  widget.NewLabel(""),
		name: widget.NewLabel(""),
		cpf:  widget.NewLabel(""),
	}

	row.codStack = container.NewStack(
		newTransparentInformationRectangle(100, row.cod.MinSize().Height),
		row.cod,
	)

	row.name.Truncation = fyne.TextTruncateEllipsis
	nameField := container.NewStack(
		newTransparentInformationRectangle(200, row.name.MinSize().Height),
		row.name,
	)
	row.nameStack = container.NewBorder(
		nil,
		nil,
		uiUtils.HPadding(150),
		uiUtils.HPadding(150),
		nameField,
	)

	row.cpfStack = container.NewStack(
		newTransparentInformationRectangle(150, row.cpf.MinSize().Height),
		row.cpf,
	)

	row.ExtendBaseWidget(row)
	return row
}

func (row *employeeRow) CreateRenderer() fyne.WidgetRenderer {
	left := container.NewHBox(uiUtils.HPadding(20), row.codStack)
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, left, row.cpfStack, row.nameStack))
}

func (row *employeeRow) SetEmployee(employee *entity.Employee) {
	row.cod.SetText(fmt.Sprintf("%d", employee.GetId()))
	row.name.SetText(fmt.Sprintf("%s %s", employee.FirstName(), employee.LastName()))
	row.cpf.SetText(employee.CPF())
}

func newEmployeesList(employees []*entity.Employee, onEmployeeSelected func(employeeID int)) *widget.List {
	employeesList := widget.NewList(
		func() int { return len(employees) },
		func() fyne.CanvasObject { return newEmployeeRow() },
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*employeeRow).SetEmployee(employees[id])
		},
	)
	employeesList.OnSelected = func(id widget.ListItemID) {
		employeesList.Unselect(id)
		if onEmployeeSelected != nil {
			onEmployeeSelected(employees[id].GetId())
		}
	}
	return employeesList
}

func newEmployeesPanel(employeeList *widget.List, listHeight float32) fyne.CanvasObject {
	headerCod := container.NewHBox(uiUtils.HPadding(20), newEmployeeHeaderCell("CÓDIGO", 100))
	headerName := container.NewHBox(uiUtils.HPadding(150), newEmployeeHeaderCell("NOME COMPLETO", 200))
	headerCPF := container.NewHBox(newEmployeeHeaderCell("CPF", 150))
	header := container.NewBorder(nil, nil, headerCod, headerCPF, headerName)

	topSpacer := canvas.NewRectangle(color.Transparent)
	topSpacer.SetMinSize(fyne.NewSize(0, 12))
	content := container.NewPadded(container.NewBorder(
		container.NewVBox(topSpacer, header, widget.NewSeparator()),
		nil,
		nil,
		nil,
		newFixedHeightObject(employeeList, listHeight),
	))

	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeColor = color.Black
	frame.StrokeWidth = 1
	frame.CornerRadius = 20

	legendBackground := canvas.NewRectangle(color.NRGBA{R: 255, G: 236, B: 153, A: 255})
	legendBackground.SetMinSize(fyne.NewSize(200, 24))
	legendBackground.CornerRadius = 5
	legendBackground.StrokeWidth = 2
	legendBackground.StrokeColor = color.NRGBA{R: 240, G: 140, B: 0, A: 255}

	legendLabel := widget.NewLabel("FUNCIONÁRIOS")
	legendLabel.TextStyle.Bold = true
	legend := container.NewStack(legendBackground, container.NewCenter(legendLabel))
	legend.Move(fyne.NewPos(28, -10))
	legend.Resize(fyne.NewSize(125, 24))

	return container.NewStack(frame, content, container.NewWithoutLayout(legend))
}

func newEmployeeHeaderCell(text string, width float32) fyne.CanvasObject {
	label := widget.NewLabel(text)
	background := canvas.NewRectangle(color.Transparent)
	background.SetMinSize(fyne.NewSize(width, label.MinSize().Height))
	return container.NewStack(background, label)
}
