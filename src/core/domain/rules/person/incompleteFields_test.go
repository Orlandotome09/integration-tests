package person

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

func TestNewIncompleteFieldsAnalyzer(t *testing.T) {
	incompleteFieldsAnalyzer := NewIncompleteFieldsAnalyzer(
		true,
		true,
		true,
		true,
		true,
		true)

	assert.NotNil(t, incompleteFieldsAnalyzer)
}

func TestIncompleteFieldsAnalyzer_Analyze_Success(t *testing.T) {
	incompleteFieldsAnalyzer := NewIncompleteFieldsAnalyzer(
		true,
		true,
		true,
		true,
		true,
		true)

	t.Run("should approve when all required fiels present", func(t *testing.T) {
		profileID := uuid.New()
		truePointer := true
		date := time.Date(2020, 01, 01, 1, 1, 1, 1, time.UTC)
		person := entity.Person{
			ProfileID: profileID,
			Email:     "aaa",
			Individual: &entity.Individual{
				Pep:                 &truePointer,
				Phones:              []entity.Phone{{Number: "123"}},
				DateOfBirthInputted: &date,
				LastName:            "bil",
			},
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedIndividual: entity.EnrichedIndividual{
					BirthDate: "11091999",
				},
			},
		}

		expected := values.ResultStatusApproved
		result, err := incompleteFieldsAnalyzer.Analyze(person)

		assert.Nil(t, err)
		assert.Equal(t, expected, result.Result)
	})

	t.Run("should be incomplete when missing required field date of birth inputed", func(t *testing.T) {
		profileID := uuid.New()
		truePointer := true
		date := time.Date(2020, 01, 01, 1, 1, 1, 1, time.UTC)
		person := entity.Person{
			ProfileID: profileID,
			Email:     "aaa",
			Individual: &entity.Individual{
				Pep:                 &truePointer,
				Phones:              []entity.Phone{{Number: "123"}},
				DateOfBirthInputted: nil,
				LastName:            "bil",
				BureauInformation:   &entity.BureauInformation{DateOfBirth: &date},
			},
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedIndividual: entity.EnrichedIndividual{
					BirthDate: "11091999",
				},
			},
		}

		expected := values.ResultStatusIncomplete
		result, err := incompleteFieldsAnalyzer.Analyze(person)

		assert.Nil(t, err)
		assert.Equal(t, expected, result.Result)
		assert.Equal(t, values.ProblemCodeDateOfBirthRequired, result.Problems[0].Code)
	})

	t.Run("should be incomplete when missing required field date of birth inputed or enriched", func(t *testing.T) {
		profileID := uuid.New()
		truePointer := true
		date := time.Date(2020, 01, 01, 1, 1, 1, 1, time.UTC)
		person := entity.Person{
			ProfileID: profileID,
			Email:     "aaa",
			Individual: &entity.Individual{
				Pep:                 &truePointer,
				Phones:              []entity.Phone{{Number: "123"}},
				DateOfBirthInputted: &date,
				DateOfBirth:         nil,
				LastName:            "bil",
			},
			EnrichedInformation: nil,
		}

		expected := values.ResultStatusIncomplete
		result, err := incompleteFieldsAnalyzer.Analyze(person)

		assert.Nil(t, err)
		assert.Equal(t, expected, result.Result)
		assert.Equal(t, values.ProblemCodeInputtedOrEnrichedDateOfBirthRequired, result.Problems[0].Code)
	})

	t.Run("should be incomplete when missing required field phone number", func(t *testing.T) {
		profileID := uuid.New()
		truePointer := true
		date := time.Date(2020, 01, 01, 1, 1, 1, 1, time.UTC)
		person := entity.Person{
			ProfileID: profileID,
			Email:     "aaa",
			Individual: &entity.Individual{
				Pep:                 &truePointer,
				Phones:              []entity.Phone{},
				DateOfBirthInputted: &date,
				DateOfBirth:         nil,
				LastName:            "bil",
			},
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedIndividual: entity.EnrichedIndividual{
					BirthDate: "11-09-1999",
				},
			},
		}

		expected := values.ResultStatusIncomplete
		result, err := incompleteFieldsAnalyzer.Analyze(person)

		assert.Nil(t, err)
		assert.Equal(t, expected, result.Result)
		assert.Equal(t, values.ProblemCodePhoneRequired, result.Problems[0].Code)
	})

	t.Run("should be incomplete when missing required field email", func(t *testing.T) {
		profileID := uuid.New()
		truePointer := true
		date := time.Date(2020, 01, 01, 1, 1, 1, 1, time.UTC)
		person := entity.Person{
			ProfileID: profileID,
			Email:     "",
			Individual: &entity.Individual{
				Pep:                 &truePointer,
				Phones:              []entity.Phone{{Number: "1234"}},
				DateOfBirthInputted: &date,
				DateOfBirth:         nil,
				LastName:            "bil",
			},
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedIndividual: entity.EnrichedIndividual{
					BirthDate: "11-09-1999",
				},
			},
		}

		expected := values.ResultStatusIncomplete
		result, err := incompleteFieldsAnalyzer.Analyze(person)

		assert.Nil(t, err)
		assert.Equal(t, expected, result.Result)
		assert.Equal(t, values.ProblemCodeEmailRequired, result.Problems[0].Code)
	})

	t.Run("should be incomplete when missing required field pep self declation", func(t *testing.T) {
		profileID := uuid.New()
		date := time.Date(2020, 01, 01, 1, 1, 1, 1, time.UTC)
		person := entity.Person{
			ProfileID: profileID,
			Email:     "somemail",
			Individual: &entity.Individual{
				Pep:                 nil,
				Phones:              []entity.Phone{{Number: "1234"}},
				DateOfBirthInputted: &date,
				DateOfBirth:         nil,
				LastName:            "bil",
			},
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedIndividual: entity.EnrichedIndividual{
					BirthDate: "11-09-1999",
				},
			},
		}

		expected := values.ResultStatusIncomplete
		result, err := incompleteFieldsAnalyzer.Analyze(person)

		assert.Nil(t, err)
		assert.Equal(t, expected, result.Result)
		assert.Equal(t, values.ProblemCodePepRequired, result.Problems[0].Code)
	})

	t.Run("should be incomplete when missing required field last name", func(t *testing.T) {
		profileID := uuid.New()
		truePointer := true
		date := time.Date(2020, 01, 01, 1, 1, 1, 1, time.UTC)
		person := entity.Person{
			ProfileID: profileID,
			Email:     "somemail",
			Individual: &entity.Individual{
				Pep:                 &truePointer,
				Phones:              []entity.Phone{{Number: "1234"}},
				DateOfBirthInputted: &date,
				DateOfBirth:         nil,
				LastName:            "",
			},
			EnrichedInformation: &entity.EnrichedInformation{
				EnrichedIndividual: entity.EnrichedIndividual{
					BirthDate: "11-09-1999",
				},
			},
		}

		expected := values.ResultStatusIncomplete
		result, err := incompleteFieldsAnalyzer.Analyze(person)

		assert.Nil(t, err)
		assert.Equal(t, expected, result.Result)
		assert.Equal(t, values.ProblemCodeLastNameRequired, result.Problems[0].Code)
	})

	t.Run("should be incomplete when missing all required fields", func(t *testing.T) {
		profileID := uuid.New()
		person := entity.Person{ProfileID: profileID}

		expected := values.ResultStatusIncomplete
		result, err := incompleteFieldsAnalyzer.Analyze(person)

		assert.Nil(t, err)
		assert.Equal(t, expected, result.Result)
		assert.Equal(t, result.Problems[0].Code, values.ProblemCodeDateOfBirthRequired)
		assert.Equal(t, result.Problems[1].Code, values.ProblemCodeInputtedOrEnrichedDateOfBirthRequired)
		assert.Equal(t, result.Problems[2].Code, values.ProblemCodePhoneRequired)
		assert.Equal(t, result.Problems[3].Code, values.ProblemCodeEmailRequired)
		assert.Equal(t, result.Problems[4].Code, values.ProblemCodePepRequired)
		assert.Equal(t, result.Problems[5].Code, values.ProblemCodeLastNameRequired)
	})
}
