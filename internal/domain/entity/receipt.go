package entity

import (
	"fmt"
	error_factory "labor-calculador-4companies/internal/domain/error"
	"strings"

	"github.com/shopspring/decimal"
)

// errors declarations
var (
	ErrInvalidReceiptSequencialID = error_factory.NewError("id sequencial do recibo inválido")
)

// enum types declarations
type ReceiptValueType int

const (
	Provento ReceiptValueType = iota
	Desconto
	Informativo
)

// compositions objects declarations
type ReceiptValueAndDescription struct {
	Description  string
	DecimalValue decimal.Decimal
	ValueType    ReceiptValueType
}

type ReceiptItem interface {
	Value() ReceiptValueAndDescription
}

// main object declaration
type Receipt struct {
	id                int
	idEmployee        int
	sumaryDescription string
	Items             []ReceiptItem
}

func NewReceipt(idEmployee int, sumaryDescription string, items []ReceiptItem) *Receipt {
	receipt := new(Receipt)
	receipt.idEmployee = idEmployee
	receipt.sumaryDescription = strings.TrimSpace(sumaryDescription)

	for _, item := range items {
		receipt.AddInformationItems(item)
	}

	return receipt
}

func LoadReceipt(id int, idEmployee int, sumaryDescription string, items []ReceiptItem) (*Receipt, error) {
	receipt := NewReceipt(idEmployee, sumaryDescription, items)

	if err := receipt.SetId(id); err != nil {
		return nil, err
	}

	return receipt, nil
}

func (r *Receipt) GetId() int {
	return r.id
}

func (r *Receipt) SetId(id int) error {
	if id <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidReceiptSequencialID, id)
	}

	r.id = id
	return nil
}

func (r *Receipt) GetEmployeeId() int {
	return r.idEmployee
}

func (r *Receipt) SetEmployeeId(id int) {
	r.idEmployee = id
}

func (r *Receipt) GetSumaryDescription() string {
	return r.sumaryDescription
}

func (r *Receipt) SetSumaryDescription(sumaryDescription string) {
	r.sumaryDescription = strings.TrimSpace(sumaryDescription)
}

func (r *Receipt) AddInformationItems(item ReceiptItem) {
	r.Items = append(r.Items, item)
}

// use this when need to display informations: page/pdf/preview
func (r *Receipt) GetInformationsFromItems() []ReceiptValueAndDescription {
	receiptValueAndDescription := make([]ReceiptValueAndDescription, 0)
	for _, item := range r.Items {
		receiptValueAndDescription = append(receiptValueAndDescription, item.Value())
	}
	return receiptValueAndDescription
}
