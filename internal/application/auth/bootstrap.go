package auth

import (
	"context"

	"poly.app/api/internal/domain"
)

type BootstrapInput struct {
	ClerkOrgID  string
	ClerkUserID string
	OrgName     string
	UserName    string
	UserEmail   string
}

type BootstrapOutput struct {
	Estudio *domain.Estudio
	Usuario *domain.Usuario
	Bancos  []*domain.Banco
}

type BootstrapUseCase struct {
	estudios domain.EstudioRepository
	usuarios domain.UsuarioRepository
	bancos   domain.BancoRepository
}

func NewBootstrapUseCase(
	estudios domain.EstudioRepository,
	usuarios domain.UsuarioRepository,
	bancos domain.BancoRepository,
) *BootstrapUseCase {
	return &BootstrapUseCase{estudios: estudios, usuarios: usuarios, bancos: bancos}
}

// Execute is idempotent: upserts estudio and usuario, returns full profile.
func (uc *BootstrapUseCase) Execute(ctx context.Context, in BootstrapInput) (*BootstrapOutput, error) {
	estudio, err := uc.estudios.UpsertByClerkOrgID(ctx, in.ClerkOrgID, in.OrgName)
	if err != nil {
		return nil, err
	}

	usuario, err := uc.usuarios.UpsertByClerkUserID(ctx, domain.UpsertUsuarioInput{
		ClerkUserID: in.ClerkUserID,
		EstudioID:   estudio.ID,
		Nombre:      in.UserName,
		Email:       in.UserEmail,
		Rol:         "ADMIN",
	})
	if err != nil {
		return nil, err
	}

	bancos, err := uc.bancos.List(ctx, estudio.ID)
	if err != nil {
		return nil, err
	}

	return &BootstrapOutput{Estudio: estudio, Usuario: usuario, Bancos: bancos}, nil
}
