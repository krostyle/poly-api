package caso

import (
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

// MotivoTermino is the reason a caso is terminated. Enum defined by the business expert.
type MotivoTermino string

const (
	MotivoImprocedente                  MotivoTermino = "IMPROCEDENTE"
	MotivoExtemporaneo                  MotivoTermino = "EXTEMPORANEO"
	MotivoBusquedasNegativas            MotivoTermino = "BUSQUEDAS_NEGATIVAS"
	MotivoDeudorFallecido               MotivoTermino = "DEUDOR_FALLECIDO"
	MotivoDesistimientoCliente          MotivoTermino = "DESISTIMIENTO_CLIENTE"
	MotivoDesistimientoBanco            MotivoTermino = "DESISTIMIENTO_BANCO"
	MotivoDesistimientoDenunciaInvalida MotivoTermino = "DESISTIMIENTO_DENUNCIA_INVALIDA"
	MotivoDesistimientoSinDenuncia      MotivoTermino = "DESISTIMIENTO_SIN_DENUNCIA"
	MotivoSentenciaFavorableBanco       MotivoTermino = "SENTENCIA_FAVORABLE_BANCO"
	MotivoSentenciaDesfavorableBanco    MotivoTermino = "SENTENCIA_DESFAVORABLE_BANCO"
	MotivoAvenimiento                   MotivoTermino = "AVENIMIENTO"
	MotivoAbandonoProcedimiento         MotivoTermino = "ABANDONO_PROCEDIMIENTO"
)

// IsValidMotivoTermino reports whether m is a known termination reason.
func IsValidMotivoTermino(m string) bool {
	switch MotivoTermino(m) {
	case MotivoImprocedente, MotivoExtemporaneo, MotivoBusquedasNegativas,
		MotivoDeudorFallecido, MotivoDesistimientoCliente, MotivoDesistimientoBanco,
		MotivoDesistimientoDenunciaInvalida, MotivoDesistimientoSinDenuncia,
		MotivoSentenciaFavorableBanco, MotivoSentenciaDesfavorableBanco,
		MotivoAvenimiento, MotivoAbandonoProcedimiento:
		return true
	}
	return false
}

// Caso is the central aggregate of the domain.
type Caso struct {
	ID             string
	EstudioID      string
	BancoID        string
	ClienteID      string
	AbogadoID      *string
	NumeroOT       *string
	Estado         estado.Estado
	FechaDJ        *time.Time
	FechaDenuncia  *time.Time
	DenunciaValida bool
	MotivoTermino  *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ValidateTransition checks that moving to target is allowed by the state machine.
func (c *Caso) ValidateTransition(target estado.Estado) error {
	return estado.Transition(c.Estado, target)
}

// RequiresTerminationReason reports whether transitioning to target requires a motivo_termino.
func RequiresTerminationReason(target estado.Estado) bool {
	return target == estado.Terminado
}
