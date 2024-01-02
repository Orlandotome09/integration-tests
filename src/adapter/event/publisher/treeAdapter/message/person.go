package message

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message/types"
)

type Person struct {
	ID                     string                 `json:"id" binding:"required"`
	Type                   types.PersonType       `json:"type" binding:"required"`
	Name                   string                 `json:"name" binding:"required"`
	BirthDate              string                 `json:"birth_date" binding:"required"`  //ddMMyyyy
	Nationality            string                 `json:"nationality" binding:"required"` //BRA
	DigitalSign            bool                   `json:"digital_sign" binding:"required"`
	Addresses              Addresses              `json:"addresses,omitempty"`
	Accounts               Accounts               `json:"accounts,omitempty"`
	Contacts               Contacts               `json:"contacts,omitempty"`
	NotificationRecipients NotificationRecipients `json:"notification_recipients"`
	Individual             Individual             `json:"individual,omitempty"`
	LegalEntity            LegalEntity            `json:"legal_entity,omitempty"`
}
