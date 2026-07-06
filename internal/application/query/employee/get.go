package employee

import "labor-calculador-4companies/internal/domain/valueobject"

type GetEmployeeWithFilter struct {
	IDEmployee int
	FirstName  string
	LastName   string
	CPF        valueobject.CPF
}
