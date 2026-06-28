package company

import (
	command "labor-calculador-4companies/internal/application/command/company"
	query "labor-calculador-4companies/internal/application/query/company"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/repository"
)

type CreateCompanyUsecase struct {
	repository repository.CompanyRepository
}
func NewCreateCompanyUsecase(repository repository.CompanyRepository) *CreateCompanyUsecase {
	return &CreateCompanyUsecase{repository: repository}
}
func (uc *CreateCompanyUsecase) Execute(cmd command.CreateCompanyCommand) error {
	return uc.repository.Create(cmd.Name, cmd.CNPJ)
}

//---------------------------------------------------------------------------------

type DeleteCompanyUsecase struct {
	repository repository.CompanyRepository
}
func NewDeleteCompanyUsecase(repository repository.CompanyRepository) *DeleteCompanyUsecase {
	return &DeleteCompanyUsecase{repository: repository}
}
func (uc *DeleteCompanyUsecase) Execute(cmd command.DeleteCompanyCommand) error {
	return uc.repository.Delete(cmd.IDCompany)
}

//---------------------------------------------------------------------------------

type UpdateCompanyUsecase struct {
	repository repository.CompanyRepository
}
func NewUpdateCompanyUsecase(repository repository.CompanyRepository) *UpdateCompanyUsecase {
	return &UpdateCompanyUsecase{repository: repository}
}
func (uc *UpdateCompanyUsecase) Execute(cmd command.UpdateCompanyCommand) error {
	company, err := entity.NewCompany(cmd.IDCompany, cmd.Name, cmd.CNPJ.String())
	if err != nil {
		return err
	}

	return uc.repository.Update(company)
}

//---------------------------------------------------------------------------------

type GetCompanyUsecase struct {
	repository repository.CompanyRepository
}
func NewGetCompanyUsecase(repository repository.CompanyRepository) *GetCompanyUsecase {
	return &GetCompanyUsecase{repository: repository}
}
func (uc *GetCompanyUsecase) Execute(qry query.GetCompanyWithFilter) ([]*entity.Company, error) {
	return uc.repository.Get(qry)
}

//---------------------------------------------------------------------------------

type GetCompanyByIdUsecase struct {
	repository repository.CompanyRepository
}
func NewGetCompanyByIdUsecase(repository repository.CompanyRepository) *GetCompanyByIdUsecase {
	return &GetCompanyByIdUsecase{repository: repository}
}
func (uc *GetCompanyByIdUsecase) Execute(cmd command.GetCompanyById) (*entity.Company, error) {
	return uc.repository.GetByID(cmd.IDCompany)
}
