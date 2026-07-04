package valueobject

import (
	"fmt"
	error_factory "labor-calculador-4companies/internal/domain/error"
	"time"

	"cloud.google.com/go/civil"
)

var (
	ErrInvalidThirteenthSalaryYear        = error_factory.NewError("ano do décimo terceiro deve estar entre 1900 e 9999")
	ErrInvalidThirteenthSalaryDate        = error_factory.NewError("data do décimo terceiro inválida")
	ErrInvalidThirteenthSalaryPeriodRange = error_factory.NewError("início das atividades deve ser anterior ou igual ao fim das atividades")
)

type MonthDay struct {
	Month time.Month
	Day   int
}

func (md MonthDay) isValid(year int) bool {
	// create and civil.Date in the supplied year to check if it is a valid date
	t := civil.Date{
		Year:  year,
		Month: md.Month,
		Day:   md.Day,
	}
	return t.IsValid()
}

func (md MonthDay) isAfter(other MonthDay) bool {
	if md.Month != other.Month {
		return md.Month > other.Month
	}

	return md.Day > other.Day
}

func newMonthDay(year int, month time.Month, day int) (*MonthDay, error) {
	md := &MonthDay{
		Month: month,
		Day:   day,
	}

	if !md.isValid(year) {
		return nil, fmt.Errorf("%w: %04d-%02d-%02d", ErrInvalidThirteenthSalaryDate, year, month, day)
	}

	return md, nil
}

type ThirteenthSalary struct {
	TargetYear        int
	StartOfActivities *MonthDay
	EndOfActivities   *MonthDay
}

func NewThirteenthSalary(year int, monthStartActivities time.Month, dayStartActivities int, montheEdOfActivities time.Month, dayEndActivities int) (*ThirteenthSalary, error) {
	//check if date is valid
	// let's consider 1900 the first valid year and 9999 the last valid year
	if year < 1900 || year > 9999 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidThirteenthSalaryYear, year)
	}

	startOfActivities, err := newMonthDay(year, monthStartActivities, dayStartActivities)

	if err != nil {
		return nil, err
	}

	endOfActivities, err := newMonthDay(year, montheEdOfActivities, dayEndActivities)

	if err != nil {
		return nil, err
	}

	if startOfActivities.isAfter(*endOfActivities) {
		return nil, fmt.Errorf("%w: %02d-%02d - %02d-%02d", ErrInvalidThirteenthSalaryPeriodRange, startOfActivities.Month, startOfActivities.Day, endOfActivities.Month, endOfActivities.Day)
	}

	return &ThirteenthSalary{
		TargetYear:        year,
		StartOfActivities: startOfActivities,
		EndOfActivities:   endOfActivities,
	}, nil
}
