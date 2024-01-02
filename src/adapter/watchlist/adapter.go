package watchlistAdapter

import (
	watchlistClient "bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/http"
	watchlistTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/translator"
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type watchlistAdapter struct {
	watchlistClient watchlistClient.WatchListClient
	translator      watchlistTranslator.WatchListTranslator
}

func New(
	watchlistClient watchlistClient.WatchListClient,
	translator watchlistTranslator.WatchListTranslator) _interfaces.WatchlistAdapter {
	return &watchlistAdapter{
		watchlistClient: watchlistClient,
		translator:      translator,
	}
}

func (ref *watchlistAdapter) SearchIndividual(documentNumber, firstName, lastName, fullName, countryCode string,
	birthYear int) (*entity.Watchlist, error) {

	responses, err := ref.watchlistClient.SearchIndividual(documentNumber, firstName, lastName, fullName, countryCode, birthYear)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sources, err := ref.translator.ToDomain(responses)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sources, nil
}
func (ref *watchlistAdapter) SearchCompany(legalName, countryCode string) (*entity.Watchlist, error) {

	responses, err := ref.watchlistClient.SearchCompany(legalName, countryCode)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	sources, err := ref.translator.ToDomain(responses)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return sources, nil
}
