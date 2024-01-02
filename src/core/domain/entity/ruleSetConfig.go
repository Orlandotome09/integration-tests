package entity

type RuleSetConfig struct {
	ManualBlockParams         *ManualBlockParams         `json:"manual_block,omitempty"`
	BlackListParams           *BlackListParams           `json:"blacklist,omitempty"`
	BureauParams              *BureauParams              `json:"bureau,omitempty"`
	IncompleteParams          *IncompleteParams          `json:"incomplete,omitempty"`
	UnderAgeParams            *UnderAgeParams            `json:"under_age,omitempty"`
	WatchListParams           *WatchListParams           `json:"watchlist,omitempty"`
	DOAParams                 *DOAParams                 `json:"doa,omitempty"`
	OwnershipStructureParams  *OwnershipStructureParams  `json:"ownership_structure,omitempty"`
	PepParams                 *PepParams                 `json:"pep,omitempty"`
	LegalRepresentativeParams *LegalRepresentativeParams `json:"legal_representative,omitempty"`
	ActivityRiskParams        *ActivityRiskParams        `json:"activity_risk,omitempty"`
	BoardOfDirectorsParams    *BoardOfDirectorsParams    `json:"board_of_directors,omitempty"`
	ManualValidationParams    *ManualValidationParams    `json:"manual_validation,omitempty"`
	CAFParams                 *CAFParams                 `json:"caf,omitempty"`
	MinimumBillingParams      *MinimumBillingParams      `json:"minimum_billing,omitempty"`
	MinimumIncomeParams       *MinimumIncomeParams       `json:"minimum_income,omitempty"`
}

func (ruleSetConfig *RuleSetConfig) Validate() error {
	if ruleSetConfig.OwnershipStructureParams != nil {
		if err := ruleSetConfig.OwnershipStructureParams.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (ruleSetConfig RuleSetConfig) HasIncompleteFieldsValidation() bool {
	return ruleSetConfig.IncompleteParams != nil &&
		(ruleSetConfig.IncompleteParams.LastNameRequired ||
			ruleSetConfig.IncompleteParams.DateOfBirthRequired ||
			ruleSetConfig.IncompleteParams.EmailRequired ||
			ruleSetConfig.IncompleteParams.InputtedOrEnrichedDateOfBirthRequired ||
			ruleSetConfig.IncompleteParams.PhoneNumberRequired ||
			ruleSetConfig.IncompleteParams.PepRequired)
}
