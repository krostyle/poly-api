package casos

import (
	"context"
	"time"

	"github.com/google/uuid"
	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/caso"
	"poly.app/api/internal/domain/estado"
)

type CreateCaseInput struct {
	EstudioID       string
	BancoID         string
	ClienteRUT      string
	ClienteNombre   string
	ClienteContacto *string
	FechaDJ         time.Time
	UsuarioID       string
}

type CreateCaseUseCase struct {
	casos    domain.CasoRepository
	clientes domain.ClienteRepository
	auditor  domain.AuditLogger
}

func NewCreateCaseUseCase(casos domain.CasoRepository, clientes domain.ClienteRepository, auditor domain.AuditLogger) *CreateCaseUseCase {
	return &CreateCaseUseCase{casos: casos, clientes: clientes, auditor: auditor}
}

func (uc *CreateCaseUseCase) Execute(ctx context.Context, in CreateCaseInput) (*domain.CasoDetalle, error) {
	cliente, err := uc.clientes.Upsert(ctx, domain.UpsertClienteInput{
		EstudioID: in.EstudioID,
		BancoID:   in.BancoID,
		RUT:       in.ClienteRUT,
		Nombre:    in.ClienteNombre,
		Contacto:  in.ClienteContacto,
	})
	if err != nil {
		return nil, err
	}

	c := &caso.Caso{
		ID:        uuid.New().String(),
		EstudioID: in.EstudioID,
		BancoID:   in.BancoID,
		ClienteID: cliente.ID,
		Estado:    estado.Llamada,
		FechaDJ:   in.FechaDJ,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := uc.casos.Create(ctx, c); err != nil {
		return nil, err
	}

	uid := in.UsuarioID
	_ = uc.auditor.Log(ctx, domain.AuditEntry{
		EstudioID: in.EstudioID,
		UsuarioID: &uid,
		CasoID:    &c.ID,
		Accion:    "CASO_CREADO",
		Detalle:   map[string]any{"estado": string(c.Estado), "cliente_id": cliente.ID},
	})

	return &domain.CasoDetalle{
		Caso:        c,
		Cliente:     cliente,
		Operaciones: nil,
	}, nil
}
