package estado

import "fmt"

type Estado string

// Los 12 estados del flujo Ley 20.009, en orden cronológico típico.
const (
	Ingreso          Estado = "INGRESO"
	Revision         Estado = "REVISION"
	Prejudicial      Estado = "PREJUDICIAL"
	PagoNormativo    Estado = "PAGO_NORMATIVO"
	Judicial         Estado = "JUDICIAL"
	Audiencia        Estado = "AUDIENCIA"
	Sentencia        Estado = "SENTENCIA"
	Apelacion        Estado = "APELACION"
	SentenciaSegunda Estado = "SENTENCIA_SEGUNDA"
	Cumplimiento     Estado = "CUMPLIMIENTO"
	Terminado        Estado = "TERMINADO"
	Cierre           Estado = "CIERRE"
)

// transitions define el grafo de transiciones permitidas.
// CIERRE viene después de TERMINADO (corregido respecto al modelo original).
var transitions = map[Estado][]Estado{
	Ingreso:          {Revision, Terminado},
	Revision:         {Prejudicial, Terminado},
	Prejudicial:      {PagoNormativo, Judicial, Terminado},
	PagoNormativo:    {Judicial, Terminado},
	Judicial:         {Audiencia, Terminado},
	Audiencia:        {Sentencia, Terminado},
	Sentencia:        {Apelacion, Cumplimiento, Terminado},
	Apelacion:        {SentenciaSegunda, Terminado},
	SentenciaSegunda: {Cumplimiento, Terminado},
	Cumplimiento:     {Terminado},
	Terminado:        {Cierre},
	Cierre:           {},
}

// IsValid reports whether the string corresponds to a known estado.
func IsValid(s string) bool {
	_, ok := transitions[Estado(s)]
	return ok
}

// Transition validates that moving from current to target is an allowed transition.
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
