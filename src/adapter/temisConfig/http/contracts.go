package http

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/pkg/errors"
)

type CadastralValidationConfigResponse []CadastralValidationConfigBase

type CadastralValidationConfigBase struct {
	PersonType      string           `json:"person_type"`
	OfferType       string           `json:"offer_type" binding:"required"`
	RoleType        string           `json:"role_type" binding:"required"`
	PartnerID       string           `json:"partner_id,omitempty"`
	ValidationSteps []ValidationStep `json:"validation_steps,omitempty"`
	ProductConfig   *ProductConfig   `json:"product_config"`
}

type ProductConfig struct {
	CreateBexsAccount           bool `json:"create_bexs_account"`
	EnrichProfileWithBureauData bool `json:"enrich_profile_with_bureau_data"`
	TreeIntegration             bool `json:"tree_integration"`
	LimitIntegration            bool `json:"limit_integration"`
}

type ValidationStep struct {
	StepNumber      int                `json:"step_number"`
	RulesConfig     *RuleSetConfigBase `json:"rules_config"`
	SkipForApproval bool               `json:"skip_for_approval"`
}

type RuleSetConfigBase struct {
	ManualBlockParams         *manualBlockParams         `json:"manual_block,omitempty"`
	BlackListParams           *blackListParams           `json:"blacklist,omitempty"`
	BureauParams              *bureauParams              `json:"bureau,omitempty"`
	IncompleteParams          *incompleteParams          `json:"incomplete,omitempty"`
	UnderAgeParams            *underAgeParams            `json:"under_age,omitempty"`
	WatchListParams           *watchListParams           `json:"watchlist,omitempty"`
	DOAParams                 *DOAParamsStruct           `json:"doa,omitempty"`
	OwnershipStructureParams  *ownershipStructureParams  `json:"ownership_structure,omitempty"`
	PepParams                 *pepParams                 `json:"pep,omitempty"`
	LegalRepresentativeParams *legalRepresentativeParams `json:"legal_representative,omitempty"`
	ActivityRiskParams        *activityRiskParams        `json:"activity_risk,omitempty"`
	BoardOfDirectorsParams    *boardOfDirectorsParams    `json:"board_of_directors,omitempty"`
	ManualValidationParams    *manualValidationParams    `json:"manual_validation,omitempty"`
	CAFParams                 *CAFParams                 `json:"caf,omitempty"`
	MinimumBillingParams      *minimumBillingParams      `json:"minimum_billing,omitempty"`
	MinimumIncomeParams       *minimumIncomeParams       `json:"minimum_income,omitempty"`
	BorderTownParams          *borderTownParams          `json:"border_town,omitempty"`
}

type manualBlockParams struct{}

type blackListParams struct{}

type manualValidationParams struct{}

type borderTownParams struct{}

type CAFParams struct{}

type bureauParams struct {
	NotFoundInSerasaStatus     *values.Result `json:"not_found_in_serasa_status,omitempty"`
	NotFoundInSerasaPending    *bool          `json:"not_found_in_serasa_pending,omitempty"`
	HasProblemsInSerasaStatus  *values.Result `json:"has_problems_in_serasa_status,omitempty"`
	HasProblemsInSerasaPending *bool          `json:"has_problems_in_serasa_pending,omitempty"`
	ApprovedStatuses           []string       `json:"approved_statuses,omitempty"`
}

type incompleteParams struct {
	DateOfBirthRequired                   bool               `json:"date_of_birth_required,omitempty"`
	InputtedOrEnrichedDateOfBirthRequired bool               `json:"inputted_or_enriched_date_of_birth_required,omitempty"`
	PhoneNumberRequired                   bool               `json:"phone_number_required,omitempty"`
	AddressRequired                       bool               `json:"address_required,omitempty"`
	EmailRequired                         bool               `json:"email_required,omitempty"`
	DocumentsRequired                     []DocumentRequired `json:"documents_required,omitempty"`
	PepRequired                           bool               `json:"pep_required,omitempty"`
	LastNameRequired                      bool               `json:"last_name_required,omitempty"`
}

type DocumentRequired struct {
	DocumentType      string      `json:"document_type,omitempty"`
	DocumentSubType   string      `json:"document_sub_type,omitempty"`
	FileRequired      bool        `json:"file_required,omitempty"`
	PendingOnApproval bool        `json:"pending_on_approval,omitempty"`
	Conditions        []Condition `json:"conditions,omitempty"`
}

type Condition struct {
	FieldName string   `json:"field_name"`
	Values    []string `json:"values"`
}

type underAgeParams struct {
	MinimumAge *int `json:"minimum_age,omitempty"`
}

type watchListParams struct {
	WantPepTag                bool           `json:"want_pep_tag,omitempty"`
	WantedSources             []string       `json:"wanted_sources,omitempty"`
	HasMatchInWatchListStatus *values.Result `json:"has_match_in_watch_list_status,omitempty"`
}

