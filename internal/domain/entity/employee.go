package entity

import (
	"fmt"
	error_factory "labor-calculador-4companies/internal/domain/error"
	"labor-calculador-4companies/internal/domain/valueobject"
	"strings"
)

var (
	ErrInvalidEmployeeFirstName = error_factory.NewError("primeiro nome do funcionário inválido")
	ErrInvalidEmployeeLastName  = error_factory.NewError("sobrenome do funcionário inválido")
)

// This entity wont be saved on bd, it's gonna be used on making receipt process
type Employee struct {
	firstName string
	lastName  string
	cpf       valueobject.CPF
}

func NewEmployee(firstName string, lastName string, cpf string) (*Employee, error) {
	employee := &Employee{}

	if err := employee.SetFirstName(firstName); err != nil {
		return nil, err
	}

	if err := employee.SetLastName(lastName); err != nil {
		return nil, err
	}

	if err := employee.SetCPF(cpf); err != nil {
		return nil, err
	}

	return employee, nil
}

func (e *Employee) FirstName() string {
	return e.firstName
}

func (e *Employee) SetFirstName(firstName string) error {
	firstName = strings.TrimSpace(firstName)

	if firstName == "" {
		return fmt.Errorf("%w: %s", ErrInvalidEmployeeFirstName, firstName)
	}

	e.firstName = firstName
	return nil
}

func (e *Employee) LastName() string {
	return e.lastName
}

func (e *Employee) SetLastName(lastName string) error {
	lastName = strings.TrimSpace(lastName)

	if lastName == "" {
		return fmt.Errorf("%w: %s", ErrInvalidEmployeeLastName, lastName)
	}

	e.lastName = lastName
	return nil
}

func (e *Employee) CPF() string {
	return e.cpf.String()
}

func (e *Employee) SetCPF(cpf string) error {
	validCPF, err := valueobject.NewCPF(cpf)
	if err != nil {
		return err
	}

	e.cpf = validCPF
	return nil
}
