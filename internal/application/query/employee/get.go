package employee

import "labor-calculador-4companies/internal/domain/valueobject"

type GetEmployeeWithFilter struct {
	CompanyID  int
	FirstName  string
	LastName   string
	CPF        valueobject.CPF
}
