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

// EmployeeDetailsComponent keeps employee presentation separate from the
// company screen while exposing explicit first-name and last-name fields.
type EmployeeDetailsComponent struct {
	View           fyne.CanvasObject
	FirstNameField *InlineEditableField
	LastNameField  *InlineEditableField
	CPFField       *InlineEditableField
	fields         []*InlineEditableField
}

func NewEmployeeDetailsComponent(employee *entity.Employee, companyName string) *EmployeeDetailsComponent {
	cpNameLabel := widget.NewLabelWithStyle(companyName, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	codLabel := widget.NewLabel("CÓDIGO")
	codValue := widget.NewLabel(strconv.Itoa(employee.GetId()))
	codStack := container.NewStack(
		newRoundedInformationRectangle(color.NRGBA{R: 240, G: 240, B: 240, A: 255}, 50, codValue.MinSize().Height),
		codValue,
	)

	firstNameLabel := widget.NewLabel("NOME")
	firstNameField := NewInlineEditableField(employee.FirstName(), 140)
	firstNameStack := container.NewStack(firstNameField)

	lastNameLabel := widget.NewLabel("SOBRENOME")
	lastNameField := NewInlineEditableField(employee.LastName(), 140)
	lastNameStack := container.NewStack(lastNameField)

	cpfLabel := widget.NewLabel("CPF")
	cpfField := NewInlineEditableField(employee.CPF(), 150)
	cpfStack := container.NewStack(cpfField)

	employeeLeft := container.NewHBox(
		codLabel,
		codStack,
		uiUtils.HPadding(10),
	)

	employeeName := container.NewHBox(
		container.NewBorder(nil, nil, firstNameLabel, nil, firstNameStack),
		uiUtils.HPadding(8),
		container.NewBorder(nil, nil, lastNameLabel, nil, lastNameStack),
	)

	employeeRight := container.NewHBox(
		uiUtils.HPadding(10),
		cpfLabel,
		cpfStack,
	)

	return &EmployeeDetailsComponent{
		View: container.NewBorder(
			cpNameLabel,
			nil,
			employeeLeft,
			employeeRight,
			employeeName,
		),
		FirstNameField: firstNameField,
		LastNameField:  lastNameField,
		CPFField:       cpfField,
		fields:         []*InlineEditableField{firstNameField, lastNameField, cpfField},
	}
}

func (details *EmployeeDetailsComponent) SetOnActivate(callback func(*InlineEditableField)) {
	for _, field := range details.fields {
		field.SetOnActivate(callback)
	}
}

func (details *EmployeeDetailsComponent) SetOnChanged(callback func(*InlineEditableField)) {
	for _, field := range details.fields {
		field.SetOnChanged(callback)
	}
}

func (details *EmployeeDetailsComponent) HasChanges() bool {
	for _, field := range details.fields {
		if field.HasChanges() {
			return true
		}
	}
	return false
}

func (details *EmployeeDetailsComponent) BeginEdit(target *InlineEditableField) {
	for _, field := range details.fields {
		field.BeginEdit()
	}
	if target != nil {
		target.focus()
	}
}

func (details *EmployeeDetailsComponent) CommitEdit() {
	for _, field := range details.fields {
		field.EndEdit()
	}
}

func (details *EmployeeDetailsComponent) CancelEdit() {
	for _, field := range details.fields {
		field.CancelEdit()
	}
}
