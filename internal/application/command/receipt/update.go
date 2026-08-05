package receipt

import "labor-calculador-4companies/internal/domain/entity"

type UpdateReceiptCommand struct {
	IDReceipt         int
	EmployeeID        int
	SumaryDescription string
	Items             []entity.ReceiptItem
}
