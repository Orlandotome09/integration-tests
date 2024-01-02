package treeAdapterEventTopicPublisher

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"encoding/json"

	treeAdapterMessageTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type treeAdapterPublisher struct {
	queuePublisher interfaces.QueuePublisher
	translator     treeAdapterMessageTranslator.TreeAdapterMessageTranslator
}

func New(queuePublisher interfaces.QueuePublisher,
	translator treeAdapterMessageTranslator.TreeAdapterMessageTranslator,
) interfaces.TreeAdapterPublisher {
	return &treeAdapterPublisher{
		queuePublisher: queuePublisher,
		translator:     translator,
	}
}

func (ref *treeAdapterPublisher) Send(profile entity2.Profile, accounts []entity2.Account, addresses []entity2.Address) error {
	treeAdapterMessage := ref.translator.Translate(profile, accounts, addresses)

	message, err := json.Marshal(treeAdapterMessage)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = ref.queuePublisher.Publish(message, "", ""); err != nil {
		return errors.WithStack(err)
	}
	logrus.
		WithField("message_data", string(message)).
		Infof("[TreeAdapterTopicPublisher] Published Event To Tree Adapter Topic")
	return nil
}
