package caso

import (
	"errors"
	"time"

	"poly.app/api/internal/domain/estado"
)

// MedioPago identifies the payment instrument of the disputed transaction.
type MedioPago string

const (
	TarjetaCredito MedioPago = "TARJETA_CREDITO"
	TarjetaDebito  MedioPago = "TARJETA_DEBITO"
	Transferencia  MedioPago = "TRANSFERENCIA"
	Cajero         MedioPago = "CAJERO"
)

// RelacionOperacion indicates the relationship between the account holder and the counterpart.
type RelacionOperacion string

const (
	CuentaPropia RelacionOperacion = "CUENTA_PROPIA"
	Familiar     RelacionOperacion = "FAMILIAR"
	Tercero      RelacionOperacion = "TERCERO"
)

// Rol of a user within the estudio.
type Rol string

const (
	RolAbogado    Rol = "ABOGADO"
	RolTramitador Rol = "TRAMITADOR"
	RolAdmin      Rol = "ADMIN"
)

// Caso is the central aggregate of the domain.
type Caso struct {
	ID             string
	EstudioID      string
	BancoID        string
	ClienteID      string
	AbogadoID      *string
	NumeroOT       *string
	Estado         estado.Estado
	FechaDJ        time.Time
	FechaDenuncia  *time.Time
	DenunciaValida bool
	MotivoTermino  *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ValidateTransition checks that moving to target is allowed and that
// business pre-conditions for the transition are met.
func (c *Caso) ValidateTransition(target estado.Estado) error {
	if err := estado.Transition(c.Estado, target); err != nil {
		return err
	}
	// JUDICIALIZACION requires a valid denuncia on record.
	if target == estado.Judicializacion && !c.DenunciaValida {
		return errors.New("valid denuncia required to enter JUDICIALIZACION")
	}
	return nil
}

// RequiresTerminationReason reports whether transitioning to target requires a motivo_termino.
func RequiresTerminationReason(target estado.Estado) bool {
	return target == estado.Terminado
}
