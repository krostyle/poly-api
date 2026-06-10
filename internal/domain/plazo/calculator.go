package plazo

import "time"

// TipoPlazo identifies the legal deadline type.
type TipoPlazo string

const (
	TipoAnalisisInterno    TipoPlazo = "ANALISIS_INTERNO"
	TipoRestitucion        TipoPlazo = "RESTITUCION"
	TipoAsignacion         TipoPlazo = "ASIGNACION"
	TipoPrecautelar        TipoPlazo = "PRECAUTELAR"
	TipoDemanda            TipoPlazo = "DEMANDA"
	TipoRestitucionRechazo TipoPlazo = "RESTITUCION_RECHAZO"
)

// Semaforo classifies the urgency level of a plazo.
type Semaforo string

const (
	Verde    Semaforo = "VERDE"
	Amarillo Semaforo = "AMARILLO"
	Rojo     Semaforo = "ROJO"
	Vencido  Semaforo = "VENCIDO"
)

// SemaforoThresholds defines remaining business-day cutoffs for each urgency level.
// Configurable; these are the spec defaults.
type SemaforoThresholds struct {
	Amarillo int // days ≤ this → AMARILLO (was VERDE)
	Rojo     int // days ≤ this → ROJO
}

var DefaultThresholds = SemaforoThresholds{
	Amarillo: 5,
	Rojo:     2,
}

func isBusinessDay(t time.Time, holidays map[time.Time]struct{}) bool {
	wd := t.Weekday()
	if wd == time.Saturday || wd == time.Sunday {
		return false
	}
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	_, isHoliday := holidays[d]
	return !isHoliday
}

func normalizeHolidays(list []time.Time) map[time.Time]struct{} {
	m := make(map[time.Time]struct{}, len(list))
	for _, h := range list {
		m[time.Date(h.Year(), h.Month(), h.Day(), 0, 0, 0, 0, time.UTC)] = struct{}{}
	}
	return m
}

// CalculateDeadline adds businessDays banking business days (Mon–Fri, excluding holidays)
// starting from the day after start.
func CalculateDeadline(start time.Time, businessDays int, holidays []time.Time) time.Time {
	holidayMap := normalizeHolidays(holidays)
	current := start
	counted := 0
	for counted < businessDays {
		current = current.AddDate(0, 0, 1)
		if isBusinessDay(current, holidayMap) {
			counted++
		}
	}
	return current
}

// RemainingBusinessDays returns how many banking business days remain between today and deadline.
// Returns a negative value if the deadline has already passed.
func RemainingBusinessDays(today, deadline time.Time, holidays []time.Time) int {
	holidayMap := normalizeHolidays(holidays)
	start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(deadline.Year(), deadline.Month(), deadline.Day(), 0, 0, 0, 0, time.UTC)

	if !end.After(start) {
		count := 0
		cur := end
		for cur.Before(start) {
			cur = cur.AddDate(0, 0, 1)
			if isBusinessDay(cur, holidayMap) {
				count--
			}
		}
		return count
	}

	count := 0
	cur := start
	for cur.Before(end) {
		cur = cur.AddDate(0, 0, 1)
		if isBusinessDay(cur, holidayMap) {
			count++
		}
	}
	return count
}

// EvaluateSemaforo classifies remaining days against the configured thresholds.
func EvaluateSemaforo(remainingDays int, thresholds SemaforoThresholds) Semaforo {
	switch {
	case remainingDays <= 0:
		return Vencido
	case remainingDays <= thresholds.Rojo:
		return Rojo
	case remainingDays <= thresholds.Amarillo:
		return Amarillo
	default:
		return Verde
	}
}
