package entity

import (
	"fmt"
	error_factory "labor-calculador-4companies/internal/domain/error"
	"labor-calculador-4companies/internal/domain/valueobject"
	"cloud.google.com/go/civil"
	"github.com/shopspring/decimal"
)

var (
	// CalculateOvertime errors.
	ErrInvalidExtraHoursFactor = error_factory.NewError("fator de hora extra deve ser maior do que 0%")
	ErrInvalidExtraHoursQtt    = error_factory.NewError("quantidade de horas extras não pode ser negativa")
	ErrInvalidExtraHoursIncome = error_factory.NewError("valor da hora não pode ser menor ou igual a zero")

	// CalculateHoursBalanceUsingDays and CalculateHoursBalanceUsingHours errors.
	ErrInvalidBaseSalary  = error_factory.NewError("valor base não pode ser menor ou igual a zero")
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

	// CalculateVacation errors.
	ErrInvalidVacationAcquisitionPeriod      = error_factory.NewError("início do período aquisitivo deve ser anterior ao fim do período aquisitivo")
	ErrInvalidVacationAcquisitionPeriodRange = error_factory.NewError("período aquisitivo não pode ser maior do que 1 ano completo")
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

func CalculateVacation(averageIncomePerMonth decimal.Decimal, startOfAcquisitionPeriod civil.Date, endOfAcquisitionPeriod civil.Date) (*valueobject.VacationAmount, error) {
	if startOfAcquisitionPeriod == endOfAcquisitionPeriod {
		zeroAmount, _ := valueobject.NewVacationAmount(decimal.Zero)
		return zeroAmount, nil
	}

	if !startOfAcquisitionPeriod.Before(endOfAcquisitionPeriod) {
		return nil, fmt.Errorf("%w: %s - %s", ErrInvalidVacationAcquisitionPeriod, startOfAcquisitionPeriod.String(), endOfAcquisitionPeriod.String())
	}

	lastValidEndDate := startOfAcquisitionPeriod.AddYears(1).AddDays(-1)
	if endOfAcquisitionPeriod.After(lastValidEndDate) {
		return nil, fmt.Errorf("%w: %s - %s", ErrInvalidVacationAcquisitionPeriodRange, startOfAcquisitionPeriod.String(), endOfAcquisitionPeriod.String())
	}

	qttMonthsWorked, qttDaysWorked := calculateMonthDayDifference(startOfAcquisitionPeriod, endOfAcquisitionPeriod)

	if qttDaysWorked >= 15 {
		qttMonthsWorked++
	}

	decimalQttMonthsWorked := decimal.NewFromInt(int64(qttMonthsWorked))
	decimalQttMonthsInAYear := decimal.NewFromInt(12)

	vacationAmount, err := valueobject.NewVacationAmount(averageIncomePerMonth.Div(decimalQttMonthsInAYear).Mul(decimalQttMonthsWorked))

	if err != nil {
		return nil, err
	}

	return vacationAmount, nil
}

func CalculateThirteenthSalary(incomePerMonth decimal.Decimal, averages decimal.Decimal, ts *valueobject.ThirteenthSalary) decimal.Decimal {
	civilStartActivities := civil.Date{
		Year:  ts.TargetYear,
		Month: ts.StartOfActivities.Month,
		Day:   ts.StartOfActivities.Day,
	}
	civilEndActivities := civil.Date{
		Year:  ts.TargetYear,
		Month: ts.EndOfActivities.Month,
		Day:   ts.EndOfActivities.Day,
	}

	qttMonthsInAYear := decimal.NewFromInt(12)

	qttMonths := calculateMonthDifferenceInsideTheSameMonth(civilStartActivities, civilEndActivities)

	qttMonthsDecimal := decimal.NewFromInt(int64(qttMonths))

	averageIncome := incomePerMonth.Add(averages)

	return averageIncome.Div(qttMonthsInAYear).Mul(qttMonthsDecimal)
}

func CalculateFGTS(baseValue decimal.Decimal) (decimal.Decimal, error) {
	if baseValue.LessThan(decimal.Zero) {
		return decimal.Zero, ErrInvalidBaseSalary
	}

	return baseValue.Mul(decimal.NewFromFloat(0.08)), nil
}

func calculateMonthDayDifference(startDate civil.Date, endDate civil.Date) (int, int) {
	qttMonths := 0
	for {
		nextMonthDate := startDate.AddMonths(qttMonths + 1)
		if nextMonthDate.After(endDate) {
			break
		}
		qttMonths++
	}

	lastMonthDate := startDate.AddMonths(qttMonths)
	qttDays := endDate.DaysSince(lastMonthDate)

	return qttMonths, qttDays
}

func calculateMonthDifferenceInsideTheSameMonth(startDate civil.Date, endDate civil.Date) int {
	qttMonths := 0
	i := -1
	for {
		//discover the last day of both months
		lastDayStartDate := discoverLastDayInAMonth(startDate.AddMonths(i+1))
		lastDayEndDate := discoverLastDayInAMonth(endDate)
		
		if lastDayStartDate == lastDayEndDate {
			//in this case the dates are in the same month
			if endDate.Day > 14 {
				qttMonths++
			}
			break
		}

		//in the first iteraction, the startDate can be broken, ex: 23/04/yyyy, so we can't assume it is day 1 at first iter
		if i == -1 {
			if (lastDayStartDate.Day - startDate.Day + 1) < 15 {
				i++
				continue
			}
		}

		qttMonths++
		i++
	}	
	return qttMonths
}

func discoverLastDayInAMonth(d civil.Date) civil.Date {
	return civil.Date{Year: d.Year, Month: d.Month, Day: 1}.AddMonths(1).AddDays(-1)
}