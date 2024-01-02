package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/pkg/errors"
)

type ManualBlockParams struct{}

type BlackListParams struct{}

type BureauParams struct {
	NotFoundInSerasaStatus     *values.Result `json:"not_found_in_serasa_status,omitempty"`
	NotFoundInSerasaPending    *bool          `json:"not_found_in_serasa_pending,omitempty"`
	HasProblemsInSerasaStatus  *values.Result `json:"has_problems_in_serasa_status,omitempty"`
	HasProblemsInSerasaPending *bool          `json:"has_problems_in_serasa_pending,omitempty"`
	ApprovedStatuses           []string       `json:"approved_statuses,omitempty"`
}

type IncompleteParams struct {
	DateOfBirthRequired                   bool              `json:"date_of_birth_required"`
	InputtedOrEnrichedDateOfBirthRequired bool              `json:"inputted_or_enriched_date_of_birth_required"`
	PhoneNumberRequired                   bool              `json:"phone_number_required"`
	LastNameRequired                      bool              `json:"last_name_required"`
	AddressRequired                       bool              `json:"address_required"`
	EmailRequired                         bool              `json:"email_required"`
	DocumentsRequired                     DocumentsRequired `json:"documents_required,omitempty"`
	PepRequired                           bool              `json:"pep_required"`
}

type Condition struct {
	FieldName string   `json:"field_name"`
	Values    []string `json:"values"`
}

type UnderAgeParams struct {
	MinimumAge *int `json:"minimum_age"`
}

type WatchListParams struct {
	WantPepTag                bool           `json:"want_pep_tag"`
	WantedSources             []string       `json:"wanted_sources,omitempty"`
	HasMatchInWatchListStatus *values.Result `json:"has_match_in_watch_list_status,omitempty"`
}

type PepParams struct{}

type DOAParams struct {
	ApprovedScore *float64 `json:"approved_score,omitempty"`
	RejectedScore *float64 `json:"rejected_score,omitempty"`
}

type OwnershipStructureParams struct {
	MinShareholdingPercentage *float64 `json:"min_shareholding_percentage,omitempty"`
}

func (params *OwnershipStructureParams) Validate() error {
	if params.MinShareholdingPercentage != nil {
		if *params.MinShareholdingPercentage > 100.0 {
			return errors.New("min shareholding percentage should be max 100.0")
		}
	}
	return nil
}

type LegalRepresentativeParams struct{}

type ActivityRiskParams struct{}

type BoardOfDirectorsParams struct{}

type ManualValidationParams struct{}

type CAFParams struct{}

type MinimumBillingParams struct {
	MinimumBilling *float64 `json:"minimum_billing,omitempty"`
}

type MinimumIncomeParams struct {
	MinimumIncome *float64 `json:"minimum_income,omitempty"`
}
