package _interfaces

import "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"

type RestrictiveListsAdapter interface {
	OccurrenceInBlackList(document string, name string) (*entity.BlacklistStatus, error)
	OccurrenceInPepList(documentNumber string) (*entity.PepInformation, error)
}
