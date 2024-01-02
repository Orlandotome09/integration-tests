package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"github.com/pkg/errors"
)

type incompleteAddressAnalyzer struct {
}

func NewIncompleteAddressAnalyzer() Analyzer {
	return &incompleteAddressAnalyzer{}
}

func (ref *incompleteAddressAnalyzer) Analyze(person entity.Person) (*entity.RuleResultV2, error) {
	addressNotFound := entity.NewRuleResultV2(values.RuleSetIncomplete, values.RuleNameAddressNotFound)

	result := make([]string, 0)
	if len(person.Addresses) == 0 {
		result = append(result, "Address Not Found")
	}

	if len(result) > 0 {
		metadata, err := json.Marshal(result)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		addressNotFound.SetResult(values.ResultStatusIncomplete).SetMetadata(metadata).
			AddProblem(values.ProblemCodeAddressNotFound, "")
		return addressNotFound, nil
	}

	addressNotFound.SetResult(values.ResultStatusApproved)
	return addressNotFound, nil
}
