package employee

import "labor-calculador-4companies/internal/domain/valueobject"

type CreateEmployeeCommand struct {
	CompanyID int
	FirstName string
	LastName  string
	CPF       valueobject.CPF
}
