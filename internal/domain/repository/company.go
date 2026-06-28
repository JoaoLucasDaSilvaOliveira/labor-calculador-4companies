package repository

import (
	query "labor-calculador-4companies/internal/application/query/company"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/valueobject"
)

type CompanyRepository interface {
	Create(name string, cnpj valueobject.CNPJ) error
	Update(company *entity.Company) error
	Delete(sequencialID int) error
	Get(filter query.GetCompanyWithFilter) ([]*entity.Company, error)
	GetByID(sequencialID int) (*entity.Company, error)
}
