package main

import (
	"fmt"
	"image/color"
	command "labor-calculador-4companies/internal/application/command/company"
	employeeQuery "labor-calculador-4companies/internal/application/query/employee"
	"labor-calculador-4companies/internal/domain/entity"
	uiUtils "labor-calculador-4companies/internal/ui/utils"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/JoaoLucasDaSilvaOliveira/learning-fyne/classes/utils"
)

func main() {
	a := app.NewWithID("teste")

	w := utils.NewWindowWithSize(a, "")

	// w.SetContent(companyInformatiosComponent()) //provides instances for testing later

	w.ShowAndRun()
}

type CompanyFinderById interface {
	Execute(cmd command.GetCompanyById) (*entity.Company, error)
}

func companyInformationsComponent(companyFinder CompanyFinderById, companyFinderCommand command.GetCompanyById, employeeFinder EmployeeFinder, onEmployeeSelected func(EmployeeId int)) *fyne.Container {
	company, err := companyFinder.Execute(companyFinderCommand)

	if err != nil {
		//todo: handle error
	}
	content := container.NewStack()

	// companyInfosDisabled := showCompanyInformationsDisabled(company)
	companyInfosEnabled, btnGravar := showCompanyInformations(company)
	content.Add(companyInfosEnabled)

	employeeListComponent := NewEmployeeListComponent(employeeFinder, onEmployeeSelected)

	btnGravar.Disable()
	btnVoltar := widget.NewButton("VOLTAR", func() { /*todo: need the router objt to implement this */ })
	btnVoltar.Disable()
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

	insideBorder := container.NewBorder(
		content,
		nil,
		nil,
		nil,
		employeeListComponentBorder,
	)

	insideBorderWithPadding := container.NewBorder(
		uiUtils.VPadding(10),
		uiUtils.VPadding(10),
		uiUtils.HPadding(10),
		uiUtils.HPadding(10),
		insideBorder,
	)

	return container.NewBorder(
		container.NewStack(widget.NewButton("", nil)),
		nil,
		container.NewStack(widget.NewButton("", nil)),
		nil,
		insideBorderWithPadding,
	)
} // main component

func showCompanyInformations(company *entity.Company) (*fyne.Container, *widget.Button) {
	codLabel := widget.NewLabel("CÓDIGO")
	codField := newInlineEditableCompanyField(strconv.Itoa(company.GetId()), 50)
	codStack := container.NewStack(codField)

	nameLabel := widget.NewLabel("NOME")
	nameField := newInlineEditableCompanyField(company.Name(), 300)
	nameStack := container.NewStack(nameField)

	cnpjLabel := widget.NewLabel("CNPJ")
	cnpjField := newInlineEditableCompanyField(company.CNPJ(), 150)
	cnpjStack := container.NewStack(cnpjField)

	companyLeft := container.NewHBox(
		codLabel,
		codStack,
		uiUtils.HPadding(10),
	)

	companyName := container.NewBorder(
		nil,
		nil,
		nameLabel,
		nil,
		nameStack,
	)

	companyRight := container.NewHBox(
		uiUtils.HPadding(10),
		cnpjLabel,
		cnpjStack,
	)

	return container.NewBorder(
			nil,
			nil,
			companyLeft,
			companyRight,
			companyName,
		), widget.NewButton("GRAVAR", func() { /*todo: create separated func to handle this event: important: the func must receive usecase by injection*/
		})
}

// inlineEditableCompanyField shows a readable Label until it is clicked. It
// then replaces it with an Entry, without calling Disable() and its muted text.
type inlineEditableCompanyField struct {
	widget.BaseWidget

	value   *widget.Label
	entry   *widget.Entry
	content *fyne.Container
}

func newInlineEditableCompanyField(text string, width float32) *inlineEditableCompanyField {
	field := &inlineEditableCompanyField{
		value: widget.NewLabel(text),
		entry: widget.NewEntry(),
	}
	field.value.Truncation = fyne.TextTruncateEllipsis
	field.entry.SetText(text)
	field.entry.Hide()

	background := newRectangleForCompanyInformations(
		utils.NewColor(240, 240, 240, 255),
		width,
		field.entry.MinSize().Height,
	)
	field.content = container.NewStack(background, field.value, field.entry)
	field.ExtendBaseWidget(field)
	return field
}

func (field *inlineEditableCompanyField) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(field.content)
}

func (field *inlineEditableCompanyField) Tapped(*fyne.PointEvent) {
	if field.entry.Visible() {
		return
	}

	field.value.Hide()
	field.entry.Show()
	if app := fyne.CurrentApp(); app != nil {
		if canvas := app.Driver().CanvasForObject(field.entry); canvas != nil {
			canvas.Focus(field.entry)
		}
	}
}

func newRectangleForCompanyInformations(fillColor color.Color, width float32, height float32) *canvas.Rectangle {
	rectangle := newRectangleWithDefaultSettings(fillColor, width, height)
	rectangle.CornerRadius = 5
	rectangle.StrokeColor = color.Black
	rectangle.StrokeWidth = 0.5
	return rectangle
}

func newRectangleForEmployeesInformations(width float32, height float32) *canvas.Rectangle {
	return newRectangleWithDefaultSettings(color.Transparent, width, height)
}

