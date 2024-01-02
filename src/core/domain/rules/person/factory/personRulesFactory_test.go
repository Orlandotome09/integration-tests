package personRulesFactory

import (
	"os"
	"testing"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	person2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/person"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/doa"
	"github.com/stretchr/testify/assert"
)

var (
	fileService             *mocks.FileAdapter
	doaService              *doa.MockDOAService
	economicActivityService *mocks.EconomicActivityService
	factory                 interfaces.PersonRulesFactory
)

func TestMain(m *testing.M) {
	fileService = &mocks.FileAdapter{}
	doaService = &doa.MockDOAService{}
	economicActivityService = &mocks.EconomicActivityService{}
	factory = New(fileService, doaService, economicActivityService)

	os.Exit(m.Run())
}

func TestGetRules(t *testing.T) {
	minimalAge := 25
	ruleSetConfig := entity.RuleSetConfig{
		BlackListParams: &entity.BlackListParams{},
		BureauParams:    &entity.BureauParams{},
		IncompleteParams: &entity.IncompleteParams{
			DateOfBirthRequired: true,
			AddressRequired:     true,
			DocumentsRequired: []entity.DocumentRequired{{
				DocumentType:      "",
				FileRequired:      false,
				PendingOnApproval: false,
			}},
		},
		UnderAgeParams:      &entity.UnderAgeParams{MinimumAge: &minimalAge},
		WatchListParams:     &entity.WatchListParams{},
		DOAParams:           &entity.DOAParams{},
		ActivityRiskParams:  &entity.ActivityRiskParams{},
		CAFParams:           &entity.CAFParams{},
		MinimumIncomeParams: &entity.MinimumIncomeParams{},
	}
	person := entity.Person{}

	expected := []entity.Rule{
		person2.NewBlackListAnalyzer(person),
		person2.NewWatchlistAnalyzer(person,
			ruleSetConfig.WatchListParams.WantPepTag,
			ruleSetConfig.WatchListParams.WantedSources,
			ruleSetConfig.WatchListParams.HasMatchInWatchListStatus),
		person2.NewBureauAnalyzer(person,
			ruleSetConfig.BureauParams.ApprovedStatuses,
			ruleSetConfig.BureauParams.NotFoundInSerasaStatus,
			ruleSetConfig.BureauParams.NotFoundInSerasaPending,
			ruleSetConfig.BureauParams.HasProblemsInSerasaStatus,
			ruleSetConfig.BureauParams.HasProblemsInSerasaPending),
		person2.NewIncompleteAnalyzer(person,
			[]person2.Analyzer{
				person2.NewIncompleteFieldsAnalyzer(
					ruleSetConfig.IncompleteParams.DateOfBirthRequired,
					ruleSetConfig.IncompleteParams.InputtedOrEnrichedDateOfBirthRequired,
					ruleSetConfig.IncompleteParams.PhoneNumberRequired,
					ruleSetConfig.IncompleteParams.EmailRequired,
					ruleSetConfig.IncompleteParams.PepRequired,
					ruleSetConfig.IncompleteParams.LastNameRequired),
				person2.NewIncompleteAddressAnalyzer(),
				person2.NewIncompleteDocumentsAnalyzer(ruleSetConfig.IncompleteParams.DocumentsRequired)}),
		person2.NewIsUnderAgeAnalyzer(person, &minimalAge),
		person2.NewDoaAnalyzer(fileService,
			doaService,
			person,
			ruleSetConfig.DOAParams.ApprovedScore,
			ruleSetConfig.DOAParams.RejectedScore),
		person2.NewHighRiskActivityAnalyzer(person, economicActivityService),
		person2.NewCafAnalyzer(person.EnrichedInformation),
		person2.NewMinimumIncomeAnalyzer(person),
	}

	received := factory.GetRules(ruleSetConfig, person)

	assert.Equal(t, expected, received)
}
