package repository

import (
	query "labor-calculador-4companies/internal/application/query/employee"
	"labor-calculador-4companies/internal/domain/entity"
<<<<<<< HEAD
	"labor-calculador-4companies/internal/domain/valueobject"
)

type EmployeeRepository interface {
	Create(firstName string, lastName string, cpf valueobject.CPF) error
=======
)

type EmployeeRepository interface {
	Create(employee *entity.Employee) error
>>>>>>> f2deb6c (feat: now employee is saved on database)
	Update(employee *entity.Employee) error
	Delete(sequencialID int) error
	Get(filter query.GetEmployeeWithFilter) ([]*entity.Employee, error)
	GetByID(sequencialID int) (*entity.Employee, error)
}
