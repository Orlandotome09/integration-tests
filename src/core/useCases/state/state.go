package state

import (
	"sync"
	"time"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type stateService struct {
	stateRepository            interfaces.StateRepository
	profileAdapter             interfaces.ProfileAdapter
	contractAdapter            interfaces.ContractAdapter
	stateEventsPublisher       interfaces.StateEventsPublisher
	complianceCommandPublisher interfaces.ComplianceCommandPublisher
	Mutex                      sync.Mutex
}

func NewStateService(repository interfaces.StateRepository,
	profileAdapter interfaces.ProfileAdapter,
	contractAdapter interfaces.ContractAdapter,
	stateEventsPublisher interfaces.StateEventsPublisher,
	complianceCommandPublisher interfaces.ComplianceCommandPublisher,
) interfaces.StateService {
	return &stateService{
		stateRepository:            repository,
		profileAdapter:             profileAdapter,
		contractAdapter:            contractAdapter,
		stateEventsPublisher:       stateEventsPublisher,
		complianceCommandPublisher: complianceCommandPublisher,
	}
}

func (ref *stateService) Get(entityID uuid.UUID) (*entity.State, bool, error) {
	return ref.stateRepository.Get(entityID)
}

func (ref *stateService) Create(entityID uuid.UUID, engineName values.EngineName) (*entity.State, error) {
	created, err := ref.stateRepository.Create(entityID, engineName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := ref.stateEventsPublisher.Send(*created, values.EventTypeStateCreated); err != nil {
		return nil, errors.WithStack(err)
	}

	return created, nil
}

func (ref *stateService) Update(state entity.State, requestDate time.Time, executionTime time.Time) error {
	saved, err := ref.stateRepository.Save(state, requestDate, executionTime)
	if err != nil {
		return errors.WithStack(err)
	}
	if saved == nil {
		logrus.Infof("[StateService] discarding state %s since there is a most recent one", state.EntityID)
		return nil
	}

	logrus.WithField("state", saved).Infof("[StateService] new state %s updated", saved.EntityID)
	if err := ref.stateEventsPublisher.Send(*saved, values.EventTypeStateChanged); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (ref *stateService) SearchProfileStates(request entity.SearchProfileStateRequest) (*entity.ProfileStateList, error) {
	states, count, err := ref.stateRepository.Search(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	profileStates := make([]entity.ProfileState, len(states))

	var wg sync.WaitGroup
	for i, state := range states {
		wg.Add(1)
		go func(idx int, stateProfile entity.State, wg *sync.WaitGroup) {
			defer wg.Done()
			profile, errGet := ref.profileAdapter.Get(stateProfile.EntityID)
			if errGet == nil {
				if profile != nil {
					profileState := entity.ProfileState{
						State:          stateProfile,
						ProfileID:      *profile.ProfileID,
						PartnerID:      profile.PartnerID,
						ParentID:       profile.ParentID,
						DocumentNumber: profile.DocumentNumber,
						Name:           profile.Name,
					}
					profileStates[idx] = profileState
				}
			}
		}(i, state, &wg)
	}
	wg.Wait()

	profileStateList := &entity.ProfileStateList{
		Count:         count,
		ProfileStates: profileStates,
	}

	return profileStateList, nil
}

func (ref *stateService) SearchContractStates(request entity.SearchProfileStateRequest) (*entity.ProfileStateList, error) {
	states, count, err := ref.stateRepository.SearchContractStates(request)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	contractStates := make([]entity.ProfileState, len(states))

	var wg sync.WaitGroup
	for i, state := range states {
		wg.Add(1)
		go func(idx int, stateContract entity.State, wg *sync.WaitGroup) {
			defer wg.Done()
			contract, _, errGet := ref.contractAdapter.Get(&stateContract.EntityID)
			if errGet == nil && contract != nil {
				profile, errGet := ref.profileAdapter.Get(*contract.ProfileID)
				if errGet == nil && profile != nil {
					contractState := entity.ProfileState{
						State:          stateContract,
						ProfileID:      *profile.ProfileID,
						PartnerID:      profile.PartnerID,
						ParentID:       profile.ParentID,
						DocumentNumber: profile.DocumentNumber,
						Name:           profile.Name,
					}
					contractStates[idx] = contractState
				}
			}

		}(i, state, &wg)
	}
	wg.Wait()

	contractStateList := &entity.ProfileStateList{
		Count:         count,
		ProfileStates: contractStates,
	}

	return contractStateList, nil
}

func (ref *stateService) FindByProfileID(profileID uuid.UUID, engine values.EngineName, result values.Result) ([]entity.State, error) {
	return ref.stateRepository.FindByProfileID(profileID, engine, result)
}

func (ref *stateService) Resync(entityIds ...string) ([]string, error) {
	if len(entityIds) == 0 {
		logrus.Infof("[StateService] No entity ids to resync")
		return nil, nil
	}

	errorGroup := new(errgroup.Group)
	errorGroup.SetLimit(20)

	var resynced []string
	for _, entityId := range entityIds {
		errorGroup.Go(func(entityId string) func() error {
			return func() error {
				id, err := uuid.Parse(entityId)
				if err != nil {
					logrus.Errorf("[StateService] Invalid entity id %s from repository: %+v", entityId, errors.WithStack(err))
					return values.NewErrorValidation(err.Error())
				}

				state, exists, err := ref.stateRepository.Get(id)
				if err != nil {
					logrus.Errorf("[StateService] Error getting state of entity id %s from repository: %+v", entityId, errors.WithStack(err))
					return err
				}
				if !exists {
					logrus.Errorf("[StateService] Error getting entity id %s not found", entityId)
					return nil
				}

				if err := ref.stateEventsPublisher.Send(*state, values.EventTypeStateResync); err != nil {
					logrus.Errorf("[StateService] Failed to send event state entity id %s : %+v", entityId, errors.WithStack(err))
					return err
				}

				ref.Mutex.Lock()
				defer ref.Mutex.Unlock()
				resynced = append(resynced, entityId)

				return nil
			}
		}(entityId))
	}

	if err := errorGroup.Wait(); err != nil {
		return nil, err
	}

	return resynced, nil
}

func (ref *stateService) Reprocess(engineName string, entityIds ...string) ([]string, error) {
	if len(entityIds) == 0 {
		logrus.Infof("[StateService] No entity ids to reprocess")
		return nil, nil
	}

	if engineName == "" {
		engineName = values.EngineNameProfile
	}

	errorGroup := new(errgroup.Group)
	errorGroup.SetLimit(20)

	var reprocessed []string
	for _, entityId := range entityIds {
		errorGroup.Go(func(entityId string) func() error {
			return func() error {
				id, err := uuid.Parse(entityId)
				if err != nil {
					logrus.Errorf("[StateService] Invalid entity id %s : %+v", entityId, errors.WithStack(err))
					return values.NewErrorValidation(err.Error())
				}

				if _, err := ref.complianceCommandPublisher.SendCommand(id, &id, engineName); err != nil {
					logrus.Errorf("[StateService] Failed to send command with entity id %s : %+v", entityId, errors.WithStack(err))
					return err
				}

				ref.Mutex.Lock()
				defer ref.Mutex.Unlock()
				reprocessed = append(reprocessed, entityId)

				return nil
			}
		}(entityId))
	}

	if err := errorGroup.Wait(); err != nil {
		return nil, err
	}

	return reprocessed, nil
}