type DOAParamsStruct struct {
	ApprovedScore *float64 `json:"approved_score,omitempty"`
	RejectedScore *float64 `json:"rejected_score,omitempty"`
}

type ownershipStructureParams struct {
	MinShareholdingPercentage *float64 `json:"min_shareholding_percentage,omitempty"`
}

type pepParams struct{}

type legalRepresentativeParams struct{}

type activityRiskParams struct{}

type boardOfDirectorsParams struct{}

type minimumBillingParams struct {
	MinimumBilling *float64 `json:"minimum_billing,omitempty"`
}

type minimumIncomeParams struct {
	MinimumIncome *float64 `json:"minimum_income,omitempty"`
}

func (responses CadastralValidationConfigResponse) ToDomain() entity.CadastralValidationConfigMap {
	configs := entity.CadastralValidationConfigMap{}
	for _, response := range responses {
		key := response.PersonType + response.RoleType + response.OfferType + response.PartnerID
		config, err := response.ToDomain()
		if err != nil {
			return nil
		}
		configs[key] = config
	}

	return configs
}

func (ref *CadastralValidationConfigBase) ToDomain() (entity.CadastralValidationConfig, error) {
	roleType, err := values.ParseToRoleType(ref.RoleType)
	if err != nil {
		return entity.CadastralValidationConfig{}, errors.WithStack(err)
	}
	personType, err := values.ParseToPersonType(ref.PersonType)
	if err != nil {
		return entity.CadastralValidationConfig{}, errors.WithStack(err)
	}

	validationSteps := make([]entity.ValidationStep, 0)
	if len(ref.ValidationSteps) > 0 {
		for _, step := range ref.ValidationSteps {
			validationStep, err := step.ToDomain()
			if err != nil {
				return entity.CadastralValidationConfig{}, errors.WithStack(err)
			}
			validationSteps = append(validationSteps, *validationStep)
		}
	}

	catalog := entity.CadastralValidationConfig{
		OfferType:       ref.OfferType,
		RoleType:        roleType,
		PersonType:      personType,
		PartnerID:       ref.PartnerID,
		ProductConfig:   ref.ProductConfig.ToDomain(),
		ValidationSteps: validationSteps,
	}
	return catalog, nil
}

func (ref *ProductConfig) ToDomain() *entity.ProductConfig {
	return &entity.ProductConfig{
		CreateBexsAccount:           ref.CreateBexsAccount,
		EnrichProfileWithBureauData: ref.EnrichProfileWithBureauData,
		TreeIntegration:             ref.TreeIntegration,
		LimitIntegration:            ref.LimitIntegration,
	}
}

func (ref *ValidationStep) ToDomain() (*entity.ValidationStep, error) {
	var rules *entity.RuleSetConfig
	if ref.RulesConfig != nil {
		result, err := ref.RulesConfig.ToDomain()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		rules = result
	}

	return &entity.ValidationStep{
		StepNumber:      ref.StepNumber,
		RulesConfig:     rules,
		SkipForApproval: ref.SkipForApproval,
	}, nil
}

func (ref *manualBlockParams) toDomain() *entity.ManualBlockParams {
	if ref == nil {
		return nil
	}
	return &entity.ManualBlockParams{}
}

func (ref *blackListParams) toDomain() *entity.BlackListParams {
	if ref == nil {
		return nil
	}
	return &entity.BlackListParams{}
}

func (ref *bureauParams) toDomain() *entity.BureauParams {
	if ref == nil {
		return nil
	}

	return &entity.BureauParams{
		NotFoundInSerasaStatus:     ref.NotFoundInSerasaStatus,
		NotFoundInSerasaPending:    ref.NotFoundInSerasaPending,
		HasProblemsInSerasaStatus:  ref.HasProblemsInSerasaStatus,
		HasProblemsInSerasaPending: ref.HasProblemsInSerasaPending,
		ApprovedStatuses:           ref.ApprovedStatuses,
	}
}

func (ref *incompleteParams) toDomain() (*entity.IncompleteParams, error) {
	if ref == nil {
		return nil, nil
	}
	var requiredDocuments []entity.DocumentRequired
	for _, required := range ref.DocumentsRequired {
		requiredDocument, err := required.toDomain()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		requiredDocuments = append(requiredDocuments, *requiredDocument)
	}
	return &entity.IncompleteParams{
		DateOfBirthRequired:                   ref.DateOfBirthRequired,
		InputtedOrEnrichedDateOfBirthRequired: ref.InputtedOrEnrichedDateOfBirthRequired,
		PhoneNumberRequired:                   ref.PhoneNumberRequired,
		AddressRequired:                       ref.AddressRequired,
		DocumentsRequired:                     requiredDocuments,
		EmailRequired:                         ref.EmailRequired,
		PepRequired:                           ref.PepRequired,
		LastNameRequired:                      ref.LastNameRequired,
	}, nil
}

