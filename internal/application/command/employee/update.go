package employee

import "labor-calculador-4companies/internal/domain/valueobject"

type UpdateEmployeeCommand struct {
	IDEmployee int
	CompanyID  int
	FirstName  string
	LastName   string
	CPF        valueobject.CPF
}
