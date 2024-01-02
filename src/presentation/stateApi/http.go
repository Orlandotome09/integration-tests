package stateApi

import (
	"fmt"
	"math"
	"net/http"
	"time"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/presentation"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/stateApi/contracts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type stateApi struct {
	stateService interfaces.StateService
	cnc          interfaces.EventProcessor
}

func RegisterStateApi(ginRouterGroup *gin.RouterGroup, service interfaces.StateService, cnc interfaces.EventProcessor) {
	stateApi := &stateApi{
		stateService: service,
		cnc:          cnc,
	}
	ginRouterGroup.GET("/states/search", stateApi.SearchProfileStates)
	ginRouterGroup.GET("/states/search/:profile_id", stateApi.FindStatesByProfile)
	ginRouterGroup.GET("/state/:profile_id", stateApi.Get)
	ginRouterGroup.GET("/result/:entity_id", stateApi.GetComplianceResult)
	ginRouterGroup.POST("/check/:entity_id", stateApi.ComplianceCheck)
	ginRouterGroup.POST("/states/resync", stateApi.StatesResync)
	ginRouterGroup.POST("/states/reprocess", stateApi.StatesReprocess)
}

// SearchProfileStates godoc
// @Summary Search profile states
// @Description search profile states
// @Tags States
// @Accept json
// @Produce json
// @Param rule_name query string false "Search By rule_name"
// @Param offset query string false "Skip rows returned by setting offset"
// @Param limit query string false "Cut rows returned by setting limit"
// @Param sort_by query string false "Sort rows by sort_by"
// @Param order_by query string false "Order rows by order_by"
// @Param result_status query []string false "Filter rows by result status of the profile"
// @Param filter query string false "Filter rows by filter using profile_id, name or document_number"
// @Param partner_ids query []string false "Filter rows by one or more partner_id"
// @Param parents_ids query []string false "Filter rows by one or more parent_id"
// @Param offer_types query []string false "Filter rows by one or more offer_type"
// @Param rule_name query string false "Filter rows by one rule name"
// @Success 200 {object} contracts.SearchPaginatedResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /states/search [get]
func (ref *stateApi) SearchProfileStates(c *gin.Context) {
	var request contracts.SearchRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	searchRequest, err := request.ToDomain()
	if err != nil {
		presentation.SendBadRequestError(c, errors.WithStack(err))
		return
	}

	var profileStateList *entity.ProfileStateList
	if request.Engine == values2.EngineNameContract {
		profileStateList, err = ref.stateService.SearchContractStates(*searchRequest)
		if err != nil {
			presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
			return
		}
	} else {
		profileStateList, err = ref.stateService.SearchProfileStates(*searchRequest)
		if err != nil {
			presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
			return
		}
	}

	responses := make([]contracts.SearchResponse, 0)
	for _, profileState := range profileStateList.ProfileStates {
		response := (contracts.SearchResponse{}).FromDomain(profileState, true)
		responses = append(responses, response)
	}

	paginatedResponse := contracts.SearchPaginatedResponse{
		Profiles: responses,
		PagingResponse: presentation.PagingResponse{
			NumberOfPages: int64(math.Ceil(float64(profileStateList.Count) / float64(request.Limit))),
			TotalCount:    profileStateList.Count,
		},
	}

	c.JSON(http.StatusOK, paginatedResponse)
}

// FindStatesByProfile godoc
// @Summary Find states by profile_id
// @Description Find states by profile_id
// @Tags States
// @Accept json
// @Produce json
// @Param profile_id path string true "Profile ID"
// @Param result_status query string true "Find by result_status"
// @Param engine_name query string true "Engine name"
// @Success 200 {object} contracts.FindResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /states/search/{profile_id} [get]
func (ref *stateApi) FindStatesByProfile(c *gin.Context) {
	var request contracts.FindRequest

	if err := c.ShouldBindUri(&request); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	if err := c.ShouldBindQuery(&request); err != nil {
		presentation.SendBadRequestError(c, err)
	}

	profileID := uuid.MustParse(request.ProfileID)
	states, err := ref.stateService.FindByProfileID(profileID, request.EngineName, request.ResultStatus)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	if len(states) == 0 {
		presentation.SendNotFoundError(c, errors.Errorf("States not found, profile : %s, engine : %s, result : %s", request.ProfileID, request.EngineName, request.ResultStatus))
		return
	}

	response := &contracts.FindResponse{
		ProfileID: profileID,
	}
	for _, state := range states {
		stateResponse := contracts.StateResponse{}
		stateResponse.FromDomain(&state, false)
		response.States = append(response.States, stateResponse)
	}

	c.JSON(http.StatusOK, response)
}

