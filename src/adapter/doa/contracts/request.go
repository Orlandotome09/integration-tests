package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type DOAExtractionRequest struct {
	ProfileID      uuid.UUID       `json:"profile_id" binding:"required"`
	Type           DocumentType    `json:"document_type" binding:"required"`
	Metadata       *MetadataDto    `json:"metadata"`
	FileParams     []FileParams    `json:"file_params"`
	OCRParams      *OCRParams      `json:"ocr_params,omitempty"`
	ScoreParams    *ScoreParams    `json:"score_params,omitempty"`
	CallbackParams *CallbackParams `json:"callback_params,omitempty"`
}

type MetadataDto struct {
	Name            string `json:"name" binding:"required"`
	MothersName     string `json:"mothers_name,omitempty"`
	Birthday        string `json:"birthday,omitempty"`
	Number          string `json:"document_number" binding:"required"`
	IssueDate       string `json:"issue_date" binding:"required"`
	IdentidadeTexto string `json:"identification,omitempty"`
}

type FileParams struct {
	FileID   string   `json:"file_id" binding:"required"`
	FileSide FileSide `json:"file_side" binding:"required"`
	FileURI  string   `json:"file_uri" binding:"required"`
	FileName string   `json:"file_name" binding:"required"`
}

type OCRParams struct {
	Provider OCRProvider `json:"provider"`
}

type ScoreParams struct {
	Calculator string `json:"calculator"`
}

type CallbackParams struct {
	URL string `json:"url"`
}

type DocumentType string

const (
	RG  DocumentType = "RG"
	CNH DocumentType = "CNH"
)

type FileSide string

const (
	FrontSide     FileSide = "FRONT"
	BackSide      FileSide = "BACK"
	BackFrontSide FileSide = "ALL"
)

type OCRProvider string

const (
	SIMPLY OCRProvider = "SIMPLY"
)

func ToDocumentType(d string) DocumentType {
	switch d {
	case "RG":
		return RG
	case "CNH":
		return CNH
	default:
		return ""

	}
}

func ToFileSide(d values.FileSide) FileSide {
	switch d {
	case values.FileSideFront:
		return FrontSide
	case values.FileSideBack:
		return BackSide
	case values.FileSideBackFront:
		return BackFrontSide
	default:
		return ""
	}
}

func ToMetadata(fields entity.DocumentFields) *MetadataDto {

	metadata := &MetadataDto{
		Name:      fields.Name,
		Number:    fields.Number,
		IssueDate: fields.IssueDate,
	}
	return metadata

}
