package offerApi

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/presentation"
	"bitbucket.org/bexstech/temis-compliance/src/presentation/offerApi/contracts"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type offerApi struct {
	offerService interfaces.OfferService
}

var parser = func(original interface{}) interface{} {
	response := contracts.OfferResponse{}
	offer := original.(*values.Offer)
	response = response.FromDomain(*offer)
	return response
}

func RegisterOfferApi(ginRouterGroup *gin.RouterGroup, offerService interfaces.OfferService) {
	offerApi := &offerApi{
		offerService: offerService,
	}
	ginRouterGroup.POST("/offers", offerApi.Create)
	ginRouterGroup.GET("/offers/:offer_type", offerApi.Get)
	ginRouterGroup.PATCH("/offers/:offer_type", offerApi.Update)
	ginRouterGroup.DELETE("/offers/:offer_type", offerApi.Delete)
	ginRouterGroup.GET("/offers", offerApi.List)
}

/*
// Create offer godoc
// @Summary Create a Offer
// @Description Create a Offer
// @Tags Offer
// @Accept json
// @Produce json
// @Param offer body contracts.CreateOfferRequest true "Body"
// @Success 200 {object} contracts.OfferResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /offers [post]
*/
func (ref *offerApi) Create(c *gin.Context) {
	var request contracts.CreateOfferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offer := request.ToDomain()

	created, err := ref.offerService.Create(*offer)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).WithParser(parser).Handle(c)
		return
	}

	response := contracts.OfferResponse{}.FromDomain(*created)
	c.JSON(http.StatusCreated, response)
}

/*
// Get Offer godoc
// @Summary Get Offer by id
// @Description Get Offer by id
// @Tags Offer
// @Accept json
// @Produce json
// @Param offer_type path string true "Offer ID"
// @Success 200 {object} contracts.OfferResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /offers/{offer_type} [get]
*/
func (ref *offerApi) Get(c *gin.Context) {
	var request contracts.GetOfferRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offer, err := ref.offerService.Get(request.OfferType)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}
	if offer == nil {
		presentation.SendNotFoundError(c, errors.New("Offer not found"))
		return
	}

	response := contracts.OfferResponse{}.FromDomain(*offer)

	c.JSON(http.StatusOK, response)
}

/*
// Update offer godoc
// @Summary Update a Offer
// @Description Update a Offer
// @Tags Offer
// @Accept json
// @Produce json
// @Param offer body contracts.UpdateOfferRequest true "Body"
// @Param offer_type path string true "Offer ID"
// @Success 200 {object} contracts.OfferResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /offers/{offer_type} [patch]
*/
func (ref *offerApi) Update(c *gin.Context) {
	var request contracts.UpdateOfferRequest
	if err := c.ShouldBindUri(&request.UpdateOfferRequestURI); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}
	if err := c.ShouldBindJSON(&request.UpdateOfferRequestJSON); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	offer := request.ToDomain()

	updated, err := ref.offerService.Update(*offer)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).WithParser(parser).Handle(c)
		return
	}

	response := contracts.OfferResponse{}.FromDomain(*updated)
	c.JSON(http.StatusOK, response)
}

/*
// Delete offer godoc
// @Summary Delete a Offer
// @Description Delete a Offer
// @Tags Offer
// @Accept json
// @Produce json
// @Param offer_type path string true "Offer ID"
// @Success 200 {object} string
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /offers/{offer_type} [patch]
*/
func (ref *offerApi) Delete(c *gin.Context) {
	var request contracts.DeleteRequest
	if err := c.ShouldBindUri(&request); err != nil {
		presentation.SendBadRequestError(c, err)
		return
	}

	err := ref.offerService.Delete(request.OfferType)
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).WithParser(parser).Handle(c)
		return
	}

	c.JSON(http.StatusOK, "ok")
}

/*
// List offer godoc
// @Summary List a Offer
// @Description List a Offer
// @Tags Offer
// @Accept json
// @Produce json
// @Success 200 {object} []contracts.OfferResponse
// @Failure 400 {object} presentation.ErrorResponse
// @Failure 404 {object} presentation.ErrorResponse
// @Failure 500 {object} presentation.ErrorResponse
// @Router /offers [get]
*/
func (ref *offerApi) List(c *gin.Context) {

	offers, err := ref.offerService.List()
	if err != nil {
		presentation.NewErrorHandler(errors.WithStack(err)).Handle(c)
		return
	}

	responses := make([]contracts.OfferResponse, len(offers))
	for i, offer := range offers {
		response := contracts.OfferResponse{}.FromDomain(offer)
		responses[i] = response
	}

	c.JSON(http.StatusOK, responses)
}
