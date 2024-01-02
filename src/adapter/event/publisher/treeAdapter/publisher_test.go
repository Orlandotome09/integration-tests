package treeAdapterEventTopicPublisher

import (
	message "bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message"
	treeAdapterMessageTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSend(t *testing.T) {
	queuePublisher := &mocks.QueuePublisher{}
	translator := &treeAdapterMessageTranslator.MockTreeAdapterMessageTranslator{}
	publisher := New(queuePublisher, translator)

	profile := entity2.Profile{}
	accounts := []entity2.Account{{AccountNumber: "xxx"}, {AccountNumber: "yyy"}}
	addresses := []entity2.Address{}

	treeAdapterMessage := &message.TreeAdapterMessage{}
	msg, _ := json.Marshal(treeAdapterMessage)

	translator.On("Translate", profile, accounts, addresses).Return(treeAdapterMessage)
	queuePublisher.On("Publish", msg, "", "").Return(nil)

	err := publisher.Send(profile, accounts, addresses)

	if err != nil {
		t.Errorf("\nExpected error nil\n")
	}

	mock.AssertExpectationsForObjects(t, translator, queuePublisher)
}
