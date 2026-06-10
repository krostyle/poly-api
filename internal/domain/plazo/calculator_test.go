package plazo_test

import (
	"testing"
	"time"

	"poly.app/api/internal/domain/plazo"
)

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func TestCalculateDeadline_NoHolidays(t *testing.T) {
	// Monday 2 Jun 2025 + 5 business days = Monday 9 Jun 2025
	start := date(2025, 6, 2)
	got := plazo.CalculateDeadline(start, 5, nil)
	want := date(2025, 6, 9)
	if !got.Equal(want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestCalculateDeadline_SkipsWeekend(t *testing.T) {
	// Thursday 5 Jun + 3 business days = Tuesday 10 Jun (skips Sat+Sun)
	start := date(2025, 6, 5)
	got := plazo.CalculateDeadline(start, 3, nil)
	want := date(2025, 6, 10)
	if !got.Equal(want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestCalculateDeadline_WithHoliday(t *testing.T) {
	// Wed 4 Jun + 3 business days; Thu 5 is a holiday
	// Thu 5 (holiday, skip) → Fri 6 (day 1) → Mon 9 (day 2) → Tue 10 (day 3)
	start := date(2025, 6, 4)
	holidays := []time.Time{date(2025, 6, 5)}
	got := plazo.CalculateDeadline(start, 3, holidays)
	want := date(2025, 6, 10)
	if !got.Equal(want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestRemainingBusinessDays_Positive(t *testing.T) {
	today := date(2025, 6, 2)    // Monday
	deadline := date(2025, 6, 6) // Friday = 4 business days
	got := plazo.RemainingBusinessDays(today, deadline, nil)
	if got != 4 {
		t.Errorf("want 4, got %d", got)
	}
}

func TestRemainingBusinessDays_Overdue(t *testing.T) {
	today := date(2025, 6, 10)
	deadline := date(2025, 6, 6)
	got := plazo.RemainingBusinessDays(today, deadline, nil)
	if got >= 0 {
		t.Errorf("want negative, got %d", got)
	}
}

func TestEvaluateSemaforo(t *testing.T) {
	th := plazo.DefaultThresholds
	if plazo.EvaluateSemaforo(0, th) != plazo.Vencido {
		t.Error("0 days should be Vencido")
	}
	if plazo.EvaluateSemaforo(1, th) != plazo.Rojo {
		t.Error("1 day should be Rojo")
	}
	if plazo.EvaluateSemaforo(3, th) != plazo.Amarillo {
		t.Error("3 days should be Amarillo")
	}
	if plazo.EvaluateSemaforo(6, th) != plazo.Verde {
		t.Error("6 days should be Verde")
	}
}
