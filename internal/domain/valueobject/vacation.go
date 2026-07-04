package valueobject

import (
	error_factory "labor-calculador-4companies/internal/domain/error"

	"github.com/shopspring/decimal"
)

var (
	ErrNegativeVacationAmount = error_factory.NewError("valor de férias não pode ser negativo")
)

type VacationAmount struct {
	VacationValue decimal.Decimal
	OneThird      decimal.Decimal
}

func NewVacationAmount(amount decimal.Decimal) (*VacationAmount, error) {
	if amount.LessThan(decimal.Zero) {
		return nil, ErrNegativeVacationAmount
	}

	//rounding values
	
	
	return &VacationAmount{
		VacationValue: amount,
		OneThird: amount.Div(decimal.NewFromInt(3)),
	}, nil
}


