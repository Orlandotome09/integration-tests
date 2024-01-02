package limitEventTopicPublisher

import (
	message "bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/limit/message"
	limitMessageTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/limit/translator"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSend(t *testing.T) {
	queuePublisher := &mocks.QueuePublisher{}
	translator := &limitMessageTranslator.MockLimitMessageTranslator{}
	publisher := New(queuePublisher, translator)

	profile := entity2.Profile{}
	state := entity2.State{}

	limitMessage := &message.LimitMessage{}
	msg, _ := json.Marshal(limitMessage)

	translator.On("Translate", profile).Return(limitMessage)

	queuePublisher.On("Publish", msg, "", "").Return(nil)

	err := publisher.Send(profile, state)

	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, translator, queuePublisher)
}
