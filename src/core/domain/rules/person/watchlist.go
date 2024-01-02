package person

import (
	"encoding/json"

	"github.com/pkg/errors"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

const (
	sourcePEP = "PEP"
)

type watchlistAnalyzer struct {
	person                    entity.Person
	wantPepTag                bool
	wantedSources             []string
	hasMatchInWatchListStatus *values.Result
}

func NewWatchlistAnalyzer(person entity.Person,
	wantPepTag bool,
	wantedSources []string,
	hasMatchInWatchListStatus *values.Result) entity.Rule {
	return &watchlistAnalyzer{
		person:                    person,
		wantPepTag:                wantPepTag,
		wantedSources:             wantedSources,
		hasMatchInWatchListStatus: hasMatchInWatchListStatus,
	}
}

func (ref *watchlistAnalyzer) Analyze() ([]entity.RuleResultV2, error) {
	results, err := ref.validate()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if results != nil {
		return results, err
	}

	occurrenceInWatchlist := entity.NewRuleResultV2(values.RuleSetWatchlist, values.RuleNameOccurrenceInWatchlist).
		SetMetadata(ref.person.Watchlist.Metadata)

	var wantedSources = ref.wantedSources
	var sources = removePEPSource(ref.person.Watchlist.Sources)
	var foundWantedSources bool

	for _, source := range sources {
		for _, wantedSource := range wantedSources {
			if source == wantedSource {
				foundWantedSources = true
			}
		}
	}

	occurrenceInWatchlist.SetResult(values.ResultStatusApproved)

	if foundWantedSources {
		if ref.hasMatchInWatchListStatus != nil {
			occurrenceInWatchlist.SetPending(false)
			occurrenceInWatchlist.SetResult(*ref.hasMatchInWatchListStatus)
		} else {
			occurrenceInWatchlist.SetPending(true)
			occurrenceInWatchlist.SetResult(values.ResultStatusAnalysing)
		}
		occurrenceInWatchlist.AddProblem(values.ProblemCodePersonFoundOnWatchlist, ref.person.Watchlist.Sources)
	}

	return []entity.RuleResultV2{*occurrenceInWatchlist}, nil
}

func (ref *watchlistAnalyzer) validate() ([]entity.RuleResultV2, error) {
	if ref.person.PersonType == values.PersonTypeCompany {
		if ref.person.Company == nil {
			return nil, errors.New("Company is empty")
		}

		if ref.person.Company.LegalName == "" {
			return nil, errors.New("Company Legal Name is required for Watchlist rule")
		}
	}

	if ref.person.Watchlist == nil {
		metadata, _ := json.Marshal("Date of birth is not present in profile and was not enriched")
		occurrenceInWatchlist := entity.NewRuleResultV2(values.RuleSetWatchlist, values.RuleNameOccurrenceInWatchlist).
			SetMetadata(metadata).SetResult(values.ResultStatusAnalysing).SetPending(true).
			AddProblem(values.ProblemCodeDateOfBirthNotInputtedOrEnriched, "")

		return []entity.RuleResultV2{*occurrenceInWatchlist}, nil
	}

	return nil, nil
}

func (ref *watchlistAnalyzer) Name() values.RuleSet {
	return values.RuleSetWatchlist
}

func removePEPSource(sources []string) []string {
	var filtered []string

	for _, source := range sources {
		if source != sourcePEP {
			filtered = append(filtered, source)
		}
	}

	return filtered
}
