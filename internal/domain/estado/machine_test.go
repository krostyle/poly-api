package estado_test

import (
	"testing"

	"poly.app/api/internal/domain/estado"
)

func TestValidTransitions(t *testing.T) {
	cases := []struct {
		current estado.Estado
		target  estado.Estado
	}{
		{estado.Llamada, estado.Revision},
		{estado.Llamada, estado.Terminado},
		{estado.Revision, estado.Suspension},
		{estado.Suspension, estado.PreJudicializacion},
		{estado.PreJudicializacion, estado.Judicializacion},
		{estado.PreJudicializacion, estado.Restitucion},
		{estado.PreJudicializacion, estado.Terminado},
		{estado.Restitucion, estado.Judicializacion},
		{estado.Restitucion, estado.Cierre},
		{estado.Judicializacion, estado.Cierre},
		{estado.Judicializacion, estado.Terminado},
	}
	for _, tc := range cases {
		if err := estado.Transition(tc.current, tc.target); err != nil {
			t.Errorf("expected valid transition %s→%s, got: %v", tc.current, tc.target, err)
		}
	}
}

func TestInvalidTransitions(t *testing.T) {
	cases := []struct {
		current estado.Estado
		target  estado.Estado
	}{
		{estado.Llamada, estado.Judicializacion},
		{estado.Cierre, estado.Llamada},
		{estado.Terminado, estado.Revision},
		{estado.Revision, estado.Cierre},
	}
	for _, tc := range cases {
		if err := estado.Transition(tc.current, tc.target); err == nil {
			t.Errorf("expected error for invalid transition %s→%s", tc.current, tc.target)
		}
	}
}