// GetState godoc
// @Summary Get state from entity
// @Description Get state from by entity_id
// @Tags States
// @Accept json
// @Produce json
// @Param entity_id path string true "Entity ID"
// @Param only_pending query string false "Search only states pending"
// @Success 200 {object} contracts.StateResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /states/{entity_id} [get]
func (ref *stateApi) Get(c *gin.Context) {
	var request contracts.GetRequest

	if err := c.ShouldBindUri(&request); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	if err := c.ShouldBindQuery(&request); err != nil {
		presentation.SendBadRequestError(c, err)
	}

	entityID := uuid.MustParse(request.ProfileID)

	state, exists, err := ref.stateService.Get(entityID)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	if !exists {
		presentation.SendNotFoundError(c, errors.Errorf("Profile %s not found", entityID))
		return
	}

	response := &contracts.StateResponse{}
	response.FromDomain(state, request.OnlyPending)
	c.JSON(http.StatusOK, response)
}

// GetComplianceResult godoc
// @Summary Get compliance result
// @Description get states result for a given entity
// @Tags States
// @Accept json
// @Produce json
// @Param entity_id path string true "Entity ID"
// @Success 200 {object} contracts.ComplianceResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /result/{entity_id} [get]
func (ref *stateApi) GetComplianceResult(c *gin.Context) {
	var request contracts.GetComplianceRequest

	if err := c.ShouldBindUri(&request); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	entityID := uuid.MustParse(request.EntityID)

	state, exists, err := ref.stateService.Get(entityID)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}
	if !exists {
		presentation.SendNotFoundError(c, errors.Errorf("Entity %s not found", entityID))
		return
	}

	response := &contracts.ComplianceResponse{EntityID: entityID, EntityType: state.EngineName, Result: state.Result}
	c.JSON(http.StatusOK, response)
}

// ComplianceCheck godoc
// @Summary Post compliance check
// @Description execute compliance validation over an entity
// @Tags States
// @Accept json
// @Produce json
// @Param entity_id path string true "Entity ID"
// @Param entity_type path string true "Name of Compliance Engine "
// @Param event_type path string false "Type of new check event"
// @Success 200 {object} contracts.ComplianceResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /check/{entity_id} [post]
func (ref *stateApi) ComplianceCheck(c *gin.Context) {
	var request contracts.GetComplianceCheckRequest

	if err := c.ShouldBindUri(&request); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	entityID, err := uuid.Parse(request.EntityID)
	if err != nil {
		message := fmt.Sprintf("%s is an invalid UUID", request.EntityID)
		presentation.SendBadRequestError(c, errors.New(message))
		return
	}

	engineName, err := values2.ParseToEngineName(request.EngineName)
	if err != nil {
		message := fmt.Sprintf("%s is an invalid Engine name", request.EngineName)
		presentation.SendBadRequestError(c, errors.New(message))
		return
	}

	now := time.Now()
	event := values2.Event{
		EngineName:  engineName,
		EventType:   request.EventType,
		ParentID:    entityID,
		EntityID:    entityID,
		Date:        now,
		RequestDate: now,
	}

	state, err := ref.cnc.ExecuteForEvent(&event)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	response := contracts.ComplianceResponse{}

	response.FromDomain(*state)

	c.JSON(http.StatusOK, response)
}

// Compliance Resync godoc
// @Summary Resync compliance states
// @Description Resync compliance states
// @Tags States
// @Accept json
// @Produce json
// @Param ids path string true "Array of entities ids to resync compliance states"
// @Success 200 {object} contracts.ComplianceResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /states/resync [post]
func (ref *stateApi) StatesResync(c *gin.Context) {
	var request contracts.ResyncRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, presentation.NewErrorResponse(err.Error()))
		return
	}

	resynced, err := ref.stateService.Resync(request.Ids...)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	c.JSON(http.StatusOK, contracts.ResyncResponse{Resynced: resynced})
}

// Compliance Resync godoc
// @Summary Resync compliance states
// @Description Resync compliance states
// @Tags States
// @Accept json
// @Produce json
// @Param engine_name path string false "Name of the Engine to reprocess"
// @Param ids path string true "Array of entities ids to resync compliance states"
// @Success 200 {object} contracts.ComplianceResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /states/reprocess [post]
func (ref *stateApi) StatesReprocess(c *gin.Context) {
	var request contracts.ReprocessRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, presentation.NewErrorResponse(err.Error()))
		return
	}

	reprocessed, err := ref.stateService.Reprocess(request.EngineName, request.Ids...)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	c.JSON(http.StatusOK, contracts.ReprocessedResponse{Reprocessed: reprocessed})
}
