package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewIncompleteAddressAnalyzer(t *testing.T) {
	analyzerIncompleteAddress := NewIncompleteAddressAnalyzer()

	assert.NotNil(t, analyzerIncompleteAddress)
}

func TestIncompleteAddressAnalyzer_Analyze_Success(t *testing.T) {
	analyzerIncompleteAddress := NewIncompleteAddressAnalyzer()

	profileID := uuid.New()
	addressID := uuid.New()
	person := entity.Person{ProfileID: profileID, EntityID: profileID, Addresses: []entity.Address{{
		AddressID: &addressID,
		ProfileID: &profileID,
	}}}

	expected := values.ResultStatusApproved
	result, err := analyzerIncompleteAddress.Analyze(person)

	assert.Nil(t, err)
	assert.Equal(t, expected, result.Result)
}

func TestIncompleteAddressAnalyzer_Analyze_Incomplete(t *testing.T) {
	analyzerIncompleteAddress := NewIncompleteAddressAnalyzer()

	profileID := uuid.New()
	person := entity.Person{ProfileID: profileID, EntityID: profileID}

	expected := values.ResultStatusIncomplete
	expectedCode := values.ProblemCodeAddressNotFound

	result, err := analyzerIncompleteAddress.Analyze(person)

	assert.Nil(t, err)
	assert.Equal(t, expected, result.Result)
	assert.Equal(t, expectedCode, result.Problems[0].Code)
}
