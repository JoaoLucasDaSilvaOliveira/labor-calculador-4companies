package entity

import (
	"fmt"
	error_factory "labor-calculador-4companies/internal/domain/error"

	"github.com/shopspring/decimal"
)

var (
	// CalculateOvertime errors.
	ErrInvalidExtraHoursFactor = error_factory.NewError("fator de hora extra deve ser maior do que 0%")
	ErrInvalidExtraHoursQtt    = error_factory.NewError("quantidade de horas extras não pode ser negativa")
	ErrInvalidExtraHoursIncome = error_factory.NewError("valor da hora não pode ser menor ou igual a zero")

	// CalculateHoursBalanceUsingDays and CalculateHoursBalanceUsingHours errors.
	ErrInvalidBaseSalary  = error_factory.NewError("salário base não pode ser menor ou igual a zero")
	ErrInvalidWorkedDays  = error_factory.NewError("quantidade de dias trabalhados deve estar entre 0 e 30")
	ErrInvalidWorkedHours = error_factory.NewError("quantidade de horas trabalhadas deve estar entre 0 e 220")

	// CalculateInsalubrity errors.
	ErrInvalidInsalubrityFactor    = error_factory.NewError("fator de insalubridade não pode ser negativo")
	ErrInvalidInsalubrityBaseValue = error_factory.NewError("valor base de insalubridade não pode ser menor ou igual a zero")

	// CalculateNightAdditional errors.
	ErrInvalidNightAdditionalHours = error_factory.NewError("quantidade de horas noturnas não pode ser negativa")

	// CalculateEffects errors.
	ErrInvalidEffectsNonUtilDays = error_factory.NewError("quantidade de dias não úteis não pode ser negativa")
	ErrInvalidEffectsUtilDays    = error_factory.NewError("quantidade de dias úteis deve ser maior do que zero")
	ErrInvalidEffectsTotalValue  = error_factory.NewError("valor total dos reflexos não pode ser negativo")

	// CalculateTransportationVoucher errors.
	ErrInvalidTransportationContribution = error_factory.NewError("fator de contribuição do vale-transporte deve estar entre 0% e 6%")
	ErrInvalidTransportationVoucherBase  = error_factory.NewError("salário base do vale-transporte não pode ser menor ou igual a zero")
)

func CalculateHoursBalanceUsingDays(baseSalary decimal.Decimal, qttWorkedDays int64) (decimal.Decimal, error) {
	if baseSalary.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, fmt.Errorf("%w: %s", ErrInvalidBaseSalary, baseSalary.String())
	}

	if qttWorkedDays < 0 || qttWorkedDays > 30 {
		return decimal.Zero, fmt.Errorf("%w: %d", ErrInvalidWorkedDays, qttWorkedDays)
	}

	// we'r gonna normalize the quantity of days in a month to 30 max days
	qttDaysInAMonth := decimal.NewFromInt(30)

	decimalQttWorkedDays := decimal.NewFromInt(qttWorkedDays)

	return baseSalary.Div(qttDaysInAMonth).Mul(decimalQttWorkedDays), nil
}

func CalculateHoursBalanceUsingHours(baseSalary decimal.Decimal, qttWorkedHours float32) (decimal.Decimal, error) {
	if baseSalary.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, fmt.Errorf("%w: %s", ErrInvalidBaseSalary, baseSalary.String())
	}

	if qttWorkedHours < 0 || qttWorkedHours > 220 {
		return decimal.Zero, fmt.Errorf("%w: %f", ErrInvalidWorkedHours, qttWorkedHours)
	}

	maxHoursOfWorkInAMonth := decimal.NewFromInt(220)

	decimalQttWorkedHours := decimal.NewFromFloat32(qttWorkedHours)

	return baseSalary.Div(maxHoursOfWorkInAMonth).Mul(decimalQttWorkedHours), nil
}

func CalculateOvertime(factor float32, qttHours float32, incomePerHour decimal.Decimal) (decimal.Decimal, error) {
	if factor < 1 {
		return decimal.Zero, fmt.Errorf("%w: %f", ErrInvalidExtraHoursFactor, factor)
	}

	if qttHours < 0 {
		return decimal.Zero, fmt.Errorf("%w: %f", ErrInvalidExtraHoursQtt, qttHours)
	}

	if incomePerHour.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, fmt.Errorf("%w: %s", ErrInvalidExtraHoursIncome, incomePerHour.String())
	}

	decimalFactor := decimal.NewFromFloat32(factor)
	decimalQttHours := decimal.NewFromFloat32(qttHours)
	return incomePerHour.Mul(decimalFactor).Mul(decimalQttHours), nil
}

func CalculateInsalubrity(baseValue decimal.Decimal, factor float32) (decimal.Decimal, error) {
	if factor < 0 {
		return decimal.Zero, fmt.Errorf("%w: %f", ErrInvalidInsalubrityFactor, factor)
	}

	if baseValue.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, fmt.Errorf("%w: %s", ErrInvalidInsalubrityBaseValue, baseValue.String())
	}

	decimalFactor := decimal.NewFromFloat32(factor)
	return baseValue.Mul(decimalFactor), nil
}

func CalculateNightAdditional(qttHours float32, incomePerHour decimal.Decimal) (decimal.Decimal, error) {
	if qttHours < 0 {
		return decimal.Zero, fmt.Errorf("%w: %f", ErrInvalidNightAdditionalHours, qttHours)
	}

	if incomePerHour.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, fmt.Errorf("%w: %s", ErrInvalidExtraHoursIncome, incomePerHour.String())
	}

	decimalQttHours := decimal.NewFromFloat32(qttHours)
	factor := decimal.NewFromFloat(0.2)
	return decimalQttHours.Mul(incomePerHour).Mul(factor), nil
}

func CalculateEffects(qttNonUtilDays int64, qttUtilDays int64, totalValue decimal.Decimal) (decimal.Decimal, error) {
	if qttNonUtilDays < 0 {
		return decimal.Zero, fmt.Errorf("%w: %d", ErrInvalidEffectsNonUtilDays, qttNonUtilDays)
	}

	if qttUtilDays <= 0 {
		return decimal.Zero, fmt.Errorf("%w: %d", ErrInvalidEffectsUtilDays, qttUtilDays)
	}

	if totalValue.IsNegative() {
		return decimal.Zero, fmt.Errorf("%w: %s", ErrInvalidEffectsTotalValue, totalValue.String())
	}

	decimalQttNonUtilDays := decimal.NewFromInt(qttNonUtilDays)
	decimalQttUtilDays := decimal.NewFromInt(qttUtilDays)

	totalValueDividedByQttUtilDays := totalValue.Div(decimalQttUtilDays)

	return totalValueDividedByQttUtilDays.Mul(decimalQttNonUtilDays), nil
}

func CalculateTransportationVoucher(baseSalary decimal.Decimal, contributionFactor float32) (decimal.Decimal, error) {
	if baseSalary.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, fmt.Errorf("%w: %s", ErrInvalidTransportationVoucherBase, baseSalary.String())
	}

	if contributionFactor < 0 || contributionFactor > 0.06 {
		return decimal.Zero, fmt.Errorf("%w: %f", ErrInvalidTransportationContribution, contributionFactor)
	}

	decimalContributionFator := decimal.NewFromFloat32(contributionFactor)

	return baseSalary.Mul(decimalContributionFator), nil
}
