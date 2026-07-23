package receipt

import "labor-calculador-4companies/internal/domain/entity"

type CreateReceiptCommand struct {
	EmployeeID int
	Items      []entity.ReceiptItem
}
