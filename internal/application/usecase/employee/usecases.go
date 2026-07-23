package employee

import (
	command "labor-calculador-4companies/internal/application/command/employee"
	query "labor-calculador-4companies/internal/application/query/employee"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/repository"
)

type CreateEmployeeUsecase struct {
	repository repository.EmployeeRepository
}

func NewCreateEmployeeUsecase(repository repository.EmployeeRepository) *CreateEmployeeUsecase {
	return &CreateEmployeeUsecase{repository: repository}
}
func (uc *CreateEmployeeUsecase) Execute(cmd command.CreateEmployeeCommand) error {
	employee, err := entity.NewEmployee(cmd.CompanyID, cmd.FirstName, cmd.LastName, cmd.CPF.String())
	if err != nil {
		return err
	}

	return uc.repository.Create(employee)
}

//---------------------------------------------------------------------------------

type DeleteEmployeeUsecase struct {
	repository repository.EmployeeRepository
}

func NewDeleteEmployeeUsecase(repository repository.EmployeeRepository) *DeleteEmployeeUsecase {
	return &DeleteEmployeeUsecase{repository: repository}
}
func (uc *DeleteEmployeeUsecase) Execute(cmd command.DeleteEmployeeCommand) error {
	return uc.repository.Delete(cmd.IDEmployee)
}

//---------------------------------------------------------------------------------

type UpdateEmployeeUsecase struct {
	repository repository.EmployeeRepository
}

func NewUpdateEmployeeUsecase(repository repository.EmployeeRepository) *UpdateEmployeeUsecase {
	return &UpdateEmployeeUsecase{repository: repository}
}
func (uc *UpdateEmployeeUsecase) Execute(cmd command.UpdateEmployeeCommand) error {
	employee, err := entity.LoadEmployee(cmd.IDEmployee, cmd.CompanyID, cmd.FirstName, cmd.LastName, cmd.CPF.String())
	if err != nil {
		return err
	}

	return uc.repository.Update(employee)
}

//---------------------------------------------------------------------------------

type GetEmployeeUsecase struct {
	repository repository.EmployeeRepository
}

func NewGetEmployeeUsecase(repository repository.EmployeeRepository) *GetEmployeeUsecase {
	return &GetEmployeeUsecase{repository: repository}
}
func (uc *GetEmployeeUsecase) Execute(qry query.GetEmployeeWithFilter) ([]*entity.Employee, error) {
	return uc.repository.GetByCompanyIdWithFilter(qry)
}

//---------------------------------------------------------------------------------

type GetEmployeeByIdUsecase struct {
	repository repository.EmployeeRepository
}

func NewGetEmployeeByIdUsecase(repository repository.EmployeeRepository) *GetEmployeeByIdUsecase {
	return &GetEmployeeByIdUsecase{repository: repository}
}
func (uc *GetEmployeeByIdUsecase) Execute(cmd command.GetEmployeeById) (*entity.Employee, error) {
	return uc.repository.GetByID(cmd.IDEmployee)
}
