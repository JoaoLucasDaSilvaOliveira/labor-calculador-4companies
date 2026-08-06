package main

import (
	"fmt"
	"image/color"
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

var displayInfoDisable = true

func main() {
	a := app.NewWithID("teste")

	w := utils.NewWindowWithSize(a, "")

	w.SetContent(employeeInformatiosComponent())

	w.ShowAndRun()
}

func employeeInformatiosComponent() *fyne.Container {
	content := container.NewStack()

	// employeeInfosDisabled := showEmployeeInformationsDisabled()
	employeeInfosEnabled := showEmployeeInformationsEnabled()
	content.Add(employeeInfosEnabled)

	employeeListComponent := NewEmployeeListComponent()
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
				widget.NewButton("VOLTAR", nil),
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

func showEmployeeInformationsDisabled() *fyne.Container {
	codLabel := widget.NewLabel("CÓDIGO")
	codEntry := widget.NewEntry()
	codEntry.SetText("1")
	codEntry.Disable()
	codRec := canvas.NewRectangle(color.Transparent)
	codRec.SetMinSize(fyne.NewSize(50, codEntry.MinSize().Height))
	codStack := container.NewStack(codRec, codEntry)

	nameLabel := widget.NewLabel("NOME")
	nameEntry := widget.NewEntry()
	nameEntry.SetText("Fulano de Tal Pinguelo")
	nameEntry.Disable()
	nameRec := canvas.NewRectangle(color.Transparent)
	nameRec.SetMinSize(fyne.NewSize(300, codEntry.MinSize().Height))
	nameStack := container.NewStack(nameRec, nameEntry)

	cnpjLabel := widget.NewLabel("CPF")
	cnpjEntry := widget.NewEntry()
	cnpjEntry.SetText("048.511.110-10")
	cnpjEntry.Disable()
	cnpjRec := canvas.NewRectangle(color.Transparent)
	cnpjRec.SetMinSize(fyne.NewSize(150, codEntry.MinSize().Height))
	cnpjStack := container.NewStack(cnpjRec, cnpjEntry)

	employeeLeft := container.NewHBox(
		codLabel,
		codStack,
		uiUtils.HPadding(10),
	)

	employeeName := container.NewBorder(
		nil,
		nil,
		nameLabel,
		nil,
		nameStack,
	)

	employeeRight := container.NewHBox(
		uiUtils.HPadding(10),
		cnpjLabel,
		cnpjStack,
	)

	return container.NewBorder(
		nil,
		nil,
		employeeLeft,
		employeeRight,
		employeeName,
	)
}

func showEmployeeInformationsEnabled() *fyne.Container {
	codLabel := widget.NewLabel("CÓDIGO")
	codValue := widget.NewLabel("1")
	codRec := newRectangleForEmployeeInformations(utils.NewColor(240, 240, 240, 255), 50, codValue.MinSize().Height)
	codStack := container.NewStack(codRec, codValue)

	nameLabel := widget.NewLabel("NOME")
	nameValue := widget.NewLabel("Fulano de Tal Pinguelo")
	nameValue.Truncation = fyne.TextTruncateEllipsis
	nameRec := newRectangleForEmployeeInformations(utils.NewColor(240, 240, 240, 255), 300, nameValue.MinSize().Height)
	nameStack := container.NewStack(nameRec, nameValue)

	cnpjLabel := widget.NewLabel("CPF")
	cnpjValue := widget.NewLabel("048.511.110-10")
	cnpjRec := newRectangleForEmployeeInformations(utils.NewColor(240, 240, 240, 255), 150, cnpjValue.MinSize().Height)
	cnpjStack := container.NewStack(cnpjRec, cnpjValue)

	employeeLeft := container.NewHBox(
		codLabel,
		codStack,
		uiUtils.HPadding(10),
	)

	employeeName := container.NewBorder(
		nil,
		nil,
		nameLabel,
		nil,
		nameStack,
	)

	employeeRight := container.NewHBox(
		uiUtils.HPadding(10),
		cnpjLabel,
		cnpjStack,
	)

	return container.NewBorder(
		nil,
		nil,
		employeeLeft,
		employeeRight,
		employeeName,
	)
}

func newRectangleForEmployeeInformations(fillColor color.Color, width float32, height float32) *canvas.Rectangle {
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

type receiptRow struct {
	widget.BaseWidget

	codStack  *fyne.Container
	descriptionStack *fyne.Container
	actionsStack  *fyne.Container

	cod  *widget.Label
	description *widget.Label
	actions  []*widget.Button
}

func newReceiptRow() *receiptRow {
	row := &receiptRow{
		cod:  widget.NewLabel(""),
		description: widget.NewLabel(""),
		actions:  make([]*widget.Button, 0),
	}

	row.codStack = container.NewStack(
		newRectangleForEmployeesInformations(
			100,
			row.cod.MinSize().Height,
		),
		row.cod,
	)

	row.description.Truncation = fyne.TextTruncateEllipsis
	descriptionField := container.NewStack(
		newRectangleForEmployeesInformations(
			200,
			row.description.MinSize().Height,
		),
		row.description,
	)
	row.descriptionStack = container.NewBorder(
		nil,
		nil,
		uiUtils.HPadding(150),
		uiUtils.HPadding(150),
		descriptionField,
	)

	row.actionsStack = container.NewStack(
		newRectangleForEmployeesInformations(
			150,
			row.actions.MinSize().Height, //todo: resolve
		),
		row.actions,
	)

	row.ExtendBaseWidget(row)
	return row
}

func (row *receiptRow) CreateRenderer() fyne.WidgetRenderer {
	left := container.NewHBox(
		uiUtils.HPadding(20),
		row.codStack,
	)

	return widget.NewSimpleRenderer(
		container.NewBorder(nil, nil, left, row.actionsStack, row.descriptionStack),
	)
}

func (row *receiptRow) SetReceipt(receipt *entity.Receipt) {
	row.cod.SetText(strconv.Itoa(receipt.GetId.GetId()))
	
	row.name.SetText(receipt.)
	row.cpf.SetText(employee.CPF())
}

func NewEmployeeListComponent() fyne.CanvasObject {
	f1, _ := entity.LoadEmployee(1, 1, "Jose", "PingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPingueloPinguelo", "04851111010")
	f2, _ := entity.LoadEmployee(2, 1, "Maria", "Pinguelo", "04851111010")
	f3, _ := entity.LoadEmployee(3, 1, "leo", "Pinguelo", "04851111010")
	f4, _ := entity.LoadEmployee(4, 1, "John", "Pinguelo", "04851111010")
	f5, _ := entity.LoadEmployee(4, 1, "John", "Pinguelo", "04851111010")
	f6, _ := entity.LoadEmployee(4, 1, "John", "Pinguelo", "04851111010")
	f7, _ := entity.LoadEmployee(4, 1, "John", "Pinguelo", "04851111010")
	f8, _ := entity.LoadEmployee(4, 1, "John", "Pinguelo", "04851111010")
	f9, _ := entity.LoadEmployee(4, 1, "John", "Pinguelo", "04851111010")
	fA, _ := entity.LoadEmployee(4, 1, "John", "Pinguelo", "04851111010")
	fB, _ := entity.LoadEmployee(4, 1, "John", "Pinguelo", "04851111010")
	fC, _ := entity.LoadEmployee(4, 1, "John", "Pinguelo", "04851111010")
	employees := []*entity.Employee{
		f1, f2, f3, f4, f5, f6, f7, f8,
		f9,
		fA,
		fB,
		fC,
	}

	employeeList := newEmployeesList(employees, func(employeeID int) {
		fmt.Println(employees[employeeID]) //pra esse teste ta ruim pq estamos tratando do array, porém na vida real isso aqui vai virar um consulta no banco de dados
	})

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

func newEmployeesPanel(employeeList *widget.List) fyne.CanvasObject {
	headerCod := container.NewHBox(
		uiUtils.HPadding(20),
		newEmployeeHeaderCell("CÓDIGO", 100),
	)
	headerName := container.NewHBox(
		uiUtils.HPadding(150),
		newEmployeeHeaderCell("DESCRIÇÃO", 200),
	)
	headerCpf := container.NewHBox(
		newEmployeeHeaderCell("AÇÕES", 150),
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

	legendLabel := widget.NewLabel("RECIBOS")
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

func newEmployeesList(employees []*entity.Employee, onEmployeeSelected func(employeeID int)) *widget.List {
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
