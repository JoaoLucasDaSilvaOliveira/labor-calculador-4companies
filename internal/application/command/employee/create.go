package employee

import "labor-calculador-4companies/internal/domain/valueobject"

type CreateEmployeeCommand struct {
	FirstName string
	LastName  string
	CPF       valueobject.CPF
}
