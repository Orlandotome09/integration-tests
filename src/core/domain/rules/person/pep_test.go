package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPepAnalyzer_Analyze_HasWatchlistPepSource(t *testing.T) {
	dateOfBirth := time.Now()
	metadata, _ := json.Marshal([]string{"OFAC"})
	person := entity.Person{
		PersonType:     values.PersonTypeIndividual,
		DocumentNumber: "1111",
		Individual: &entity.Individual{
			FirstName:   "Dilma",
			LastName:    "Roussef",
			DateOfBirth: &dateOfBirth,
		},
		Watchlist: &entity.Watchlist{
			Sources: []string{"PEP"}, Metadata: metadata,
		},
	}

	expected := values.ResultStatusAnalysing
	analyzer := NewPepAnalyzer(person)

	rulesResult, err := analyzer.Analyze()
	occurrenceInPep := rulesResult[0]

	assert.Nil(t, err)
	assert.Equal(t, expected, occurrenceInPep.Result)
	assert.Equal(t, values.RuleNamePep, occurrenceInPep.RuleName)
	assert.True(t, occurrenceInPep.Pending)
}

func TestPepAnalyzer_Analyze_PepPerson(t *testing.T) {
	dateOfBirth := time.Now()
	metadata, _ := json.Marshal([]string{"OFAC"})
	pep := true
	person := entity.Person{
		PersonType:     values.PersonTypeIndividual,
		DocumentNumber: "1111",
		Individual: &entity.Individual{
			FirstName:   "Dilma",
			LastName:    "Roussef",
			DateOfBirth: &dateOfBirth,
			Pep:         &pep,
		},
		Watchlist: &entity.Watchlist{
			Sources: []string{"OFAC"}, Metadata: metadata,
		},
	}

	expected := values.ResultStatusAnalysing
	analyzer := NewPepAnalyzer(person)

	rulesResult, err := analyzer.Analyze()
	occurrenceInPep := rulesResult[0]

	assert.Nil(t, err)
	assert.Equal(t, expected, occurrenceInPep.Result)
	assert.Equal(t, values.RuleNamePep, occurrenceInPep.RuleName)
	assert.True(t, occurrenceInPep.Pending)
}

func TestPepAnalyzer_Analyze_NotPepApproved(t *testing.T) {
	dateOfBirth := time.Now()
	metadata, _ := json.Marshal([]string{"OFAC"})
	person := entity.Person{
		PersonType:     values.PersonTypeIndividual,
		DocumentNumber: "1111",
		Individual: &entity.Individual{
			FirstName:   "Dilma",
			LastName:    "Roussef",
			DateOfBirth: &dateOfBirth,
		},
		Watchlist: &entity.Watchlist{
			Sources: []string{"OFAC"}, Metadata: metadata,
		},
	}
	
	expected := values.ResultStatusApproved
	analyzer := NewPepAnalyzer(person)

	rulesResult, err := analyzer.Analyze()
	occurrenceInPep := rulesResult[0]

	assert.Nil(t, err)
	assert.Equal(t, expected, occurrenceInPep.Result)
	assert.Equal(t, values.RuleNamePep, occurrenceInPep.RuleName)
	assert.False(t, occurrenceInPep.Pending)
}
