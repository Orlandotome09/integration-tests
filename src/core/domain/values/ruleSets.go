package values

import (
	"fmt"
)

type RuleSet string

const (
	RuleSetIsUnderAge           RuleSet = "IS_UNDER_AGE"
	RuleSetSerasa               RuleSet = "SERASA_BUREAU"
	RuleSetWatchlist            RuleSet = "WATCHLIST"
	RuleSetBlacklist            RuleSet = "BLACKLIST"
	RuleSetState                RuleSet = "STATE"
	RuleSetIncomplete           RuleSet = "INCOMPLETE"
	RuleSetIncompleteContract   RuleSet = "INCOMPLETE_CONTRACT"
	RuleSetPreconditionContract RuleSet = "PRECONDITION_CONTRACT_SET"
	RuleSetTest                 RuleSet = "TEST"
	RuleSetDOA                  RuleSet = "DOA"
	RuleSetManualBlock          RuleSet = "MANUAL_BLOCK"
	RuleSetOwnershipStructure   RuleSet = "OWNERSHIP_STRUCTURE"
	RuleSetPep                  RuleSet = "PEP"
	RuleSetLegalRepresentatives RuleSet = "LEGAL_REPRESENTATIVES"
	RuleSetActivityRisk         RuleSet = "ACTIVITY_RISK"
	RuleSetBoardOfDirectors     RuleSet = "BOARD_OF_DIRECTORS"
	RuleSetManualValidation     RuleSet = "MANUAL_VALIDATION"
	RuleSetCafAnalysis          RuleSet = "CAF_ANALYSIS"
	RuleSetMinimumBilling       RuleSet = "MINIMUM_BILLING"
	RuleSetMinimumIncome        RuleSet = "MINIMUM_INCOME"
)

var validRuleSets = map[string]RuleSet{
	RuleSetIsUnderAge.ToString():           RuleSetIsUnderAge,
	RuleSetSerasa.ToString():               RuleSetSerasa,
	RuleSetWatchlist.ToString():            RuleSetWatchlist,
	RuleSetBlacklist.ToString():            RuleSetBlacklist,
	RuleSetState.ToString():                RuleSetState,
	RuleSetIncomplete.ToString():           RuleSetIncomplete,
	RuleSetIncompleteContract.ToString():   RuleSetIncompleteContract,
	RuleSetPreconditionContract.ToString(): RuleSetPreconditionContract,
	RuleSetTest.ToString():                 RuleSetTest,
	RuleSetDOA.ToString():                  RuleSetDOA,
	RuleSetManualBlock.ToString():          RuleSetManualBlock,
	RuleSetOwnershipStructure.ToString():   RuleSetOwnershipStructure,
	RuleSetPep.ToString():                  RuleSetPep,
	RuleSetLegalRepresentatives.ToString(): RuleSetLegalRepresentatives,
	RuleSetActivityRisk.ToString():         RuleSetActivityRisk,
	RuleSetBoardOfDirectors.ToString():     RuleSetBoardOfDirectors,
	RuleSetManualValidation.ToString():     RuleSetManualValidation,
	RuleSetMinimumBilling.ToString():       RuleSetMinimumBilling,
	RuleSetMinimumIncome.ToString():        RuleSetMinimumIncome,
}

func (ruleSet RuleSet) Validate() error {
	_, in := validRuleSets[ruleSet.ToString()]
	if !in {
		return NewErrorValidation(fmt.Sprintf("%s is an invalid rule set name", ruleSet))
	}
	return nil
}

func (ruleSet RuleSet) ToString() string {
	return string(ruleSet)
}
