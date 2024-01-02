package contract

import (
	"github.com/google/uuid"
	"time"
)

type Catalog struct {
	CatalogID       *uuid.UUID       `json:"catalog_id"`
	OfferType       string           `json:"offer_type"`
	RoleType        string           `json:"role_type"`
	PersonType      string           `json:"profile_type"`
	PartnerID       string           `json:"partner_id,omitempty"`
	ValidationSteps []ValidationStep `json:"validation_steps,omitempty"`
	ProductConfig   *ProductConfig   `json:"product_config,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type ValidationStep struct {
	StepNumber      int            `json:"step_number"`
	RulesConfig     *RuleSetConfig `json:"rules_config,omitempty"`
	SkipForApproval bool           `json:"skip_for_approval"`
}

type ProductConfig struct {
	CreateBexsAccount           bool `json:"create_bexs_account"`
	EnrichProfileWithBureauData bool `json:"enrich_profile_with_bureau_data"`
	TreeIntegration             bool `json:"tree_integration"`
	LimitIntegration            bool `json:"limit_integration"`
}

type RuleSetConfig struct {
	ManualBlockParams         *ManualBlockParams         `json:"manual_block_params,omitempty"`
	BlackListParams           *BlackListParams           `json:"black_list_params,omitempty"`
	BureauParams              *BureauParams              `json:"bureau_params,omitempty"`
	IncompleteParams          *IncompleteParams          `json:"incomplete_params,omitempty"`
	UnderAgeParams            *UnderAgeParams            `json:"under_age_params,omitempty"`
	WatchListParams           *WatchListParams           `json:"watch_list_params,omitempty"`
	DOAParams                 *DOAParams                 `json:"doa_params,omitempty"`
	OwnershipStructureParams  *OwnershipStructureParams  `json:"ownership_structure_params,omitempty"`
	PepParams                 *PepParams                 `json:"pep_params,omitempty"`
	LegalRepresentativeParams *LegalRepresentativeParams `json:"legal_representative_params,omitempty"`
	ActivityRiskParams        *ActivityRiskParams        `json:"activity_risk_params,omitempty"`
	BoardOfDirectorsParams    *BoardOfDirectorsParams    `json:"board_of_directors_params,omitempty"`
	ManualValidationParams    *ManualValidationParams    `json:"manual_validation_params,omitempty"`
	CAFParams                 *CAFParams                 `json:"caf_params,omitempty"`
}
