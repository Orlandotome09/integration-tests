package enginefactory

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/pkg/errors"
)

type engineFactory struct {
	profileEngine  _interfaces.Engine
	contractEngine _interfaces.Engine
}

func NewEngineFactory(
	profileEngine _interfaces.Engine,
	contractEngine _interfaces.Engine) _interfaces.EngineFactory {
	return &engineFactory{
		profileEngine:  profileEngine,
		contractEngine: contractEngine,
	}
}

func (ref *engineFactory) CreateEngine(engineName string) (_interfaces.Engine, error) {
	switch engineName {
	case values.EngineNameProfile:
		return ref.profileEngine.NewInstance(), nil
	case values.EngineNameContract:
		return ref.contractEngine.NewInstance(), nil
	default:
		return nil, errors.Errorf("Cannot choose an engine: %s", engineName)
	}
}
