package contracts

type DeleteRequest struct {
	OfferType string `uri:"offer_type" binding:"required"`
}
