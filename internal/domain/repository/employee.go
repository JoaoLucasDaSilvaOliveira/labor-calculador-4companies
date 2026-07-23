package repository

import (
	query "labor-calculador-4companies/internal/application/query/employee"
	"labor-calculador-4companies/internal/domain/entity"
)

type EmployeeRepository interface {
	Create(employee *entity.Employee) error
	Update(employee *entity.Employee) error
	Delete(sequencialID int) error
	GetByCompanyIdWithFilter(filter query.GetEmployeeWithFilter) ([]*entity.Employee, error) //GET
	GetByID(sequencialID int) (*entity.Employee, error)
}
