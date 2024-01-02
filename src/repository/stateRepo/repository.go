package stateRepo

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/repository"
	"bitbucket.org/bexstech/temis-compliance/src/repository/errorHandler"
	"bitbucket.org/bexstech/temis-compliance/src/repository/model"
)

type stateRepository struct {
	db           *gorm.DB
	errorHandler errorHandler.ErrorHandler
}

func NewStateSqlRepository(db *gorm.DB, errorHandler errorHandler.ErrorHandler) _interfaces.StateRepository {
	return &stateRepository{
		db:           db,
		errorHandler: errorHandler,
	}
}

func (ref *stateRepository) Get(entityID uuid.UUID) (*entity.State, bool, error) {
	record, exists, err := ref.get(entityID)
	if err != nil {
		return nil, false, errors.WithStack(err)
	}
	if !exists {
		return nil, false, nil
	}

	return record.ToDomain(), true, nil
}

func (ref *stateRepository) Create(entityID uuid.UUID, engineName values.EngineName) (*entity.State, error) {
	now := time.Now()
	record := &model.State{
		EntityID:   entityID,
		EngineName: engineName,
		Result:     string(values.ResultStatusCreated),
		Version:    1,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	result := ref.db.Create(record)
	if ref.errorHandler.IsRecordDuplicated(result) {
		logrus.Infof("[StateRepository] Create State duplicated. Using current")
		current, _, err := ref.get(entityID)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return current.ToDomain(), nil
	}

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return record.ToDomain(), errors.WithStack(result.Error)
}

func (ref *stateRepository) Save(state entity.State, requestDate time.Time, executionTime time.Time) (*entity.State, error) {

	record := model.StateFromDomain(&state)
	record.Version++
	record.RequestDate = requestDate
	record.ExecutionTime = executionTime
	record.UpdatedAt = time.Now()

	result := ref.db.Save(&record).Where("version = ?", record.Version-1)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	if result.RowsAffected == 0 {
		logrus.WithField("entityId", state.EntityID).Error("[StateRepository] returning message to queue since state was already modified by another routine")
		return nil, fmt.Errorf("entity %s has a more recente validation", record.EntityID.String())
	}

	saved, _, err := ref.Get(record.EntityID)
	if err != nil {
		return nil, nil
	}

	return saved, nil
}

func (ref *stateRepository) Search(request entity.SearchProfileStateRequest) (states []entity.State, count int64, err error) {
	var records []model.State

	queryBuilder := &repository.QueryBuilder{}
	queryBuilder = queryBuilder.WithAnd("not jsonb_typeof(profile_states.validation_steps_results) = 'null'")
	queryBuilder = queryBuilder.WithAnd("pending")
	queryBuilder = queryBuilder.WithAnd("engine_name = 'PROFILE'")

	if request.Name != "" {
		queryBuilder = queryBuilder.WithAnd(fmt.Sprintf("profiles.name ILIKE '%%%s%%'", request.Name))
	}

	if request.DocumentNumber != "" {
		queryBuilder = queryBuilder.WithAnd(fmt.Sprintf("profiles.document_number = '%s'", request.DocumentNumber))
	}

	if request.RuleName != "" {
		queryBuilder = queryBuilder.WithAnd(fmt.Sprintf("'%s' = any(rule_names)", request.RuleName))
	}

	if request.ProfileID != nil {
		queryBuilder = queryBuilder.WithAnd(fmt.Sprintf("profile_states.entity_id = '%s'", request.ProfileID.String()))
	}

	subQuery1 := &repository.QueryBuilder{}

	if len(request.PartnerIDs) > 0 {
		partnerIds := "'" + strings.Join(request.PartnerIDs, "','") + "'"
		subQuery1 = subQuery1.WithOr(fmt.Sprintf("profiles.partner_id in (%s)", partnerIds))
	}

	if len(request.ParentIDs) > 0 {
		var uuids []string
		for _, parentID := range request.ParentIDs {
			uuids = append(uuids, parentID.String())
		}
		parentIds := "'" + strings.Join(uuids, "','") + "'"
		subQuery1 = subQuery1.WithOr(fmt.Sprintf("profiles.parent_id in (%s)", parentIds))
	}

	queryBuilder = queryBuilder.WithAndSubQuery(subQuery1)

	if len(request.OfferTypes) > 0 {
		offers := "'" + strings.Join(request.OfferTypes, "','") + "'"
		queryBuilder = queryBuilder.WithAnd(fmt.Sprintf("profiles.offer_type in (%s)", offers))
	}

	if len(request.ResultStatus) > 0 {
		var statuses []string
		for _, resultStatus := range request.ResultStatus {
			statuses = append(statuses, resultStatus.ToString())
		}
		resultStatus := "'" + strings.Join(statuses, "','") + "'"
		queryBuilder = queryBuilder.WithAnd(fmt.Sprintf("profile_states.result in (%s)", resultStatus))
	}

	where := queryBuilder.Build()

	if request.DocumentNumber != "" || request.Name != "" || len(request.PartnerIDs) > 0 || len(request.ParentIDs) > 0 || len(request.OfferTypes) > 0 || request.SortBy == "profiles.name" {

		errCount := ref.db.Table("profile_states").Joins("JOIN profiles ON profiles.profile_id = profile_states.entity_id").Where(where).
			Count(&count).Error
		if errCount != nil {
			return nil, 0, errCount
		}

		if count < request.OffSet {
			request.OffSet = 0
		}

		errFind := ref.db.Joins("JOIN profiles ON profiles.profile_id = profile_states.entity_id").Where(where).
			Limit(int(request.Limit)).
			Offset(int(request.OffSet)).
			Order(request.SortBy + " " + request.OrderBy).
			Find(&records).Error
		if errFind != nil {
			return nil, 0, errCount
		}
	} else {

		errCount := ref.db.Table("profile_states").Where(where).
			Count(&count).Error
		if errCount != nil {
			return nil, 0, errCount
		}

		if count < request.OffSet {
			request.OffSet = 0
		}

		errFind := ref.db.Where(where).
			Limit(int(request.Limit)).
			Offset(int(request.OffSet)).
			Order(request.SortBy + " " + request.OrderBy).
			Find(&records).Error
		if errFind != nil {
			return nil, 0, errCount
		}
	}

	if count == 0 {
		return nil, 0, nil
	}

	for _, record := range records {
		state := *record.ToDomain()
		states = append(states, state)
	}

	return states, count, nil
}

func (ref *stateRepository) SearchContractStates(request entity.SearchProfileStateRequest) (states []entity.State, count int64, err error) {
	var records []model.State

	queryBuilder := &repository.QueryBuilder{}
	queryBuilder = queryBuilder.WithAnd("not jsonb_typeof(profile_states.validation_steps_results) = 'null'")
	queryBuilder = queryBuilder.WithAnd("profile_states.pending")
	queryBuilder = queryBuilder.WithAnd("profile_states.engine_name = 'CONTRACT'")
	queryBuilder = queryBuilder.WithAnd("profile_states.result = 'APPROVED'")

	if request.Name != "" {
		queryBuilder = queryBuilder.WithAnd(fmt.Sprintf("profiles.name ILIKE '%%%s%%'", request.Name))
	}

	if request.DocumentNumber != "" {
		queryBuilder = queryBuilder.WithAnd(fmt.Sprintf("profiles.document_number = '%s'", request.DocumentNumber))
	}

	if request.RuleName != "" {
		queryBuilder = queryBuilder.WithAnd(fmt.Sprintf("'%s' = any(rule_names)", request.RuleName))
	}

	if request.ProfileID != nil {
		queryBuilder = queryBuilder.WithAnd(fmt.Sprintf("ps2.entity_id = '%s'", request.ProfileID.String()))
	}

	subQuery1 := &repository.QueryBuilder{}

	if len(request.PartnerIDs) > 0 {
		partnerIds := "'" + strings.Join(request.PartnerIDs, "','") + "'"
		subQuery1 = subQuery1.WithOr(fmt.Sprintf("profiles.partner_id in (%s)", partnerIds))
	}

	if len(request.ParentIDs) > 0 {
		var uuids []string
		for _, parentID := range request.ParentIDs {
			uuids = append(uuids, parentID.String())
		}
		parentIds := "'" + strings.Join(uuids, "','") + "'"
		subQuery1 = subQuery1.WithOr(fmt.Sprintf("profiles.parent_id in (%s)", parentIds))
	}

	queryBuilder = queryBuilder.WithAndSubQuery(subQuery1)

	if len(request.OfferTypes) > 0 {
		offers := "'" + strings.Join(request.OfferTypes, "','") + "'"
		queryBuilder = queryBuilder.WithAnd(fmt.Sprintf("profiles.offer_type in (%s)", offers))
	}

	where := queryBuilder.Build()

	errCount := ref.db.Table("profile_states").
		Joins("JOIN contracts ON contracts.contract_id = profile_states.entity_id").
		Joins("JOIN profiles ON profiles.profile_id = contracts.profile_id").
		Joins("JOIN profile_states ps2 ON ps2.entity_id = profiles.profile_id").
		Where(where).
		Where("ps2.result != 'REJECTED'").
		Count(&count).Error
	if errCount != nil {
		return nil, 0, errCount
	}

	if count < request.OffSet {
		request.OffSet = 0
	}

	errFind := ref.db.Joins("JOIN contracts ON contracts.contract_id = profile_states.entity_id").
		Joins("JOIN profiles ON profiles.profile_id = contracts.profile_id").
		Joins("JOIN profile_states ps2 ON ps2.entity_id = profiles.profile_id").
		Where(where).
		Where("ps2.result != 'REJECTED'").
		Limit(int(request.Limit)).
		Offset(int(request.OffSet)).
		Order(request.SortBy + " " + request.OrderBy).
		Find(&records).Error
	if errFind != nil {
		return nil, 0, errCount
	}

	if count == 0 {
		return nil, 0, nil
	}

	for _, record := range records {
		state := *record.ToDomain()
		states = append(states, state)
	}

	return states, count, nil
}

func (ref *stateRepository) FindByProfileID(profileID uuid.UUID,
	engine values.EngineName, result values.Result) ([]entity.State, error) {
	var records []model.State

	switch engine {
	case values.EngineNameProfile:
		if db := ref.db.Find(&records, "entity_id = ?", profileID); db.Error != nil {
			return nil, db.Error
		}
	default:
		return nil, fmt.Errorf("invalid engine name: %s", engine)
	}

	states := make([]entity.State, len(records))
	for i, record := range records {
		states[i] = *record.ToDomain()
	}

	return states, nil
}

func (ref *stateRepository) get(entityID uuid.UUID) (profileState *model.State, exists bool, err error) {
	record := &model.State{}

	result := ref.db.Find(record, "entity_id = ?", entityID)
	if result.Error != nil {
		return nil, false, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, false, nil
	}

	return record, true, nil
}
