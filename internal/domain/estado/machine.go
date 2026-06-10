package estado

import "fmt"

type Estado string

const (
	Llamada            Estado = "LLAMADA"
	Revision           Estado = "REVISION"
	Suspension         Estado = "SUSPENSION"
	PreJudicializacion Estado = "PRE_JUDICIALIZACION"
	Restitucion        Estado = "RESTITUCION"
	Judicializacion    Estado = "JUDICIALIZACION"
	Cierre             Estado = "CIERRE"
	Terminado          Estado = "TERMINADO"
)

var transitions = map[Estado][]Estado{
	Llamada:            {Revision, Terminado},
	Revision:           {Suspension, Terminado},
	Suspension:         {PreJudicializacion, Terminado},
	PreJudicializacion: {Judicializacion, Restitucion, Terminado},
	Restitucion:        {Judicializacion, Cierre},
	Judicializacion:    {Cierre, Terminado},
	Cierre:             {},
	Terminado:          {},
}

// IsValid reports whether the string corresponds to a known estado.
func IsValid(s string) bool {
	_, ok := transitions[Estado(s)]
	return ok
}

// Transition validates that moving from current to target is an allowed transition.
// Returns an error if the transition is not in the table.
func Transition(current, target Estado) error {
	allowed, ok := transitions[current]
	if !ok {
		return fmt.Errorf("unknown source estado: %q", current)
	}
	for _, a := range allowed {
		if a == target {
			return nil
		}
	}
	return fmt.Errorf("transition not allowed: %q → %q", current, target)
}

// AvailableTransitions returns the estados reachable from current.
func AvailableTransitions(current Estado) []Estado {
	return transitions[current]
}
