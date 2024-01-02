package watchlistConstructor

import (
	"bitbucket.org/bexstech/temis-compliance/src/core"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
	"strings"
)

type watchlistConstructor struct {
	watchlistService interfaces.WatchlistAdapter
}

func New(service interfaces.WatchlistAdapter) interfaces.PersonConstructor {
	return &watchlistConstructor{watchlistService: service}
}

func (ref *watchlistConstructor) Assemble(personWrapper *entity.PersonWrapper) error {
	if !personWrapper.Person.ShouldGetWatchlist() {
		return nil
	}

	var watchlist *entity.Watchlist = nil

	if personWrapper.Person.IsValidWatchlistIndividual() {
		documentNumber := core.NormalizeDocument(personWrapper.Person.DocumentNumber)
		firstName, lastName, fullName := getNames(personWrapper)
		result, err := ref.watchlistService.SearchIndividual(documentNumber,
			firstName, lastName, fullName, extractCountryCode(personWrapper), personWrapper.Person.Individual.DateOfBirth.Year())
		if err != nil {
			return errors.WithStack(err)
		}

		watchlist = result
	} else {
		if personWrapper.Person.IsValidWatchlistCompany() {
			result, err := ref.watchlistService.SearchCompany(personWrapper.Person.Company.LegalName, "BR")
			if err != nil {
				return errors.WithStack(err)
			}

			watchlist = result
		}
	}

	personWrapper.Mutex.Lock()
	defer personWrapper.Mutex.Unlock()
	personWrapper.Person.Watchlist = watchlist

	return nil
}

func getNames(personWrapper *entity.PersonWrapper) (first, last, full string) {
	if personWrapper.Person.EnrichedInformation != nil && personWrapper.Person.EnrichedInformation.Name != "" {
		first, last = extractFirstLastName(personWrapper.Person.EnrichedInformation.Name)
		return first, last, personWrapper.Person.EnrichedInformation.Name
	}

	if personWrapper.Person.Name != "" {
		first, last = extractFirstLastName(personWrapper.Person.Name)
		return first, last, personWrapper.Person.Name
	}

	return "", "", ""
}

func extractFirstLastName(fullName string) (string, string) {
	names := strings.Fields(fullName)
	if len(names) == 1 {
		return names[0], ""
	}
	return names[0], names[len(names)-1]
}

func extractCountryCode(personWrapper *entity.PersonWrapper) string {
	if personWrapper != nil && personWrapper.Person.Individual != nil && len(personWrapper.Person.Individual.Nationality) == 2 {
		return personWrapper.Person.Individual.Nationality
	}
	return ""
}
