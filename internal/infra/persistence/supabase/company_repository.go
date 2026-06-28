package supabase

import (
	query "labor-calculador-4companies/internal/application/query/company"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/repository"
	"labor-calculador-4companies/internal/domain/valueobject"

	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

type companyModel struct {
	ID   int    `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name"`
	CNPJ string `gorm:"column:cnpj"`
}

func (companyModel) TableName() string {
	return "company"
}

func NewCompanyRepository(db *gorm.DB) repository.CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(name string, cnpj valueobject.CNPJ) error {
	model := companyModel{
		Name: name,
		CNPJ: cnpj.String(),
	}

	return r.db.Create(&model).Error
}

func (r *CompanyRepository) Update(company *entity.Company) error {
	model := companyModel{
		ID:   company.SequencialID(),
		Name: company.Name(),
		CNPJ: company.CNPJ(),
	}

	return r.db.Save(&model).Error
}

func (r *CompanyRepository) Delete(sequencialID int) error {
	return r.db.Delete(&companyModel{}, sequencialID).Error
}

func (r *CompanyRepository) Get(filter query.GetCompanyWithFilter) ([]*entity.Company, error) {
	dbQuery := r.db.Model(&companyModel{})

	if filter.IDCompany > 0 {
		dbQuery = dbQuery.Where("id = ?", filter.IDCompany)
	}

	if filter.Name != "" {
		dbQuery = dbQuery.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.CNPJ != "" {
		dbQuery = dbQuery.Where("cnpj = ?", filter.CNPJ.String())
	}

	var models []companyModel
	if err := dbQuery.Find(&models).Error; err != nil {
		return nil, err
	}

	return toCompanyEntities(models)
}

func (r *CompanyRepository) GetByID(sequencialID int) (*entity.Company, error) {
	var model companyModel
	if err := r.db.First(&model, sequencialID).Error; err != nil {
		return nil, err
	}

	return toCompanyEntity(model)
}

func toCompanyEntities(models []companyModel) ([]*entity.Company, error) {
	companies := make([]*entity.Company, 0, len(models))

	for _, model := range models {
		company, err := toCompanyEntity(model)
		if err != nil {
			return nil, err
		}

		companies = append(companies, company)
	}

	return companies, nil
}

func toCompanyEntity(model companyModel) (*entity.Company, error) {
	return entity.NewCompany(model.ID, model.Name, model.CNPJ)
}
