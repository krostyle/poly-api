package documentos

import (
	"context"
	"errors"
	"fmt"

	"poly.app/api/internal/domain"
)

var allowedMIME = map[string]bool{
	"application/pdf":                                                   true,
	"image/jpeg":                                                        true,
	"image/png":                                                         true,
	"application/msword":                                                true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/vnd.ms-excel":                                          true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
}

var allowedTipos = map[string]bool{
	"CARTOLA": true, "EVIDENCIA": true, "DJ": true, "DENUNCIA": true,
	"CARTA_BANCO": true, "DEMANDA": true, "RESOLUCION": true, "OTRO": true,
}

const maxFileSize = 20 * 1024 * 1024 // 20 MB

type SubirDocumentoInput struct {
	CasoID      string
	Tipo        string
	Nombre      string
	Content     []byte
	ContentType string
	SubidoPor   *string
}

type SubirDocumentoUseCase struct {
	storage domain.DocumentStorage
	repo    domain.DocumentoRepository
}

func NewSubirDocumentoUseCase(storage domain.DocumentStorage, repo domain.DocumentoRepository) *SubirDocumentoUseCase {
	return &SubirDocumentoUseCase{storage: storage, repo: repo}
}

func (uc *SubirDocumentoUseCase) Execute(ctx context.Context, in SubirDocumentoInput) (*domain.Documento, error) {
	if !allowedTipos[in.Tipo] {
		return nil, errors.New("tipo de documento no permitido")
	}
	if !allowedMIME[in.ContentType] {
		return nil, errors.New("tipo de archivo no permitido")
	}
	if len(in.Content) > maxFileSize {
		return nil, errors.New("archivo supera el límite de 20 MB")
	}

	blobName := fmt.Sprintf("casos/%s/%s", in.CasoID, in.Nombre)
	url, err := uc.storage.Upload(ctx, blobName, in.Content, in.ContentType)
	if err != nil {
		return nil, fmt.Errorf("error al subir archivo: %w", err)
	}

	return uc.repo.Create(ctx, domain.NewDocumentoInput{
		CasoID:    in.CasoID,
		Tipo:      in.Tipo,
		BlobURL:   url,
		Nombre:    in.Nombre,
		SubidoPor: in.SubidoPor,
	})
}
