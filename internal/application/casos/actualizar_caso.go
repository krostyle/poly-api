package casos

import (
	"context"
	"time"

	"poly.app/api/internal/domain"
	"poly.app/api/internal/domain/caso"
)

type UpdateCaseInput struct {
	EstudioID          string
	CasoID             string
	UsuarioID          string
	AbogadoID          *string
	NumeroOT           *string
	EstadoDenuncia     *caso.EstadoDenuncia
	FechaDenuncia      *time.Time
	FechaDJ            *time.Time
	ClearFechaDenuncia bool
	NumeroRol          *string
	Tribunal           *string
	Region             *string
	ResultadoJPL       *caso.ResultadoJPL
	FechaResolucionJPL *time.Time
	ClearResultadoJPL  bool
}

type UpdateCaseUseCase struct {
	casos    domain.CasoRepository
	plazos   domain.PlazoRepository
	feriados domain.FeriadoProvider
	auditor  domain.AuditLogger
}

func NewUpdateCaseUseCase(
	casos domain.CasoRepository,
	plazos domain.PlazoRepository,
	feriados domain.FeriadoProvider,
	auditor domain.AuditLogger,
) *UpdateCaseUseCase {
	return &UpdateCaseUseCase{casos: casos, plazos: plazos, feriados: feriados, auditor: auditor}
}

func (uc *UpdateCaseUseCase) Execute(ctx context.Context, in UpdateCaseInput) (*domain.CasoDetalle, error) {
	c, err := uc.casos.GetByID(ctx, in.EstudioID, in.CasoID)
	if err != nil {
		return nil, err
	}

	if in.AbogadoID != nil {
		c.AbogadoID = in.AbogadoID
	}
	if in.NumeroOT != nil {
		c.NumeroOT = in.NumeroOT
	}
	if in.EstadoDenuncia != nil {
		c.EstadoDenuncia = *in.EstadoDenuncia
	}
	if in.ClearFechaDenuncia {
		c.FechaDenuncia = nil
	} else if in.FechaDenuncia != nil {
		c.FechaDenuncia = in.FechaDenuncia
	}
	if in.FechaDJ != nil {
		c.FechaDJ = in.FechaDJ
	}
	if in.NumeroRol != nil {
		c.NumeroRol = in.NumeroRol
	}
	if in.Tribunal != nil {
		c.Tribunal = in.Tribunal
	}
	if in.Region != nil {
		c.Region = in.Region
	}
	if in.ClearResultadoJPL {
		c.ResultadoJPL = nil
		c.FechaResolucionJPL = nil
	} else {
		if in.ResultadoJPL != nil {
			c.ResultadoJPL = in.ResultadoJPL
		}
		if in.FechaResolucionJPL != nil {
			c.FechaResolucionJPL = in.FechaResolucionJPL
		}
	}

	if err := uc.casos.Update(ctx, c); err != nil {
		return nil, err
	}

	cambios := map[string]any{}
	if in.AbogadoID != nil {
		if in.AbogadoID == nil || *in.AbogadoID == "" {
			cambios["abogado"] = nil
		} else {
			cambios["abogado"] = *in.AbogadoID
		}
	}
	if in.NumeroOT != nil {
		cambios["numero_ot"] = *in.NumeroOT
	}
	if in.EstadoDenuncia != nil {
		cambios["estado_denuncia"] = string(*in.EstadoDenuncia)
	}
	if in.ClearFechaDenuncia {
		cambios["fecha_denuncia"] = nil
	} else if in.FechaDenuncia != nil {
		cambios["fecha_denuncia"] = in.FechaDenuncia.Format("2006-01-02")
	}
	if in.FechaDJ != nil {
		cambios["fecha_dj"] = in.FechaDJ.Format("2006-01-02")
	}
	if in.NumeroRol != nil {
		cambios["numero_rol"] = *in.NumeroRol
	}
	if in.Tribunal != nil {
		cambios["tribunal"] = *in.Tribunal
	}
	if in.Region != nil {
		cambios["region"] = *in.Region
	}
	if in.ClearResultadoJPL {
		cambios["resultado_jpl"] = nil
		cambios["fecha_resolucion_jpl"] = nil
	} else {
		if in.ResultadoJPL != nil {
			cambios["resultado_jpl"] = string(*in.ResultadoJPL)
		}
		if in.FechaResolucionJPL != nil {
			cambios["fecha_resolucion_jpl"] = in.FechaResolucionJPL.Format("2006-01-02")
		}
	}

	uid := in.UsuarioID
	_ = uc.auditor.Log(ctx, domain.AuditEntry{
		EstudioID: in.EstudioID,
		UsuarioID: &uid,
		CasoID:    &in.CasoID,
		Accion:    "CASO_ACTUALIZADO",
		Detalle:   map[string]any{"cambios": cambios},
	})

	return uc.casos.GetDetalle(ctx, in.EstudioID, in.CasoID)
}

