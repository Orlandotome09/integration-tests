package limitEventTopicPublisher

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/limit/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type limitPublisher struct {
	queuePublisher interfaces.QueuePublisher
	translator     limitMessageTranslator.LimitMessageTranslator
}

func New(queuePublisher interfaces.QueuePublisher,
	translator limitMessageTranslator.LimitMessageTranslator) interfaces.LimitPublisher {
	return &limitPublisher{
		queuePublisher: queuePublisher,
		translator:     translator,
	}
}

func (ref *limitPublisher) Send(profile entity2.Profile, state entity2.State) error {
	limitMessage := ref.translator.Translate(profile, state)

	message, err := json.Marshal(limitMessage)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = ref.queuePublisher.Publish(message, "", ""); err != nil {
		return errors.WithStack(err)
	}
	logrus.
		WithField("message_data", string(message)).
		Infof("[LimitTopicPublisher] Published Event To LimitTopic")
	return nil
}
