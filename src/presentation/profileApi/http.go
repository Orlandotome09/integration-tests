package profileApi

import (
	"net/http"
	"time"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/presentation"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/profileApi/contracts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type profileApi struct {
	complianceProfileService interfaces.ComplianceProfileService
	eventProcessor           interfaces.EventProcessor
}

func RegisterProfileApi(ginRouterGroup *gin.RouterGroup, complianceProfileService interfaces.ComplianceProfileService, eventProcessor interfaces.EventProcessor) {
	var profileApi = &profileApi{
		complianceProfileService: complianceProfileService,
		eventProcessor:           eventProcessor,
	}
	ginRouterGroup.GET("/profile/:profile_id", profileApi.Get)
}

// Get Profile godoc
// @Summary Get Profile by id
// @Description Get Profile by id
// @Tags Profile
// @Accept json
// @Produce json
// @Param profile_id path string true "Profile ID"
// @Success 200 {object} contracts.ProfileResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /profile/{profile_id} [get]
func (ref *profileApi) Get(c *gin.Context) {
	var request contracts.GetProfileRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profileID := uuid.MustParse(request.ProfileID)

	profile, err := ref.complianceProfileService.Get(profileID)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	if profile != nil {
		response := contracts.ProfileResponse{}.FromDomain(*profile)
		c.JSON(http.StatusOK, response)
		return
	}

	//If profile snapshot does not exist, run compliance engine to create it
	err = ref.executeCompliance(profileID)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	profile, err = ref.complianceProfileService.Get(profileID)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	if profile == nil {
		presentation.SendNotFoundError(c, errors.New("Profile not found"))
		return
	}

	response := contracts.ProfileResponse{}.FromDomain(*profile)
	c.JSON(http.StatusOK, response)

}

func (ref *profileApi) executeCompliance(profileID uuid.UUID) error {

	now := time.Now()
	event := values2.Event{
		EngineName:  values2.EngineNameProfile,
		EventType:   values2.EventTypeProfileResync,
		ParentID:    profileID,
		EntityID:    profileID,
		Date:        now,
		RequestDate: now,
	}

	_, err := ref.eventProcessor.ExecuteForEvent(&event)
	return err
}
