package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"github.com/pkg/errors"
)

type pepAnalyzer struct {
	person entity.Person
}

func NewPepAnalyzer(person entity.Person) entity.Rule {
	return &pepAnalyzer{
		person: person,
	}
}

func (ref *pepAnalyzer) Analyze() ([]entity.RuleResultV2, error) {
	results, err := ref.validate()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if results != nil {
		return results, err
	}

	occurrenceInPepFields := entity.NewRuleResultV2(values.RuleSetPep, values.RuleNamePep)

	pepSource := make([]string, 0)

	var isPep bool
	for _, source := range ref.person.Watchlist.Sources {
		if source == "PEP" {
			pepSource = append(pepSource, values.PepSourceWatchlist)
			isPep = true
			break
		}
	}

	if ref.person.Individual.Pep != nil && *ref.person.Individual.Pep {
		pepSource = append(pepSource, values.PepSourceSelfDeclared)
		isPep = true
	}

	if ref.person.PEPInformation != nil {
		pepSource = append(pepSource, values.PepSourceCOAF)
		isPep = true
	}

	if isPep {
		metadata := ref.buildMetadata(pepSource)
		occurrenceInPepFields.
			SetPending(true).
			SetMetadata(metadata).
			SetResult(values.ResultStatusAnalysing)
	} else {
		occurrenceInPepFields.
			SetPending(false).
			SetResult(values.ResultStatusApproved)
	}

	return []entity.RuleResultV2{*occurrenceInPepFields}, nil
}

func (ref *pepAnalyzer) validate() ([]entity.RuleResultV2, error) {
	if ref.person.PersonType == values.PersonTypeCompany {
		return nil, errors.New("PEP person must be individual")
	}

	if ref.person.Watchlist == nil {
		metadata, _ := json.Marshal("Date of birth is not present in profile and was not enriched")
		occurrenceInPep := entity.NewRuleResultV2(values.RuleSetPep, values.RuleNamePep).
			SetMetadata(metadata).SetResult(values.ResultStatusAnalysing).SetPending(true).
			AddProblem(values.ProblemCodeDateOfBirthNotInputtedOrEnriched, "")

		return []entity.RuleResultV2{*occurrenceInPep}, nil
	}

	return nil, nil
}

func (ref *pepAnalyzer) Name() values.RuleSet {
	return values.RuleSetPep
}

func (ref *pepAnalyzer) buildMetadata(pepSources []string) []byte {

	type Metadata struct {
		PepSources     []string               `json:"pep_sources"`
		PepInformation *entity.PepInformation `json:"pep_information,omitempty"`
	}

	metadata := Metadata{
		PepSources:     pepSources,
		PepInformation: ref.person.PEPInformation,
	}

	metadataJson, _ := json.Marshal(metadata)
	return metadataJson
}
