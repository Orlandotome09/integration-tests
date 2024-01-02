package complianceCommandPublisher

import (
	"encoding/json"
	"time"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type complianceTopicPublisher struct {
	queuePublisher interfaces.QueuePublisher
}

func New(queuePublisher interfaces.QueuePublisher) interfaces.ComplianceCommandPublisher {
	return &complianceTopicPublisher{
		queuePublisher: queuePublisher,
	}
}

func (ref *complianceTopicPublisher) SendCommand(entityID uuid.UUID, parentID *uuid.UUID, engineName values.EngineName) (*entity.State, error) {
	event := &complianceCommand{
		EntityID:   entityID,
		ParentID:   parentID,
		EngineName: engineName,
		Date:       time.Now(),
	}

	message, err := json.Marshal(event)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err = ref.queuePublisher.Publish(message, "", entityID.String()); err != nil {
		return nil, errors.WithStack(err)
	}
	logrus.
		WithField("message_data", string(message)).
		Infof("[ComplianceTopicPublisher] Published command To Compliance Topic")
	return nil, nil
}

type complianceCommand struct {
	EngineName string      `json:"engine_name"`
	EventType  string      `json:"event_type"`
	ParentID   *uuid.UUID  `json:"parent_id"`
	EntityID   uuid.UUID   `json:"entity_id"`
	Date       time.Time   `json:"date"`
	Content    interface{} `json:"content"`
}
