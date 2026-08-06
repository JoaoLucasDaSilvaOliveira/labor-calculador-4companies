package components

import (
	query "labor-calculador-4companies/internal/application/query/employee"
	"labor-calculador-4companies/internal/domain/entity"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

const (
	employeesLoadErrorMessage   = "Não foi possível carregar os funcionários."
	employeesEmptyMessage       = "Nenhum funcionário cadastrada."
	employeesSearchEmptyMessage = "Nenhum funcionário encontrada para essa empresa."
)

// EmployeeFinder is the query contract required by the employee list.
// GetEmployeeUsecase implements this interface without an adapter.
type EmployeeFinder interface {
	Execute(query.GetEmployeeWithFilter) ([]*entity.Employee, error)
}

// NewCompaniesListComponent creates the complete company list.
func NewEmployeesListComponent(finder EmployeeFinder, onEmployeeSelected func(employeeID int)) fyne.CanvasObject {
	return newEmployeesList(
		finder,
		query.GetEmployeeWithFilter{},
		employeesEmptyMessage,
		onEmployeeSelected,
	)
}

func newEmployeesList(finder EmployeeFinder, filter query.GetEmployeeWithFilter, emptyMessage string, onEmployeeSelected func(employeeID int)) fyne.CanvasObject {
	employees, err := finder.Execute(filter)
	if err != nil {
		return newStatusLabel(employeesLoadErrorMessage)
	}
	if len(employees) == 0 {
		return newStatusLabel(emptyMessage)
	}

	employeesList := widget.NewList(
		func() int {
			return len(employees)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			label := object.(*widget.Label)
			label.SetText(employees[id].FirstName())
		},
	)
	employeesList.OnSelected = func(id widget.ListItemID) {
		// The page is kept in router history. Clearing the selection allows the
		// same company to be opened again after the user navigates back.
		employeesList.Unselect(id)
		if onEmployeeSelected != nil {
			onEmployeeSelected(employees[id].GetId())
		}
	}

	return employeesList
}
