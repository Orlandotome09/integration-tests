package doaApi

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/presentation"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/doaApi/contracts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

type doaCallback struct {
	doaService _interfaces.DOAService
}

func RegisterDOACallback(ginRouterGroup *gin.RouterGroup, service _interfaces.DOAService) {
	doaCallback := &doaCallback{
		service,
	}
	ginRouterGroup.POST("/doa/callback", doaCallback.ReceiveCallback)
	ginRouterGroup.GET("/doa/result/:entity_id", doaCallback.GetLastResult)
	ginRouterGroup.POST("/doa/result", doaCallback.Create)
}

/*
// ReceiveCallback doa godoc
// @Summary Receive DOA callback
// @Description Receive DOA callback
// @Tags DOA
// @Accept json
// @Produce json
// @Param callback body contracts.DOACallback true "Body"
// @Success 200 {object} entity.DOAResult
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /doa/callback [post]
*/
func (ref *doaCallback) ReceiveCallback(c *gin.Context) {
	request := &contracts.DOACallback{}
	if err := c.ShouldBindJSON(request); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	id := request.RequestID
	scores := request.Scores.ToDomain()

	enriched, err := ref.doaService.EnrichWithScores(&id, scores)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).
			Handle(c)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if enriched == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot find result to enrich"})
		return
	}

	c.JSON(http.StatusOK, enriched)
}

/*
// GetLastResult Get doa godoc
// @Summary Get Last DOA result by entity_id
// @Description Get Last DOA result by entity_id
// @Tags DOA
// @Accept json
// @Produce json
// @Param entity_id path string true "Entity ID"
// @Success 200 {object} contracts.DOAResult
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /doa/result/{entity_id} [get]
*/
func (ref *doaCallback) GetLastResult(c *gin.Context) {
	entityID := uuid.MustParse(c.Param("entity_id"))
	documentID := uuid.MustParse(c.Query("document_id"))

	result, err := ref.doaService.FindLastResult(&entityID, &documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if result == nil {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	doaResult := &contracts.DOAResult{}
	doaResult.FromDomain(*result)

	c.JSON(http.StatusOK, doaResult)
}

/*
// Create doa godoc
// @Summary Create a DOA result
// @Description Create a DOA result
// @Tags DOA
// @Accept json
// @Produce json
// @Param doa body contracts.DOAResult true "Body"
// @Success 200 {object} contracts.DOAResult
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /doa/result [post]
*/
func (ref *doaCallback) Create(c *gin.Context) {
	body := &contracts.DOAResult{}
	if err := c.ShouldBindJSON(body); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	doaResult := body.ToDomain()
	saved, err := ref.doaService.Save(&doaResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := &contracts.DOAResult{}
	result.FromDomain(*saved)

	c.JSON(http.StatusOK, result)

}
