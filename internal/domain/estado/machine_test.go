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
		{estado.Ingreso, estado.Revision},
		{estado.Ingreso, estado.Terminado},
		{estado.Revision, estado.Prejudicial},
		{estado.Prejudicial, estado.PagoNormativo},
		{estado.Prejudicial, estado.Judicial},
		{estado.Prejudicial, estado.Terminado},
		{estado.PagoNormativo, estado.Judicial},
		{estado.Judicial, estado.Audiencia},
		{estado.Audiencia, estado.Sentencia},
		{estado.Sentencia, estado.Apelacion},
		{estado.Sentencia, estado.Cumplimiento},
		{estado.Terminado, estado.Cierre},
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
		{estado.Ingreso, estado.Judicial},
		{estado.Cierre, estado.Ingreso},
		{estado.Terminado, estado.Revision},
		{estado.Revision, estado.Cierre},
	}
	for _, tc := range cases {
		if err := estado.Transition(tc.current, tc.target); err == nil {
			t.Errorf("expected error for invalid transition %s→%s", tc.current, tc.target)
		}
	}
}
