package idgenerator

import (
	"bitbucket.org/bexstech/temis-compliance/src/core"
	"github.com/google/uuid"
)

type IdGenerator interface {
	Generate(legalEntityID string, documentNumber string) uuid.UUID
}

type idgenerator struct {
}

func NewIdGenerator() IdGenerator {
	return &idgenerator{}
}

func (ref *idgenerator) Generate(legalEntityID string, documentNumber string) uuid.UUID {
	shareholderID := uuid.NewSHA1(core.GetUuidNamespace(), []byte(documentNumber+legalEntityID))
	return shareholderID
}
