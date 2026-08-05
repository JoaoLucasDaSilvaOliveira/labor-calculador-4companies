package receipt

import (
	command "labor-calculador-4companies/internal/application/command/receipt"
	query "labor-calculador-4companies/internal/application/query/receipt"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/repository"
)

type CreateReceiptUsecase struct {
	repository repository.ReceiptRepository
}

func NewCreateReceiptUsecase(repository repository.ReceiptRepository) *CreateReceiptUsecase {
	return &CreateReceiptUsecase{repository: repository}
}

func (uc *CreateReceiptUsecase) Execute(cmd command.CreateReceiptCommand) error {
	receipt := entity.NewReceipt(cmd.EmployeeID, cmd.SumaryDescription, cmd.Items)

	return uc.repository.Create(receipt)
}

//---------------------------------------------------------------------------------

type DeleteReceiptUsecase struct {
	repository repository.ReceiptRepository
}

func NewDeleteReceiptUsecase(repository repository.ReceiptRepository) *DeleteReceiptUsecase {
	return &DeleteReceiptUsecase{repository: repository}
}

func (uc *DeleteReceiptUsecase) Execute(cmd command.DeleteReceiptCommand) error {
	return uc.repository.Delete(cmd.IDReceipt)
}

//---------------------------------------------------------------------------------

type UpdateReceiptUsecase struct {
	repository repository.ReceiptRepository
}

func NewUpdateReceiptUsecase(repository repository.ReceiptRepository) *UpdateReceiptUsecase {
	return &UpdateReceiptUsecase{repository: repository}
}

func (uc *UpdateReceiptUsecase) Execute(cmd command.UpdateReceiptCommand) error {
	receipt, err := entity.LoadReceipt(cmd.IDReceipt, cmd.EmployeeID, cmd.SumaryDescription, cmd.Items)
	if err != nil {
		return err
	}

	return uc.repository.Update(receipt)
}

//---------------------------------------------------------------------------------

type GetReceiptUsecase struct {
	repository repository.ReceiptRepository
}

func NewGetReceiptUsecase(repository repository.ReceiptRepository) *GetReceiptUsecase {
	return &GetReceiptUsecase{repository: repository}
}

func (uc *GetReceiptUsecase) Execute(qry query.GetReceiptWithFilter) ([]*entity.Receipt, error) {
	return uc.repository.GetByEmployeeIdWithFilter(qry)
}

//---------------------------------------------------------------------------------

type GetReceiptByIdUsecase struct {
	repository repository.ReceiptRepository
}

func NewGetReceiptByIdUsecase(repository repository.ReceiptRepository) *GetReceiptByIdUsecase {
	return &GetReceiptByIdUsecase{repository: repository}
}

func (uc *GetReceiptByIdUsecase) Execute(cmd command.GetReceiptById) (*entity.Receipt, error) {
	return uc.repository.GetByID(cmd.IDReceipt)
}
