package values

type EventType = string

const (
	EventTypeProfileCreated EventType = "PROFILE_CREATED"
	EventTypeProfileChanged EventType = "PROFILE_CHANGED"
	EventTypeProfileResync  EventType = "PROFILE_RESYNC"

	EventTypeAddressCreated EventType = "ADDRESS_CREATED"
	EventTypeAddressChanged EventType = "ADDRESS_CHANGED"
	EventTypeAddressDeleted EventType = "ADDRESS_DELETED"

	EventTypeDocumentCreated EventType = "DOCUMENT_CREATED"
	EventTypeDocumentChanged EventType = "DOCUMENT_CHANGED"

	EventTypeDocumentFileCreated EventType = "DOCUMENT_FILE_CREATED"

	EventTypeDoaResultChanged EventType = "DOA_RESULT_CHANGED"

	EventTypeForeignAccountCreated EventType = "FOREIGN_ACCOUNT_CREATED"
	EventTypeForeignAccountChanged EventType = "FOREIGN_ACCOUNT_CHANGED"

	EventTypeAccountCreated EventType = "ACCOUNT_CREATED"
	EventTypeAccountChanged EventType = "ACCOUNT_CHANGED"

	EventTypeOverrideCreated EventType = "OVERRIDE_CREATED"
	EventTypeOverrideDeleted EventType = "OVERRIDE_DELETED"

	EventTypeStateCreated EventType = "STATE_CREATED"
	EventTypeStateChanged EventType = "STATE_CHANGED"
	EventTypeStateResync  EventType = "STATE_RESYNC"

	EventTypeContractCreated EventType = "CONTRACT_CREATED"
	EventTypeContractChanged EventType = "CONTRACT_CHANGED"

	EventTypeLegalRepresentativeCreated EventType = "LEGAL_REPRESENTATIVE_CREATED"
	EventTypeLegalRepresentativeChanged EventType = "LEGAL_REPRESENTATIVE_CHANGED"

	EventTypeQuestionFormCreated EventType = "QUESTION_FORM_CREATED"
	EventTypeQuestionFormChanged EventType = "QUESTION_FORM_CHANGED"

	EventTypeShareholderCreated EventType = "SHAREHOLDER_CREATED"
	EventTypeShareholderChanged EventType = "SHAREHOLDER_CHANGED"
	EventTypeShareholderDeleted EventType = "SHAREHOLDER_DELETED"

	EventTypePersonEnriched EventType = "PERSON_ENRICHED"
)
