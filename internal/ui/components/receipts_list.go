package components

import (
	"image/color"
	"strconv"

	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/ui/assets"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// NewReceiptsListComponent is the LEGO block for the receipt panel in screen
// 1.2. The row actions remain visible as in the prototype.
func NewReceiptsListComponent(receipts []*entity.Receipt, onView, onEdit, onDelete func(receiptID int)) fyne.CanvasObject {
	if len(receipts) == 0 {
		return NewEmptyState("Nenhum recibo", "Este funcionário ainda não possui recibos associados.")
	}

	receiptList := newReceiptList(receipts, onView, onEdit, onDelete)
	// List.MinSize is the default viewport minimum, not the height of one
	// receipt row. The prototype sizes the panel from the rows themselves.
	rowHeight := newReceiptRow([]receiptAction{
		{icon: assets.BinocularsIcon},
		{icon: assets.EditIcon},
		{icon: assets.DeleteDocumentIcon},
	}).MinSize().Height
	receiptCount := receiptList.Length()
	listHeight := rowHeight * float32(receiptCount)
	if receiptCount > 1 {
		listHeight += theme.Padding() * float32(receiptCount-1)
	}
	const receiptListMaxHeight float32 = 300
	if listHeight > receiptListMaxHeight {
		listHeight = receiptListMaxHeight
	}
	receiptsPanel := newReceiptPanel(receiptList, listHeight)
	panelHeight := listHeight + 70

	// Border gives its center the whole available height. Keep the LEGO panel
	// at the prototype height and let the invisible spacer consume the rest.
	return container.NewVBox(
		newFixedHeightObject(receiptsPanel, panelHeight),
		layout.NewSpacer(),
	)
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
}

type receiptAction struct {
	icon     fyne.Resource
	onTapped func(receiptID int)
}

func newReceiptRow(actionDefinitions []receiptAction) *receiptRow {
	row := &receiptRow{
		cod:         widget.NewLabel(""),
		description: widget.NewLabel(""),
	}
	row.actionsBar = container.NewHBox()

	for _, definition := range actionDefinitions {
		definition := definition
		button := widget.NewButtonWithIcon("", definition.icon, func() {
			if definition.onTapped != nil {
				definition.onTapped(row.receiptID)
			}
		})
		row.actionsBar.Add(button)
	}

	rowHeight := row.actionsBar.MinSize().Height
	row.codStack = container.NewStack(newTransparentInformationRectangle(100, rowHeight), row.cod)

	row.description.Truncation = fyne.TextTruncateEllipsis
	descriptionField := container.NewStack(newTransparentInformationRectangle(200, rowHeight), row.description)
	row.descriptionStack = container.NewBorder(
		nil,
		nil,
		uiUtils.HPadding(150),
		uiUtils.HPadding(150),
		descriptionField,
	)
	row.actionsStack = container.NewStack(newTransparentInformationRectangle(150, rowHeight), row.actionsBar)
	row.ExtendBaseWidget(row)
	return row
}

func (row *receiptRow) CreateRenderer() fyne.WidgetRenderer {
	left := container.NewHBox(uiUtils.HPadding(20), row.codStack)
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, left, row.actionsStack, row.descriptionStack))
}

func (row *receiptRow) SetReceipt(receipt *entity.Receipt) {
	row.receiptID = receipt.GetId()
	row.cod.SetText(strconv.Itoa(row.receiptID))
	row.description.SetText(receipt.GetSumaryDescription())
}

func newReceiptList(receipts []*entity.Receipt, onView, onEdit, onDelete func(receiptID int)) *widget.List {
	receiptList := widget.NewList(
		func() int { return len(receipts) },
		func() fyne.CanvasObject {
			return newReceiptRow([]receiptAction{
				{icon: assets.BinocularsIcon, onTapped: onView},
				{icon: assets.EditIcon, onTapped: onEdit},
				{icon: assets.DeleteDocumentIcon, onTapped: onDelete},
			})
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*receiptRow).SetReceipt(receipts[id])
		},
	)
	receiptList.OnSelected = func(id widget.ListItemID) { receiptList.Unselect(id) }
	return receiptList
}

func newReceiptPanel(receiptList *widget.List, listHeight float32) fyne.CanvasObject {
	headerCod := container.NewHBox(uiUtils.HPadding(20), newReceiptHeaderCell("CÓDIGO", 100))
	headerDescription := container.NewHBox(uiUtils.HPadding(150), newReceiptHeaderCell("DESCRIÇÃO", 200))
	headerActions := container.NewHBox(newReceiptHeaderCell("AÇÕES", 150))
	header := container.NewBorder(nil, nil, headerCod, headerActions, headerDescription)

	topSpacer := canvas.NewRectangle(color.Transparent)
	topSpacer.SetMinSize(fyne.NewSize(0, 12))
	content := container.NewPadded(container.NewBorder(
		container.NewVBox(topSpacer, header, widget.NewSeparator()),
		nil,
		nil,
		nil,
		newFixedHeightObject(receiptList, listHeight),
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

	legendLabel := widget.NewLabel("RECIBOS")
	legendLabel.TextStyle.Bold = true
	legend := container.NewStack(legendBackground, container.NewCenter(legendLabel))
	legend.Move(fyne.NewPos(28, -10))
	legend.Resize(fyne.NewSize(125, 24))

	return container.NewStack(frame, content, container.NewWithoutLayout(legend))
}

func newReceiptHeaderCell(text string, width float32) fyne.CanvasObject {
	label := widget.NewLabel(text)
	background := canvas.NewRectangle(color.Transparent)
	background.SetMinSize(fyne.NewSize(width, label.MinSize().Height))
	return container.NewStack(background, label)
}
