package overrideApi

import (
	"net/http"
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/presentation"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/overrideApi/contracts"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type overrideApi struct {
	ginRouterGroup  *gin.RouterGroup
	overrideService _interfaces.OverrideService
	cnc             _interfaces.EventProcessor
}

func RegisterOverrideApi(ginRouterGroup *gin.RouterGroup, overrideService _interfaces.OverrideService, cnc _interfaces.EventProcessor) {
	override := &overrideApi{
		ginRouterGroup:  ginRouterGroup,
		overrideService: overrideService,
		cnc:             cnc,
	}
	ginRouterGroup.POST("/override", override.Save)
	ginRouterGroup.DELETE("/override", override.Delete)
}

/*
// Save Override godoc
// @Summary Save an override: Create an override; if the override already exists, update it instead
// @Description Save an override
// @Tags Overrides
// @Accept json
// @Produce json
// @Param override body contracts.SaveRequest true "Save Override"
// @Success 200 {object} string
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 412 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /override [post]
*/
func (ref *overrideApi) Save(c *gin.Context) {
	var request contracts.SaveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	override := request.ToDomain()

	if err := ref.overrideService.Save(override); err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	_, err := ref.cnc.ExecuteForEvent(translateOverrideToEvent(override))
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	c.JSON(http.StatusOK, "ok")
}

/*
// Delete Override godoc
// @Summary Delete an override
// @Description Delete an override
// @Tags Overrides
// @Accept json
// @Produce json
// @Param override body contracts.DeleteRequest true "delete Override"
// @Success 200 {object} string
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /override [delete]
*/
func (ref *overrideApi) Delete(c *gin.Context) {
	var request contracts.DeleteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	override := request.ToDomain()

	err := ref.overrideService.Delete(override)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	_, err = ref.cnc.ExecuteForEvent(translateOverrideToEvent(override))
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	c.JSON(http.StatusOK, "deleted")
}

func translateOverrideToEvent(override entity.Override) *values2.Event {

	engineName := ""

	if override.EntityType.ToString() == values2.EngineNameProfile ||
		override.EntityType.ToString() == values2.EngineNameContract {
		engineName = override.EntityType.ToString()
	} else {
		engineName = values2.EngineNameProfile
	}

	now := time.Now()
	event := values2.Event{
		EngineName:  engineName,
		EventType:   values2.EventTypeOverrideDeleted,
		ParentID:    *override.ParentID,
		EntityID:    override.EntityID,
		Date:        now,
		RequestDate: now,
	}

	return &event
}
