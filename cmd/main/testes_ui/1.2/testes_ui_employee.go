package main

import (
	"fmt"
	"image/color"
	employeeCommand "labor-calculador-4companies/internal/application/command/employee"
	receiptQuery "labor-calculador-4companies/internal/application/query/receipt"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/ui/assets"
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
	fmt.Println("antes de criar o app")
	a := app.NewWithID("teste")
	fmt.Println("depois de criar o app")
	a.Settings().SetTheme(theme.LightTheme())

	fmt.Println("antes de criar a janela")
	w := utils.NewWindowWithSize(a, "")
	fmt.Println("depois de criar a janela")

	fmt.Println("antes de setar o conteudo")
	// w.SetContent(employeeInformationsComponent()) //provide the instances later for testing
	fmt.Println("depois setar o conteudo")

	fmt.Println("executando show and run")
	w.ShowAndRun()
}

type EmployeeFinderById interface {
	Execute(cmd employeeCommand.GetEmployeeById) (*entity.Employee, error)
}

func employeeInformationsComponent(employeeFinder EmployeeFinderById, employeeFinderCommand employeeCommand.GetEmployeeById, receiptFinder ReceiptFinder, onShowReceipt, onEditReceipt, onDeleteReceipt func(receiptID int)) *fyne.Container {
	employee, err := employeeFinder.Execute(employeeFinderCommand)

	if err != nil {
		//todo: handle error
	}
	content := container.NewStack()

	// employeeInfosDisabled := showEmployeeInformationsDisabled(employee)
	employeeInfosEnabled := showEmployeeInformationsEnabled(employee)
	content.Add(employeeInfosEnabled)

	employeeListComponent := NewReceiptListComponent(receiptFinder, onShowReceipt, onEditReceipt, onDeleteReceipt)
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

func showEmployeeInformationsDisabled(employee *entity.Employee) *fyne.Container {
	codLabel := widget.NewLabel("CÓDIGO")
	codEntry := widget.NewEntry()
	codEntry.SetText(strconv.Itoa(employee.GetId()))
	codEntry.Disable()
	codRec := canvas.NewRectangle(color.Transparent)
	codRec.SetMinSize(fyne.NewSize(50, codEntry.MinSize().Height))
	codStack := container.NewStack(codRec, codEntry)

	nameLabel := widget.NewLabel("NOME")
	nameEntry := widget.NewEntry()
	nameEntry.SetText(fmt.Sprintf("%s %s", employee.FirstName(), employee.LastName()))
	nameEntry.Disable()
	nameRec := canvas.NewRectangle(color.Transparent)
	nameRec.SetMinSize(fyne.NewSize(300, nameEntry.MinSize().Height))
	nameStack := container.NewStack(nameRec, nameEntry)

	cpfLabel := widget.NewLabel("CPF")
	cpfEntry := widget.NewEntry()
	cpfEntry.SetText(employee.CPF())
	cpfEntry.Disable()
	cpfRec := canvas.NewRectangle(color.Transparent)
	cpfRec.SetMinSize(fyne.NewSize(150, cpfEntry.MinSize().Height))
	cpfStack := container.NewStack(cpfRec, cpfEntry)

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
		cpfLabel,
		cpfStack,
	)

	return container.NewBorder(
		nil,
		nil,
		employeeLeft,
		employeeRight,
		employeeName,
	)
}

func showEmployeeInformationsEnabled(employee *entity.Employee) *fyne.Container {
	codLabel := widget.NewLabel("CÓDIGO")
	codValue := widget.NewLabel(strconv.Itoa(employee.GetId()))
	codRec := newRectangleForEmployeeInformations(utils.NewColor(240, 240, 240, 255), 50, codValue.MinSize().Height)
	codStack := container.NewStack(codRec, codValue)

	nameLabel := widget.NewLabel("NOME")
	nameValue := widget.NewLabel(fmt.Sprintf("%s %s", employee.FirstName(), employee.LastName()))
	nameValue.Truncation = fyne.TextTruncateEllipsis
	nameRec := newRectangleForEmployeeInformations(utils.NewColor(240, 240, 240, 255), 300, nameValue.MinSize().Height)
	nameStack := container.NewStack(nameRec, nameValue)

	cpfLabel := widget.NewLabel("CPF")
	cpfValue := widget.NewLabel(employee.CPF())
	cpfRec := newRectangleForEmployeeInformations(utils.NewColor(240, 240, 240, 255), 150, cpfValue.MinSize().Height)
	cpfStack := container.NewStack(cpfRec, cpfValue)

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
		cpfLabel,
		cpfStack,
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

	receiptID        int
	codStack         *fyne.Container
	descriptionStack *fyne.Container
	actionsStack     *fyne.Container
	actionsBar       *fyne.Container

	cod         *widget.Label
	description *widget.Label
	actions     []*widget.Button
}

