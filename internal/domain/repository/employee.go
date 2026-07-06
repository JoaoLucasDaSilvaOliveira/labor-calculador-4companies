package repository

import (
	query "labor-calculador-4companies/internal/application/query/employee"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/valueobject"
)

type EmployeeRepository interface {
	Create(firstName string, lastName string, cpf valueobject.CPF) error
	Update(employee *entity.Employee) error
	Delete(sequencialID int) error
	Get(filter query.GetEmployeeWithFilter) ([]*entity.Employee, error)
	GetByID(sequencialID int) (*entity.Employee, error)
}
