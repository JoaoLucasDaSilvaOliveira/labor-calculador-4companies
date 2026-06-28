package entity_test

import (
	"errors"
	"testing"

	"labor-calculador-4companies/internal/domain/entity"

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