func (ref DocumentRequired) toDomain() (*entity.DocumentRequired, error) {
	documentType, err := values.ParseToDocumentType(ref.DocumentType)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var conditions []entity.Condition
	for _, condition := range ref.Conditions {
		documentCondition := condition.toDomain()
		conditions = append(conditions, *documentCondition)
	}

	return &entity.DocumentRequired{
		DocumentType:      documentType,
		DocumentSubType:   ref.DocumentSubType,
		FileRequired:      ref.FileRequired,
		PendingOnApproval: ref.PendingOnApproval,
		Conditions:        conditions,
	}, nil
}

func (ref Condition) toDomain() *entity.Condition {
	return &entity.Condition{
		FieldName: ref.FieldName,
		Values:    ref.Values,
	}
}

func (ref *underAgeParams) toDomain() *entity.UnderAgeParams {
	if ref == nil {
		return nil
	}
	return &entity.UnderAgeParams{MinimumAge: ref.MinimumAge}
}

func (ref *watchListParams) toDomain() *entity.WatchListParams {
	if ref == nil {
		return nil
	}
	return &entity.WatchListParams{
		WantPepTag:                ref.WantPepTag,
		WantedSources:             ref.WantedSources,
		HasMatchInWatchListStatus: ref.HasMatchInWatchListStatus,
	}
}

func (ref *DOAParamsStruct) toDomain() *entity.DOAParams {
	if ref == nil {
		return nil
	}
	return &entity.DOAParams{
		ApprovedScore: ref.ApprovedScore,
		RejectedScore: ref.RejectedScore,
	}
}

func (ref *ownershipStructureParams) toDomain() *entity.OwnershipStructureParams {
	if ref == nil {
		return nil
	}
	return &entity.OwnershipStructureParams{
		MinShareholdingPercentage: ref.MinShareholdingPercentage,
	}
}

func (ref *pepParams) toDomain() *entity.PepParams {
	if ref == nil {
		return nil
	}
	return &entity.PepParams{}
}

func (ref *legalRepresentativeParams) toDomain() *entity.LegalRepresentativeParams {
	if ref == nil {
		return nil
	}
	return &entity.LegalRepresentativeParams{}
}

func (ref *activityRiskParams) toDomain() *entity.ActivityRiskParams {
	if ref == nil {
		return nil
	}
	return &entity.ActivityRiskParams{}
}

func (ref *boardOfDirectorsParams) toDomain() *entity.BoardOfDirectorsParams {
	if ref == nil {
		return nil
	}
	return &entity.BoardOfDirectorsParams{}
}

func (ref *manualValidationParams) toDomain() *entity.ManualValidationParams {
	if ref == nil {
		return nil
	}
	return &entity.ManualValidationParams{}
}

func (ref *CAFParams) toDomain() *entity.CAFParams {
	if ref == nil {
		return nil
	}
	return &entity.CAFParams{}
}

func (ref *minimumBillingParams) toDomain() *entity.MinimumBillingParams {
	if ref == nil {
		return nil
	}
	return &entity.MinimumBillingParams{
		MinimumBilling: ref.MinimumBilling,
	}
}

func (ref *minimumIncomeParams) toDomain() *entity.MinimumIncomeParams {
	if ref == nil {
		return nil
	}
	return &entity.MinimumIncomeParams{
		MinimumIncome: ref.MinimumIncome,
	}
}

func (ref *RuleSetConfigBase) ToDomain() (*entity.RuleSetConfig, error) {
	incompleteParams, err := ref.IncompleteParams.toDomain()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ruleSetConfig := &entity.RuleSetConfig{
		ManualBlockParams:         ref.ManualBlockParams.toDomain(),
		BlackListParams:           ref.BlackListParams.toDomain(),
		BureauParams:              ref.BureauParams.toDomain(),
		IncompleteParams:          incompleteParams,
		UnderAgeParams:            ref.UnderAgeParams.toDomain(),
		WatchListParams:           ref.WatchListParams.toDomain(),
		DOAParams:                 ref.DOAParams.toDomain(),
		OwnershipStructureParams:  ref.OwnershipStructureParams.toDomain(),
		PepParams:                 ref.PepParams.toDomain(),
		LegalRepresentativeParams: ref.LegalRepresentativeParams.toDomain(),
		ActivityRiskParams:        ref.ActivityRiskParams.toDomain(),
		BoardOfDirectorsParams:    ref.BoardOfDirectorsParams.toDomain(),
		ManualValidationParams:    ref.ManualValidationParams.toDomain(),
		CAFParams:                 ref.CAFParams.toDomain(),
		MinimumBillingParams:      ref.MinimumBillingParams.toDomain(),
		MinimumIncomeParams:       ref.MinimumIncomeParams.toDomain(),
	}
	if err := ruleSetConfig.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}
	return ruleSetConfig, nil
}
