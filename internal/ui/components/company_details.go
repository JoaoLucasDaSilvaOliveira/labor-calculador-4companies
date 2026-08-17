package components

import (
	"image/color"
	"strconv"

	"labor-calculador-4companies/internal/domain/entity"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewCompanyDetailsComponent preserves the prototype composition of the
// company information row while keeping the editable fields reusable.
func NewCompanyDetailsComponent(company *entity.Company) fyne.CanvasObject {
	codLabel := widget.NewLabel("CÓDIGO")
	codStack := container.NewStack(NewInlineEditableField(strconv.Itoa(company.GetId()), 50))

	nameLabel := widget.NewLabel("NOME")
	nameStack := container.NewStack(NewInlineEditableField(company.Name(), 300))

	cnpjLabel := widget.NewLabel("CNPJ")
	cnpjStack := container.NewStack(NewInlineEditableField(company.CNPJ(), 150))

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
		widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil,
		companyLeft,
		companyRight,
		companyName,
	)
}

// InlineEditableField is the clickable field used by the company prototype.
// Persistence is intentionally left to the future save flow.
type InlineEditableField struct {
	widget.BaseWidget

	value   *widget.Label
	entry   *widget.Entry
	content *fyne.Container
}

func NewInlineEditableField(text string, width float32) *InlineEditableField {
	field := &InlineEditableField{
		value: widget.NewLabel(text),
		entry: widget.NewEntry(),
	}
	field.value.Truncation = fyne.TextTruncateEllipsis
	field.entry.SetText(text)
	field.entry.Hide()

	background := newRoundedInformationRectangle(
		color.NRGBA{R: 240, G: 240, B: 240, A: 255},
		width,
		field.entry.MinSize().Height,
	)
	field.content = container.NewStack(background, field.value, field.entry)
	field.ExtendBaseWidget(field)
	return field
}

func (field *InlineEditableField) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(field.content)
}

func (field *InlineEditableField) Tapped(*fyne.PointEvent) {
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

func newRoundedInformationRectangle(fillColor color.Color, width, height float32) *canvas.Rectangle {
	rectangle := canvas.NewRectangle(fillColor)
	rectangle.SetMinSize(fyne.NewSize(width, height))
	rectangle.CornerRadius = 5
	rectangle.StrokeColor = color.Black
	rectangle.StrokeWidth = 0.5
	return rectangle
}

func newTransparentInformationRectangle(width, height float32) *canvas.Rectangle {
	rectangle := canvas.NewRectangle(color.Transparent)
	rectangle.SetMinSize(fyne.NewSize(width, height))
	return rectangle
}
