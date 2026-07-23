package entity_test

import (
	"errors"
	"fmt"
	"testing"

	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/valueobject"

	"cloud.google.com/go/civil"
	"github.com/shopspring/decimal"
)

func TestCalculateHoursBalanceUsingDays(t *testing.T) {
	t.Run("calculates proportional salary by worked days", func(t *testing.T) {
		result, err := entity.CalculateHoursBalanceUsingDays(decimal.NewFromInt(3000), 15)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		assertDecimalEqual(t, result, "1500")
	})

	t.Run("returns error when worked days is invalid", func(t *testing.T) {
		_, err := entity.CalculateHoursBalanceUsingDays(decimal.NewFromInt(3000), 31)
		assertErrorIs(t, err, entity.ErrInvalidWorkedDays)
	})
}

func TestCalculateHoursBalanceUsingHours(t *testing.T) {
	t.Run("calculates proportional salary by worked hours", func(t *testing.T) {
		result, err := entity.CalculateHoursBalanceUsingHours(decimal.NewFromInt(2200), 110)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		assertDecimalEqual(t, result, "1100")
	})

	t.Run("returns error when base salary is invalid", func(t *testing.T) {
		_, err := entity.CalculateHoursBalanceUsingHours(decimal.Zero, 110)
		assertErrorIs(t, err, entity.ErrInvalidBaseSalary)
	})
}

func TestCalculateOvertime(t *testing.T) {
	t.Run("calculates overtime", func(t *testing.T) {
		result, err := entity.CalculateOvertime(1.5, 2, decimal.NewFromInt(10))
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		assertDecimalEqual(t, result, "30")
	})

	t.Run("returns error when factor is invalid", func(t *testing.T) {
		_, err := entity.CalculateOvertime(0.5, 2, decimal.NewFromInt(10))
		assertErrorIs(t, err, entity.ErrInvalidExtraHoursFactor)
	})

	t.Run("here", func(t *testing.T) {
		value, _ := entity.CalculateOvertime(2, 2, decimal.NewFromInt(10))
		fmt.Println(value)
	})
}

func TestCalculateInsalubrity(t *testing.T) {
	t.Run("calculates insalubrity", func(t *testing.T) {
		result, err := entity.CalculateInsalubrity(decimal.NewFromInt(1000), 0.2)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		assertDecimalEqual(t, result, "200")
	})

	t.Run("returns error when factor is invalid", func(t *testing.T) {
		_, err := entity.CalculateInsalubrity(decimal.NewFromInt(1000), -0.1)
		assertErrorIs(t, err, entity.ErrInvalidInsalubrityFactor)
	})
}

func TestCalculateNightAdditional(t *testing.T) {
	t.Run("calculates night additional", func(t *testing.T) {
		result, err := entity.CalculateNightAdditional(5, decimal.NewFromInt(10))
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		assertDecimalEqual(t, result, "10")
	})

	t.Run("returns error when hours is invalid", func(t *testing.T) {
		_, err := entity.CalculateNightAdditional(-1, decimal.NewFromInt(10))
		assertErrorIs(t, err, entity.ErrInvalidNightAdditionalHours)
	})

	t.Run("here", func(t *testing.T) {
		value, _ := entity.CalculateNightAdditional(73.68, decimal.NewFromFloat(10.11))
		fmt.Println(value)
	})
}

func TestCalculateEffects(t *testing.T) {
	t.Run("calculates effects", func(t *testing.T) {
		result, err := entity.CalculateEffects(4, 20, decimal.NewFromInt(1000))
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		assertDecimalEqual(t, result, "200")
	})

	t.Run("returns error when util days is invalid", func(t *testing.T) {
		_, err := entity.CalculateEffects(4, 0, decimal.NewFromInt(1000))
		assertErrorIs(t, err, entity.ErrInvalidEffectsUtilDays)
	})
}

func TestCalculateTransportationVoucher(t *testing.T) {
	t.Run("calculates transportation voucher", func(t *testing.T) {
		result, err := entity.CalculateTransportationVoucher(decimal.NewFromInt(1000), 0.06)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		assertDecimalEqual(t, result, "60")
	})

	t.Run("returns error when contribution factor is invalid", func(t *testing.T) {
		_, err := entity.CalculateTransportationVoucher(decimal.NewFromInt(1000), 0.07)
		assertErrorIs(t, err, entity.ErrInvalidTransportationContribution)
	})
}

func TestCalculateVacation(t *testing.T) {
	t.Run("accepts acquisition period with one complete year", func(t *testing.T) {
		_, err := entity.CalculateVacation(
			decimal.NewFromInt(3000),
			civil.Date{Year: 2025, Month: 2, Day: 23},
			civil.Date{Year: 2026, Month: 2, Day: 22},
		)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("returns error when start is not before end", func(t *testing.T) {
		_, err := entity.CalculateVacation(
			decimal.NewFromInt(3000),
			civil.Date{Year: 2025, Month: 2, Day: 23},
			civil.Date{Year: 2025, Month: 2, Day: 23},
		)
		assertErrorIs(t, err, entity.ErrInvalidVacationAcquisitionPeriod)
	})

	t.Run("returns error when acquisition period is greater than one complete year", func(t *testing.T) {
		_, err := entity.CalculateVacation(
			decimal.NewFromInt(3000),
			civil.Date{Year: 2025, Month: 2, Day: 23},
			civil.Date{Year: 2026, Month: 2, Day: 23},
		)
		assertErrorIs(t, err, entity.ErrInvalidVacationAcquisitionPeriodRange)
	})

	t.Run("shows on log the diference of months and days between two acquisitive period dates", func(t *testing.T) {
		_, _ = entity.CalculateVacation(
			decimal.NewFromInt(3000),
			civil.Date{Year: 2025, Month: 2, Day: 23},
			civil.Date{Year: 2026, Month: 2, Day: 22},
		)
	})

	t.Run("returns the amount of vacation", func(t *testing.T) {
		amount, err := entity.CalculateVacation(
			decimal.NewFromInt(2000),
			civil.Date{Year: 2025, Month: 2, Day: 23},
			civil.Date{Year: 2025, Month: 3, Day: 10},
		)

		if err != nil {
			fmt.Println(err) 
			return
		}
		
		fmt.Println(amount.VacationValue.Round(2))
	})
}

func TestCalculateThirteenthSalary(t *testing.T) {
	t.Run("here", func(t *testing.T) {
		ts, _ := valueobject.NewThirteenthSalary(2026, 1, 18, 12, 31)
		fmt.Println(entity.CalculateThirteenthSalary(decimal.NewFromFloat32(1500.00), decimal.Zero, ts))
	})
}

func TestCalculateFGTS(t *testing.T) {
	t.Run("here", func(t *testing.T) {
		value, err := entity.CalculateFGTS(decimal.NewFromFloat32(-1))

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)
	})
}

func assertDecimalEqual(t *testing.T, got decimal.Decimal, want string) {
	t.Helper()

	expected := decimal.RequireFromString(want)
	if !got.Equal(expected) {
		t.Fatalf("expected %s, got %s", expected.String(), got.String())
	}
}

func assertErrorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("expected error %v, got %v", target, err)
	}
}
