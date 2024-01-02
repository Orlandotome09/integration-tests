package person

import (
	"testing"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
)

func TestBlacklistAnalyzer_Analyze_IsPresentInBlackList(t *testing.T) {
	person := entity.Person{
		DocumentNumber:  "111",
		PartnerID:       "2222",
		BlacklistStatus: &entity.BlacklistStatus{},
	}
	analyzer := NewBlackListAnalyzer(person)

	rulesResult, err := analyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusAnalysing, rulesResult[0].Result)
	assert.True(t, rulesResult[0].Pending)
}

func TestBlacklistAnalyzer_Analyze_IsNoPresentInBlackList(t *testing.T) {
	person := entity.Person{
		DocumentNumber:  "333",
		PartnerID:       "444",
		BlacklistStatus: nil,
	}
	analyzer := NewBlackListAnalyzer(person)

	rulesResult, err := analyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusApproved, rulesResult[0].Result)
	assert.False(t, rulesResult[0].Pending)
}
