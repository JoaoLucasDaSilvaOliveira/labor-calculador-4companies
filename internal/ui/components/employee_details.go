package components

import (
	"fmt"
	"image/color"
	"strconv"

	"labor-calculador-4companies/internal/domain/entity"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewEmployeeDetailsComponent preserves the enabled/read-only employee
// information composition from the prototype.
func NewEmployeeDetailsComponent(employee *entity.Employee, companyName string) fyne.CanvasObject {
	cpNameLabel := widget.NewLabelWithStyle(companyName, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	codLabel := widget.NewLabel("CÓDIGO")
	codValue := widget.NewLabel(strconv.Itoa(employee.GetId()))
	codStack := container.NewStack(
		newRoundedInformationRectangle(color.NRGBA{R: 240, G: 240, B: 240, A: 255}, 50, codValue.MinSize().Height),
		codValue,
	)

	nameLabel := widget.NewLabel("NOME")
	nameValue := widget.NewLabel(fmt.Sprintf("%s %s", employee.FirstName(), employee.LastName()))
	nameValue.Truncation = fyne.TextTruncateEllipsis
	nameStack := container.NewStack(
		newRoundedInformationRectangle(color.NRGBA{R: 240, G: 240, B: 240, A: 255}, 300, nameValue.MinSize().Height),
		nameValue,
	)

	cpfLabel := widget.NewLabel("CPF")
	cpfValue := widget.NewLabel(employee.CPF())
	cpfStack := container.NewStack(
		newRoundedInformationRectangle(color.NRGBA{R: 240, G: 240, B: 240, A: 255}, 150, cpfValue.MinSize().Height),
		cpfValue,
	)

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
		cpNameLabel,
		nil,
		employeeLeft,
		employeeRight,
		employeeName,
	)
}
