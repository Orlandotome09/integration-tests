package treeAdapterMessageTranslator

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"strconv"

	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message/types"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/accountTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/addressTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/companyTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/individualTranslator"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator/timeTranslator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TreeAdapterMessageTranslator interface {
	Translate(profile entity2.Profile, accounts []entity2.Account, addresses []entity2.Address) *message.TreeAdapterMessage
}

type treeAdapterMessageTranslator struct {
	timeTranslator       timeTranslator.TimeTranslator
	accountTranslator    accountTranslator.AccountTranslator
	addressTranslator    addressTranslator.AddressTranslator
	individualTranslator individualTranslator.IndividualTranslator
	companyTranslator    companyTranslator.CompanyTranslator
}

func New(timeTranslator timeTranslator.TimeTranslator,
	accountTranslator accountTranslator.AccountTranslator,
	addressTranslator addressTranslator.AddressTranslator,
	individualTranslator individualTranslator.IndividualTranslator,
	companyTranslator companyTranslator.CompanyTranslator,
) TreeAdapterMessageTranslator {
	return &treeAdapterMessageTranslator{
		timeTranslator:       timeTranslator,
		accountTranslator:    accountTranslator,
		addressTranslator:    addressTranslator,
		individualTranslator: individualTranslator,
		companyTranslator:    companyTranslator,
	}
}

func (ref *treeAdapterMessageTranslator) Translate(profile entity2.Profile, accounts []entity2.Account, addresses []entity2.Address) *message.TreeAdapterMessage {
	treeAdapterMessage := &message.TreeAdapterMessage{
		ProfileID:          translateProfileID(profile.ProfileID),
		Partner:            profile.PartnerID,
		ProfileCreatedDate: ref.timeTranslator.TranslateTime(profile.CreatedAt),
		Negotiation:        types.NegotiationTypeMesa,
		Person:             ref.translatePerson(profile),
		Date:               ref.timeTranslator.GenerateMessageDate(),
	}

	treeAdapterMessage.Person.Accounts = ref.accountTranslator.Translate(accounts)
	treeAdapterMessage.Person.Addresses = ref.addressTranslator.Translate(addresses)

	return treeAdapterMessage
}

func (ref *treeAdapterMessageTranslator) translatePerson(profile entity2.Profile) message.Person {
	switch profile.Person.PersonType {
	case values.PersonTypeIndividual:
		return ref.translateProfileIndividual(profile)
	case values.PersonTypeCompany:
		return ref.translateProfileCompany(profile)
	}
	return message.Person{}
}

func (ref *treeAdapterMessageTranslator) translateProfileIndividual(profile entity2.Profile) message.Person {
	return message.Person{
		ID:                     profile.DocumentNumber,
		Type:                   types.PersonTypeIndividual,
		Name:                   ref.individualTranslator.TranslateName(profile.Individual),
		BirthDate:              ref.individualTranslator.TranslateDateOfBirth(profile.Individual),
		Nationality:            ref.individualTranslator.TranslateNationality(profile.Individual),
		DigitalSign:            false,
		Addresses:              message.Addresses{},
		Contacts:               ref.translateContacts(profile),
		NotificationRecipients: ref.translateNotificationRecipients(profile.NotificationRecipients),
		Individual: message.Individual{
			Pep:    ref.individualTranslator.TranslatePep(profile.Individual),
			Income: ref.individualTranslator.TranslateIncome(profile.Individual),
			Assets: ref.individualTranslator.TranslateAssets(profile.Individual),
		},
	}
}

func (ref *treeAdapterMessageTranslator) translateContacts(profile entity2.Profile) message.Contacts {
	messageContacts := make([]message.Contact, len(profile.Contacts))

	maxPhoneNumbers := 2

	for i, contact := range profile.Contacts {
		messageContacts[i] = message.Contact{
			Name:  contact.Name,
			Email: contact.Email,
		}
		for j, phone := range contact.Phones {
			if j >= maxPhoneNumbers {
				break
			}

			number, err := strconv.ParseInt(phone.Number, 10, 64)
			if err != nil {
				logrus.Errorf("[TranslateContacts] Error converting contact. Profile ID: %v, Contact ID: %v, Phone: %v", profile.ProfileID, contact.ContactID, contact.Phone)
				continue
			}

			messageContacts[i].Phones[j] = message.Phone{
				Number: int(number),
			}
		}
	}

	return messageContacts
}

func (ref *treeAdapterMessageTranslator) translateNotificationRecipients(notificationRecipients []entity2.NotificationRecipient) message.NotificationRecipients {
	messageNotificationRecipients := make([]message.NotificationRecipient, len(notificationRecipients))
	for i, notificationRecipient := range notificationRecipients {
		messageNotificationRecipients[i] = message.NotificationRecipient{
			NotificationType: notificationRecipient.NotificationType,
			EmailTo:          notificationRecipient.EmailTo,
			CopyEmail:        notificationRecipient.CopyEmail,
		}
	}
	return messageNotificationRecipients
}

func (ref *treeAdapterMessageTranslator) translateProfileCompany(profile entity2.Profile) message.Person {
	return message.Person{
		ID:                     profile.DocumentNumber,
		Type:                   types.PersonTypeLegalEntity,
		Name:                   ref.companyTranslator.TranslateLegalName(profile.Company),
		BirthDate:              ref.companyTranslator.TranslateDateOfIncorporation(profile.Company),
		Nationality:            ref.companyTranslator.TranslatePlaceOfIncorporation(profile.Company),
		DigitalSign:            false,
		Addresses:              message.Addresses{},
		Contacts:               ref.translateContacts(profile),
		NotificationRecipients: ref.translateNotificationRecipients(profile.NotificationRecipients),
		LegalEntity: message.LegalEntity{
			Size:             types.CompanySizeNotInformed,
			BusinessActivity: types.BusinessActivityA,
			MonthlyIncome:    ref.companyTranslator.TranslateIncome(profile.Company),
		},
	}
}

func translateProfileID(profileID *uuid.UUID) string {
	if profileID == nil {
		return ""
	}
	return profileID.String()
}
