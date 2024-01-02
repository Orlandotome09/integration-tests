package eventProcessor

import (
	"time"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type eventProcessor struct {
	complianceAnalyzer interfaces.ComplianceAnalyzer
	engineFactory      interfaces.EngineFactory
}

func New(
	complianceAnalyzer interfaces.ComplianceAnalyzer,
	engineFactory interfaces.EngineFactory,
) interfaces.EventProcessor {
	return &eventProcessor{
		complianceAnalyzer: complianceAnalyzer,
		engineFactory:      engineFactory,
	}
}

func (ref *eventProcessor) ExecuteForEvent(event *values.Event) (*entity.State, error) {
	if event.ParentID == uuid.Nil || event.EngineName == "" {
		return nil, errors.New("ParentID or EngineName are empty")
	}

	logrus.
		WithField("engine_name", event.EngineName).
		WithField("entity_id", event.ParentID).
		WithField("date", event.Date).
		WithField("request_data", event.RequestDate).
		Info("Running Compliance Analysis for Event")

	engine, err := ref.engineFactory.CreateEngine(event.EngineName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = engine.Prepare(event.ParentID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	now := time.Now()
	executionTime := now

	state, err := ref.complianceAnalyzer.RunComplianceAnalysis(engine, event.ParentID, event.RequestDate, executionTime, false)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return state, nil
}
