package components

import (
	"image/color"
	"strconv"

	"labor-calculador-4companies/internal/domain/entity"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CompanyDetailsComponent preserves the prototype composition and exposes
// the business fields to the owning page's edit controller.
type CompanyDetailsComponent struct {
	View      fyne.CanvasObject
	NameField *InlineEditableField
	CNPJField *InlineEditableField
	CodeLabel *widget.Label
	fields    []*InlineEditableField
}

func NewCompanyDetailsComponent(company *entity.Company) *CompanyDetailsComponent {
	codLabel := widget.NewLabel("CÓDIGO")
	codStack := container.NewStack(newReadOnlyInformationField(strconv.Itoa(company.GetId()), 50))

	nameLabel := widget.NewLabel("NOME")
	nameField := NewInlineEditableField(company.Name(), 300)
	nameStack := container.NewStack(nameField)

	cnpjLabel := widget.NewLabel("CNPJ")
	cnpjField := NewInlineEditableField(company.CNPJ(), 150)
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

	return &CompanyDetailsComponent{
		View: container.NewBorder(
			widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			nil,
			companyLeft,
			companyRight,
			companyName,
		),
		NameField: nameField,
		CNPJField: cnpjField,
		CodeLabel: codLabel,
		fields:    []*InlineEditableField{nameField, cnpjField},
	}
}

func (details *CompanyDetailsComponent) SetOnActivate(callback func(*InlineEditableField)) {
	for _, field := range details.fields {
		field.SetOnActivate(callback)
	}
}

func (details *CompanyDetailsComponent) SetOnChanged(callback func(*InlineEditableField)) {
	for _, field := range details.fields {
		field.SetOnChanged(callback)
	}
}

func (details *CompanyDetailsComponent) HasChanges() bool {
	for _, field := range details.fields {
		if field.HasChanges() {
			return true
		}
	}
	return false
}

func (details *CompanyDetailsComponent) BeginEdit(target *InlineEditableField) {
	for _, field := range details.fields {
		field.BeginEdit()
	}
	if target != nil {
		target.focus()
	}
}

func (details *CompanyDetailsComponent) CommitEdit() {
	for _, field := range details.fields {
		field.EndEdit()
	}
}

func (details *CompanyDetailsComponent) CancelEdit() {
	for _, field := range details.fields {
		field.CancelEdit()
	}
}

func newReadOnlyInformationField(text string, width float32) fyne.CanvasObject {
	value := widget.NewLabel(text)
	return container.NewStack(
		newRoundedInformationRectangle(color.NRGBA{R: 240, G: 240, B: 240, A: 255}, width, value.MinSize().Height),
		value,
	)
}
