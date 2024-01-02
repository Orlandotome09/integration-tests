package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCafAnalyzer_should_approved_when_caf_approved(t *testing.T) {
	enrichedInformation := &entity.EnrichedInformation{
		Providers: []entity.Provider{
			{
				ProviderName: entity.ProviderCAF,
				Status:       "APPROVED",
			},
		},
	}

	cafAnalyzer := NewCafAnalyzer(enrichedInformation)

	expected := []entity.RuleResultV2{
		{
			RuleSet:  values.RuleSetCafAnalysis,
			RuleName: values.RuleNameCafAnalysis,
			Result:   values.ResultStatusApproved,
		},
	}

	received, err := cafAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestCafAnalyzer_should_reject_when_caf_rejected(t *testing.T) {
	enrichedInformation := &entity.EnrichedInformation{
		Providers: []entity.Provider{
			{
				ProviderName: entity.ProviderCAF,
				Status:       "REJECTED",
			},
		},
	}

	cafAnalyzer := NewCafAnalyzer(enrichedInformation)

	expected := []entity.RuleResultV2{
		{
			RuleSet:  values.RuleSetCafAnalysis,
			RuleName: values.RuleNameCafAnalysis,
			Result:   values.ResultStatusRejected,
		},
	}

	received, err := cafAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestCafAnalyzer_should_set_to_analysis_when_caf_analysing(t *testing.T) {
	requestID := uuid.New().String()
	enrichedInformation := &entity.EnrichedInformation{
		Providers: []entity.Provider{
			{
				ProviderName: entity.ProviderCAF,
				RequestID:    requestID,
				Status:       "ANALYSING",
			},
		},
	}

	cafAnalyzer := NewCafAnalyzer(enrichedInformation)

	expected := []entity.RuleResultV2{
		{
			RuleSet:  values.RuleSetCafAnalysis,
			RuleName: values.RuleNameCafAnalysis,
			Result:   values.ResultStatusAnalysing,
			Pending:  true,
			Problems: []entity.Problem{{
				Code:   values.ProblemCodeCafAnalysisPending,
				Detail: []string{requestID},
			}},
		},
	}

	received, err := cafAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestCafAnalyzer_should_set_to_analysis_when_not_found_caf_provider_analysis(t *testing.T) {
	enrichedInformation := &entity.EnrichedInformation{
		Providers: []entity.Provider{
			{
				ProviderName: "BUREAU_ENRICHER",
				Status:       "ANALYSING",
			},
		},
	}

	cafAnalyzer := NewCafAnalyzer(enrichedInformation)

	expected := []entity.RuleResultV2{
		{
			RuleSet:  values.RuleSetCafAnalysis,
			RuleName: values.RuleNameCafAnalysis,
			Result:   values.ResultStatusAnalysing,
			Pending:  true,
			Problems: []entity.Problem{{
				Code:   values.ProblemCodeNotFoundCafAnalysis,
				Detail: "",
			}},
		},
	}

	received, err := cafAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestCafAnalyzer_should_set_to_analysis_when_not_found_enriched_information(t *testing.T) {
	var enrichedInformation *entity.EnrichedInformation

	cafAnalyzer := NewCafAnalyzer(enrichedInformation)

	expected := []entity.RuleResultV2{
		{
			RuleSet:  values.RuleSetCafAnalysis,
			RuleName: values.RuleNameCafAnalysis,
			Result:   values.ResultStatusAnalysing,
			Pending:  true,
			Problems: []entity.Problem{{
				Code:   values.ProblemCodeNotFoundEnrichedInformation,
				Detail: "",
			}},
		},
	}

	received, err := cafAnalyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}