type receiptAction struct {
	icon     fyne.Resource
	onTapped func(receiptID int)
}

func newReceiptRow(actionDefinitions []receiptAction) *receiptRow {
	row := &receiptRow{
		cod:         widget.NewLabel(""),
		description: widget.NewLabel(""),
		actions:     make([]*widget.Button, 0, len(actionDefinitions)),
	}
	row.actionsBar = container.NewHBox()

	for _, definition := range actionDefinitions {
		action := definition
		button := widget.NewButtonWithIcon("", action.icon, func() {
			action.onTapped(row.receiptID)
		})

		row.actions = append(row.actions, button)
		row.actionsBar.Add(button)
	}

	rowHeight := row.actionsBar.MinSize().Height

	row.codStack = container.NewStack(
		newRectangleForEmployeesInformations(
			100,
			rowHeight,
		),
		row.cod,
	)

	row.description.Truncation = fyne.TextTruncateEllipsis
	descriptionField := container.NewStack(
		newRectangleForEmployeesInformations(
			200,
			rowHeight,
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
			rowHeight,
		),
		row.actionsBar,
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
	row.receiptID = receipt.GetId()
	row.cod.SetText(strconv.Itoa(row.receiptID))
	row.description.SetText(receipt.GetSumaryDescription())
}

type ReceiptFinder interface {
	Execute(qry receiptQuery.GetReceiptWithFilter) ([]*entity.Receipt, error)
}

func NewReceiptListComponent(finder ReceiptFinder, onViewReceipt, onEditReceipt, OnDeleteReceipt func(receiptID int)) fyne.CanvasObject {
	receiptList := newReceiptList(finder, receiptQuery.GetReceiptWithFilter{}, onViewReceipt, onEditReceipt, OnDeleteReceipt)

	rowHeight := receiptList.MinSize().Height
	employeeCount := receiptList.Length()

	listHeight := rowHeight * float32(employeeCount)

	if employeeCount > 1 {
		listHeight += theme.Padding() * float32(employeeCount-1)
	}
	var employeeListMaxHeight float32 = 300

	if listHeight > employeeListMaxHeight {
		listHeight = employeeListMaxHeight
	}

	employeesPanel := newReceiptPanel(receiptList)
	panelHeight := listHeight + 60

	panelHeightSpacer := canvas.NewRectangle(color.Transparent)
	panelHeightSpacer.SetMinSize(fyne.NewSize(0, panelHeight))

	return container.NewStack(
		panelHeightSpacer,
		employeesPanel,
	)
}

func newReceiptList(finder ReceiptFinder, filter receiptQuery.GetReceiptWithFilter, onView, onEdit, onDelete func(receiptID int)) *widget.List {
	receipts, err := finder.Execute(filter)

	if err != nil {
		//todo: handle this later, maybe an err popup/dialog
	}
	if len(receipts) == 0 {
		//todo: this is not an error, but its prefearable to display a empty message
	}
	receiptList := widget.NewList(
		func() int {
			return len(receipts)
		},
		func() fyne.CanvasObject {
			return newReceiptRow([]receiptAction{
				{
					icon:     assets.BinocularsIcon,
					onTapped: onView,
				},
				{
					icon:     assets.EditIcon,
					onTapped: onEdit,
				},
				{
					icon:     assets.DeleteDocumentIcon,
					onTapped: onDelete,
				},
			})
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			eRow := object.(*receiptRow)
			eRow.SetReceipt(receipts[id])
		},
	)
	receiptList.OnSelected = func(id widget.ListItemID) { // prevent default
		receiptList.Unselect(id)
	}
	return receiptList
}

func newReceiptPanel(receiptList *widget.List) fyne.CanvasObject {
	headerCod := container.NewHBox(
		uiUtils.HPadding(20),
		newReceiptHeaderCell("CÓDIGO", 100),
	)
	headerDescription := container.NewHBox(
		uiUtils.HPadding(150),
		newReceiptHeaderCell("DESCRIÇÃO", 200),
	)
	headerActions := container.NewHBox(
		newReceiptHeaderCell("AÇÕES", 150),
	)
	header := container.NewBorder(
		nil,
		nil,
		headerCod,
		headerActions,
		headerDescription,
	)

	topSpacer := canvas.NewRectangle(color.Transparent)
	topSpacer.SetMinSize(fyne.NewSize(0, 12))

	content := container.NewPadded(
		container.NewBorder(
			container.NewVBox(topSpacer, header, widget.NewSeparator()),
			nil,
			nil,
			nil,
			receiptList,
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

func newReceiptHeaderCell(text string, width float32) fyne.CanvasObject {
	label := widget.NewLabel(text)

	background := canvas.NewRectangle(color.Transparent)
	background.SetMinSize(fyne.NewSize(width, label.MinSize().Height))

	return container.NewStack(background, label)
}
