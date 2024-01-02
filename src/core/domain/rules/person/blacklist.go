package person

import (
	"encoding/json"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/pkg/errors"
)

type blacklistAnalyzer struct {
	person entity.Person
}

func NewBlackListAnalyzer(person entity.Person) entity.Rule {
	return &blacklistAnalyzer{
		person: person,
	}
}

func (ref *blacklistAnalyzer) Analyze() ([]entity.RuleResultV2, error) {
	occurrenceInBlacklist := entity.NewRuleResultV2(values.RuleSetBlacklist, values.RuleNameOccurrenceInBlacklist)

	if ref.person.BlacklistStatus != nil {
		metadata, err := json.Marshal(ref.person.BlacklistStatus.Justification)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		occurrenceInBlacklist.SetResult(values.ResultStatusAnalysing).SetPending(true).SetMetadata(metadata)
		return []entity.RuleResultV2{*occurrenceInBlacklist}, nil
	}

	occurrenceInBlacklist.SetResult(values.ResultStatusApproved)
	return []entity.RuleResultV2{*occurrenceInBlacklist}, nil
}

func (ref *blacklistAnalyzer) Name() values.RuleSet {
	return values.RuleSetBlacklist
}
