package person

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

func TestWatchlistAnalyzer_Analyze(t *testing.T) {
	t.Run("should find occurrence for individual", func(t *testing.T) {
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
		wantedSources := []string{"OFAC"}

		analyzer := NewWatchlistAnalyzer(person, true, wantedSources, nil)

		rulesResult, err := analyzer.Analyze()

		assert.Nil(t, err)
		assert.Equal(t, values.ResultStatusAnalysing, rulesResult[0].Result)
		assert.True(t, rulesResult[0].Pending)
	})

	t.Run("should find occurrence for individual and override result", func(t *testing.T) {
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
		wantedSources := []string{"OFAC"}
		hasMatchInWatchListStatus := values.ResultStatusRejected
		analyzer := NewWatchlistAnalyzer(person, true, wantedSources, &hasMatchInWatchListStatus)

		rulesResult, err := analyzer.Analyze()

		assert.Nil(t, err)
		assert.Equal(t, values.ResultStatusRejected, rulesResult[0].Result)
		assert.False(t, rulesResult[0].Pending)
	})

	t.Run("should not find occurrences for individual", func(t *testing.T) {
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
		wantedSources := []string{"PEP"}

		analyzer := NewWatchlistAnalyzer(person, true, wantedSources, nil)

		rulesResult, err := analyzer.Analyze()

		assert.Nil(t, err)
		assert.Equal(t, values.ResultStatusApproved, rulesResult[0].Result)
		assert.False(t, rulesResult[0].Pending)
	})

	t.Run("should not find occurrences for individual when wanted PEP Source", func(t *testing.T) {
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
		wantedSources := []string{"PEP"}

		analyzer := NewWatchlistAnalyzer(person, true, wantedSources, nil)

		rulesResult, err := analyzer.Analyze()

		assert.Nil(t, err)
		assert.Equal(t, values.ResultStatusApproved, rulesResult[0].Result)
		assert.False(t, rulesResult[0].Pending)
	})

	t.Run("should find occurrence for company", func(t *testing.T) {
		metadata, _ := json.Marshal([]string{"OFAC"})
		person := entity.Person{
			PersonType:     values.PersonTypeCompany,
			DocumentNumber: "1111",
			Company: &entity.Company{
				LegalName: "SomeName",
			},
			Watchlist: &entity.Watchlist{
				Sources: []string{"OFAC"}, Metadata: metadata,
			},
		}
		wantedSources := []string{"OFAC"}

		analyzer := NewWatchlistAnalyzer(person, false, wantedSources, nil)

		rulesResult, err := analyzer.Analyze()

		assert.Nil(t, err)
		assert.Equal(t, values.ResultStatusAnalysing, rulesResult[0].Result)
		assert.True(t, rulesResult[0].Pending)
	})
}