func newRectangleWithDefaultSettings(fillColor color.Color, width float32, height float32) *canvas.Rectangle {
	rectangle := canvas.NewRectangle(fillColor)
	rectangle.SetMinSize(fyne.NewSize(width, height))
	return rectangle
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
		newRectangleForEmployeesInformations(
			100,
			row.cod.MinSize().Height,
		),
		row.cod,
	)

	row.name.Truncation = fyne.TextTruncateEllipsis
	nameField := container.NewStack(
		newRectangleForEmployeesInformations(
			200,
			row.name.MinSize().Height,
		),
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
		newRectangleForEmployeesInformations(
			150,
			row.cpf.MinSize().Height,
		),
		row.cpf,
	)

	row.ExtendBaseWidget(row)
	return row
}

func (row *employeeRow) CreateRenderer() fyne.WidgetRenderer {
	left := container.NewHBox(
		uiUtils.HPadding(20),
		row.codStack,
	)

	return widget.NewSimpleRenderer(
		container.NewBorder(nil, nil, left, row.cpfStack, row.nameStack),
	)
}

func (row *employeeRow) SetEmployee(employee *entity.Employee) {
	row.cod.SetText(strconv.Itoa(employee.GetId()))
	row.name.SetText(fmt.Sprintf(
		"%s %s",
		employee.FirstName(),
		employee.LastName(),
	))
	row.cpf.SetText(employee.CPF())
}

type EmployeeFinder interface {
	Execute(qry employeeQuery.GetEmployeeWithFilter) ([]*entity.Employee, error)
}

func NewEmployeeListComponent(finder EmployeeFinder, onEmployeeSelected func(employeeId int)) fyne.CanvasObject {
	employeeList := newEmployeesList(finder, employeeQuery.GetEmployeeWithFilter{}, onEmployeeSelected)

	rowHeight := employeeList.MinSize().Height
	employeeCount := employeeList.Length()

	listHeight := rowHeight * float32(employeeCount)

	if employeeCount > 1 {
		listHeight += theme.Padding() * float32(employeeCount-1)
	}
	var employeeListMaxHeight float32 = 300

	if listHeight > employeeListMaxHeight {
		listHeight = employeeListMaxHeight
	}

	employeesPanel := newEmployeesPanel(employeeList)
	panelHeight := listHeight + 60

	panelHeightSpacer := canvas.NewRectangle(color.Transparent)
	panelHeightSpacer.SetMinSize(fyne.NewSize(0, panelHeight))

	return container.NewStack(
		panelHeightSpacer,
		employeesPanel,
	)
}

func newEmployeesList(finder EmployeeFinder, filter employeeQuery.GetEmployeeWithFilter, onEmployeeSelected func(employeeID int)) *widget.List {
	employees, err := finder.Execute(filter)

	if err != nil {
		//todo: handle error
	}

	employeesList := widget.NewList(
		func() int {
			return len(employees)
		},
		func() fyne.CanvasObject {
			return newEmployeeRow()
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			eRow := object.(*employeeRow)
			eRow.SetEmployee(employees[id])
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

func newEmployeesPanel(employeeList *widget.List) fyne.CanvasObject {
	headerCod := container.NewHBox(
		uiUtils.HPadding(20),
		newEmployeeHeaderCell("CÓDIGO", 100),
	)
	headerName := container.NewHBox(
		uiUtils.HPadding(150),
		newEmployeeHeaderCell("NOME COMPLETO", 200),
	)
	headerCpf := container.NewHBox(
		newEmployeeHeaderCell("CPF", 150),
	)
	header := container.NewBorder(
		nil,
		nil,
		headerCod,
		headerCpf,
		headerName,
	)

	topSpacer := canvas.NewRectangle(color.Transparent)
	topSpacer.SetMinSize(fyne.NewSize(0, 12))

	content := container.NewPadded(
		container.NewBorder(
			container.NewVBox(topSpacer, header, widget.NewSeparator()),
			nil,
			nil,
			nil,
			employeeList,
		),
	)

	frame := canvas.NewRectangle(color.Transparent)
	frame.StrokeColor = color.Black
	frame.StrokeWidth = 1
	frame.CornerRadius = 20

	legendBackground := canvas.NewRectangle(utils.NewColor(255, 236, 153, 255))
	legendBackground.SetMinSize(fyne.NewSize(200, 24))
	legendBackground.CornerRadius = 5
	legendBackground.StrokeWidth = 2
	legendBackground.StrokeColor = utils.NewColor(240, 140, 0, 255)

	legendLabel := widget.NewLabel("FUNCIONÁRIOS")
	legendLabel.TextStyle.Bold = true
	legend := container.NewStack(
		legendBackground,
		container.NewCenter(legendLabel),
	)
	legend.Move(fyne.NewPos(28, -10))
	legend.Resize(fyne.NewSize(125, 24))

	legendLayer := container.NewWithoutLayout(legend)

	return container.NewStack(frame, content, legendLayer)
}

func newEmployeeHeaderCell(text string, width float32) fyne.CanvasObject {
	label := widget.NewLabel(text)

	background := canvas.NewRectangle(color.Transparent)
	background.SetMinSize(fyne.NewSize(width, label.MinSize().Height))

	return container.NewStack(background, label)
}
