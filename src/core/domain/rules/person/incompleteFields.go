package person

import (
	"encoding/json"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/pkg/errors"
)

type incompleteFieldsAnalyzer struct {
	dateOfBirthRequired                   bool
	inputtedOrEnrichedDateOfBirthRequired bool
	phoneNumberRequired                   bool
	emailRequired                         bool
	pepRequired                           bool
	lastNameRequired                      bool
}

func NewIncompleteFieldsAnalyzer(
	dateOfBirthRequired bool,
	inputtedOrEnrichedDateOfBirthRequired bool,
	phoneNumberRequired bool,
	emailRequired bool,
	pepRequired bool,
	lastNameRequired bool,
) Analyzer {
	return &incompleteFieldsAnalyzer{
		dateOfBirthRequired:                   dateOfBirthRequired,
		inputtedOrEnrichedDateOfBirthRequired: inputtedOrEnrichedDateOfBirthRequired,
		phoneNumberRequired:                   phoneNumberRequired,
		emailRequired:                         emailRequired,
		pepRequired:                           pepRequired,
		lastNameRequired:                      lastNameRequired,
	}
}

func (ref *incompleteFieldsAnalyzer) Analyze(person entity.Person) (*entity.RuleResultV2, error) {
	requiredFieldsNotFound := entity.NewRuleResultV2(values.RuleSetIncomplete, values.RuleNameRequiredFieldsNotFound)

	result := make([]string, 0)
	if ref.dateOfBirthRequired {
		if person.Individual == nil || person.Individual.DateOfBirthInputted == nil {
			result = append(result, "Date of birth is required")
			requiredFieldsNotFound.AddProblem(values.ProblemCodeDateOfBirthRequired, "")
		}
	}

	if ref.inputtedOrEnrichedDateOfBirthRequired {
		if person.Individual == nil ||
			(person.Individual.DateOfBirth == nil &&
				(person.EnrichedInformation == nil || person.EnrichedInformation.EnrichedIndividual.BirthDate == "")) {
			result = append(result, "Inputted or Enriched Date of birth is required")
			requiredFieldsNotFound.AddProblem(values.ProblemCodeInputtedOrEnrichedDateOfBirthRequired, "")
		}
	}

	if ref.phoneNumberRequired {
		if person.Individual == nil || len(person.Individual.Phones) == 0 {
			result = append(result, "Phone number is required")
			requiredFieldsNotFound.AddProblem(values.ProblemCodePhoneRequired, "")
		}
	}

	if ref.emailRequired {
		if person.Email == "" {
			result = append(result, "Email is required")
			requiredFieldsNotFound.AddProblem(values.ProblemCodeEmailRequired, "")
		}
	}

	if ref.pepRequired {
		if person.Individual == nil || person.Individual.Pep == nil {
			result = append(result, "PEP is required")
			requiredFieldsNotFound.AddProblem(values.ProblemCodePepRequired, "")
		}
	}

	if ref.lastNameRequired {
		if person.Individual == nil || person.Individual.LastName == "" {
			result = append(result, "Last name is required")
			requiredFieldsNotFound.AddProblem(values.ProblemCodeLastNameRequired, "")
		}
	}

	if len(result) > 0 {
		metadata, err := json.Marshal(result)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		requiredFieldsNotFound.SetResult(values.ResultStatusIncomplete).SetMetadata(metadata)

		return requiredFieldsNotFound, nil
	}

	requiredFieldsNotFound.SetResult(values.ResultStatusApproved)
	return requiredFieldsNotFound, nil
}
