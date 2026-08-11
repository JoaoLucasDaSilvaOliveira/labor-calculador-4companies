package repository

import (
	query "labor-calculador-4companies/internal/application/query/receipt"
	"labor-calculador-4companies/internal/domain/entity"
)

type ReceiptRepository interface {
	Create(receipt *entity.Receipt) error
	Update(receipt *entity.Receipt) error
	Delete(sequencialID int) error
	GetByEmployeeIdWithFilter(filter query.GetReceiptWithFilter) ([]*entity.Receipt, error)
	GetByID(sequencialID int) (*entity.Receipt, error)
}
