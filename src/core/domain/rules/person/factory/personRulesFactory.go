package personRulesFactory

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	person2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/person"
	"github.com/sirupsen/logrus"
)

type personRulesFactory struct {
	fileService               interfaces.FileAdapter
	doaService                interfaces.DOAService
	economicalActivityService interfaces.EconomicActivityService
}

func New(
	fileService interfaces.FileAdapter,
	doaService interfaces.DOAService,
	economicalActivityService interfaces.EconomicActivityService,
) interfaces.PersonRulesFactory {

	return &personRulesFactory{
		fileService:               fileService,
		doaService:                doaService,
		economicalActivityService: economicalActivityService,
	}
}

func (ref *personRulesFactory) GetRules(ruleSetConfig entity2.RuleSetConfig, person entity2.Person) []entity2.Rule {
	var resultRules []entity2.Rule

	if ruleSetConfig.ManualBlockParams != nil {
		resultRules = append(resultRules, person2.NewManualBlockAnalyzer(person))
	}
	if ruleSetConfig.BlackListParams != nil {
		resultRules = append(resultRules, person2.NewBlackListAnalyzer(person))
	}
	if ruleSetConfig.WatchListParams != nil {
		resultRules = append(resultRules, person2.NewWatchlistAnalyzer(person,
			ruleSetConfig.WatchListParams.WantPepTag,
			ruleSetConfig.WatchListParams.WantedSources,
			ruleSetConfig.WatchListParams.HasMatchInWatchListStatus))
	}
	if ruleSetConfig.PepParams != nil {
		resultRules = append(resultRules, person2.NewPepAnalyzer(person))
	}
	if ruleSetConfig.BureauParams != nil {
		resultRules = append(resultRules, person2.NewBureauAnalyzer(person,
			ruleSetConfig.BureauParams.ApprovedStatuses,
			ruleSetConfig.BureauParams.NotFoundInSerasaStatus,
			ruleSetConfig.BureauParams.NotFoundInSerasaPending,
			ruleSetConfig.BureauParams.HasProblemsInSerasaStatus,
			ruleSetConfig.BureauParams.HasProblemsInSerasaPending))
	}

	if ruleSetConfig.IncompleteParams != nil {
		var analyzers []person2.Analyzer

		if ruleSetConfig.HasIncompleteFieldsValidation() {
			analyzers = append(analyzers, person2.NewIncompleteFieldsAnalyzer(
				ruleSetConfig.IncompleteParams.DateOfBirthRequired,
				ruleSetConfig.IncompleteParams.InputtedOrEnrichedDateOfBirthRequired,
				ruleSetConfig.IncompleteParams.PhoneNumberRequired,
				ruleSetConfig.IncompleteParams.EmailRequired,
				ruleSetConfig.IncompleteParams.PepRequired,
				ruleSetConfig.IncompleteParams.LastNameRequired,
			))
		}
		if ruleSetConfig.IncompleteParams.AddressRequired {
			analyzers = append(analyzers, person2.NewIncompleteAddressAnalyzer())
		}
		if len(ruleSetConfig.IncompleteParams.DocumentsRequired) > 0 {
			analyzers = append(analyzers, person2.NewIncompleteDocumentsAnalyzer(ruleSetConfig.IncompleteParams.DocumentsRequired))
		}
		resultRules = append(resultRules, person2.NewIncompleteAnalyzer(person, analyzers))
	}
	if ruleSetConfig.UnderAgeParams != nil {
		resultRules = append(resultRules, person2.NewIsUnderAgeAnalyzer(person, ruleSetConfig.UnderAgeParams.MinimumAge))
	}
	if ruleSetConfig.DOAParams != nil {
		resultRules = append(resultRules,
			person2.NewDoaAnalyzer(
				ref.fileService,
				ref.doaService,
				person,
				ruleSetConfig.DOAParams.ApprovedScore,
				ruleSetConfig.DOAParams.RejectedScore))
	}
	if ruleSetConfig.ActivityRiskParams != nil {
		resultRules = append(resultRules, person2.NewHighRiskActivityAnalyzer(person,
			ref.economicalActivityService))
	}
	if ruleSetConfig.CAFParams != nil {
		resultRules = append(resultRules, person2.NewCafAnalyzer(
			person.EnrichedInformation))
	}
	if ruleSetConfig.MinimumIncomeParams != nil {
		resultRules = append(resultRules, person2.NewMinimumIncomeAnalyzer(person))
	}

	logrus.WithField("ruleSetConfig", ruleSetConfig).
		WithField("rules", resultRules).
		Infof("[personRulesFactory] Rules created for person %v", person.EntityID)

	return resultRules
}
