package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type WatchlistAdapter interface {
	SearchIndividual(documentNumber, firstName, lastName, fullName, countryCode string, birthYear int) (*entity.Watchlist, error)
	SearchCompany(legalName, countryCode string) (*entity.Watchlist, error)
}
