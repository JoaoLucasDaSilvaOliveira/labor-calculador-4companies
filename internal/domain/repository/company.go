package repository

import (
	query "labor-calculador-4companies/internal/application/query/company"
	"labor-calculador-4companies/internal/domain/entity"
)

type CompanyRepository interface {
	Create(company *entity.Company) error
	Update(company *entity.Company) error
	Delete(sequencialID int) error
	Get(filter query.GetCompanyWithFilter) ([]*entity.Company, error)
	GetByID(sequencialID int) (*entity.Company, error)
}
