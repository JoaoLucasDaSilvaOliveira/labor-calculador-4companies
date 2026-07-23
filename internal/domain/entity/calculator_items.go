package entity

import (
	"labor-calculador-4companies/internal/domain/valueobject"

	"cloud.google.com/go/civil"
	"github.com/shopspring/decimal"
)

// SalaryBalanceItem is the proportional salary due in the receipt period.
type SalaryBalanceItem struct {
	amount decimal.Decimal
}

func NewSalaryBalanceItemByDays(baseSalary decimal.Decimal, qttWorkedDays int64) (*SalaryBalanceItem, error) {
	amount, err := CalculateHoursBalanceUsingDays(baseSalary, qttWorkedDays)
	if err != nil {
		return nil, err
	}

	return &SalaryBalanceItem{amount: amount}, nil
}

func NewSalaryBalanceItemByHours(baseSalary decimal.Decimal, qttWorkedHours float32) (*SalaryBalanceItem, error) {
	amount, err := CalculateHoursBalanceUsingHours(baseSalary, qttWorkedHours)
	if err != nil {
		return nil, err
	}

	return &SalaryBalanceItem{amount: amount}, nil
}

func (i SalaryBalanceItem) Value() ReceiptValueAndDescription {
	return ReceiptValueAndDescription{
		Description:  "Saldo de salário",
		DecimalValue: i.amount,
		ValueType:    Provento,
	}
}

// OvertimeItem is the compensation for overtime worked.
type OvertimeItem struct {
	amount decimal.Decimal
}

func NewOvertimeItem(factor float32, qttHours float32, incomePerHour decimal.Decimal) (*OvertimeItem, error) {
	amount, err := CalculateOvertime(factor, qttHours, incomePerHour)
	if err != nil {
		return nil, err
	}

	return &OvertimeItem{amount: amount}, nil
}

func (i OvertimeItem) Value() ReceiptValueAndDescription {
	return ReceiptValueAndDescription{
		Description:  "Horas extras",
		DecimalValue: i.amount,
		ValueType:    Provento,
	}
}

// InsalubrityItem is the unhealthy-work additional compensation.
type InsalubrityItem struct {
	amount decimal.Decimal
}

func NewInsalubrityItem(baseValue decimal.Decimal, factor float32) (*InsalubrityItem, error) {
	amount, err := CalculateInsalubrity(baseValue, factor)
	if err != nil {
		return nil, err
	}

	return &InsalubrityItem{amount: amount}, nil
}

func (i InsalubrityItem) Value() ReceiptValueAndDescription {
	return ReceiptValueAndDescription{
		Description:  "Adicional de insalubridade",
		DecimalValue: i.amount,
		ValueType:    Provento,
	}
}

// NightAdditionalItem is the additional compensation for night work.
type NightAdditionalItem struct {
	amount decimal.Decimal
}

func NewNightAdditionalItem(qttHours float32, incomePerHour decimal.Decimal) (*NightAdditionalItem, error) {
	amount, err := CalculateNightAdditional(qttHours, incomePerHour)
	if err != nil {
		return nil, err
	}

	return &NightAdditionalItem{amount: amount}, nil
}

func (i NightAdditionalItem) Value() ReceiptValueAndDescription {
	return ReceiptValueAndDescription{
		Description:  "Adicional noturno",
		DecimalValue: i.amount,
		ValueType:    Provento,
	}
}

// EffectsItem is the compensation resulting from effects on non-working days.
type EffectsItem struct {
	amount decimal.Decimal
}

func NewEffectsItem(qttNonUtilDays int64, qttUtilDays int64, totalValue decimal.Decimal) (*EffectsItem, error) {
	amount, err := CalculateEffects(qttNonUtilDays, qttUtilDays, totalValue)
	if err != nil {
		return nil, err
	}

	return &EffectsItem{amount: amount}, nil
}

func (i EffectsItem) Value() ReceiptValueAndDescription {
	return ReceiptValueAndDescription{
		Description:  "Reflexos",
		DecimalValue: i.amount,
		ValueType:    Provento,
	}
}

// TransportationVoucherItem is the employee contribution deducted for transportation vouchers.
type TransportationVoucherItem struct {
	amount decimal.Decimal
}

func NewTransportationVoucherItem(baseSalary decimal.Decimal, contributionFactor float32) (*TransportationVoucherItem, error) {
	amount, err := CalculateTransportationVoucher(baseSalary, contributionFactor)
	if err != nil {
		return nil, err
	}

	return &TransportationVoucherItem{amount: amount}, nil
}

func (i TransportationVoucherItem) Value() ReceiptValueAndDescription {
	return ReceiptValueAndDescription{
		Description:  "Desconto de vale-transporte",
		DecimalValue: i.amount,
		ValueType:    Desconto,
	}
}

// VacationItem represents the vacation pay, including the constitutional one-third additional.
type VacationItem struct {
	amount decimal.Decimal
}

func NewVacationItem(averageIncomePerMonth decimal.Decimal, startOfAcquisitionPeriod civil.Date, endOfAcquisitionPeriod civil.Date) (*VacationItem, error) {
	vacationAmount, err := CalculateVacation(averageIncomePerMonth, startOfAcquisitionPeriod, endOfAcquisitionPeriod)
	if err != nil {
		return nil, err
	}

	return &VacationItem{amount: vacationAmount.VacationValue.Add(vacationAmount.OneThird)}, nil
}

func (i VacationItem) Value() ReceiptValueAndDescription {
	return ReceiptValueAndDescription{
		Description:  "Férias com adicional de 1/3",
		DecimalValue: i.amount,
		ValueType:    Provento,
	}
}

// ThirteenthSalaryItem is the proportional thirteenth salary due in the target year.
type ThirteenthSalaryItem struct {
	amount decimal.Decimal
}

func NewThirteenthSalaryItem(incomePerMonth decimal.Decimal, averages decimal.Decimal, thirteenthSalary *valueobject.ThirteenthSalary) *ThirteenthSalaryItem {
	return &ThirteenthSalaryItem{
		amount: CalculateThirteenthSalary(incomePerMonth, averages, thirteenthSalary),
	}
}

func (i ThirteenthSalaryItem) Value() ReceiptValueAndDescription {
	return ReceiptValueAndDescription{
		Description:  "13º salário",
		DecimalValue: i.amount,
		ValueType:    Provento,
	}
}

// FGTSItem is an employer charge displayed as informational in the receipt.
type FGTSItem struct {
	amount decimal.Decimal
}

func NewFGTSItem(baseValue decimal.Decimal) (*FGTSItem, error) {
	amount, err := CalculateFGTS(baseValue)
	if err != nil {
		return nil, err
	}

	return &FGTSItem{amount: amount}, nil
}

func (i FGTSItem) Value() ReceiptValueAndDescription {
	return ReceiptValueAndDescription{
		Description:  "FGTS",
		DecimalValue: i.amount,
		ValueType:    Informativo,
	}
}
