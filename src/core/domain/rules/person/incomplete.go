package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/pkg/errors"
)

type incompleteAnalyzer struct {
	person    entity.Person
	analyzers []Analyzer
}

func NewIncompleteAnalyzer(person entity.Person,
	analyzers []Analyzer) entity.Rule {
	return &incompleteAnalyzer{
		person:    person,
		analyzers: analyzers,
	}
}

func (ref *incompleteAnalyzer) Analyze() ([]entity.RuleResultV2, error) {
	var rulesResult []entity.RuleResultV2

	for _, analyzer := range ref.analyzers {
		result, err := analyzer.Analyze(ref.person)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		rulesResult = append(rulesResult, *result)
	}

	return rulesResult, nil
}

func (ref *incompleteAnalyzer) Name() values.RuleSet {
	return values.RuleSetIncomplete
}
