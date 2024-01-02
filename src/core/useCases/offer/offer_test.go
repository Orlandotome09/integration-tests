package offer

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	validate        *validator.Validate
	offerRepository *mocks.OfferRepository
	service         interfaces.OfferService
)

func TestMain(m *testing.M) {
	validate = validator.New()
	offerRepository = &mocks.OfferRepository{}
	service = NewOfferService(validate, offerRepository)
	os.Exit(m.Run())
}

func TestCreate(t *testing.T) {
	offer := values.Offer{Type: "testOffer", Product: "Webpayments"}
	saved := &values.Offer{}

	offerRepository.On("Create", offer).Return(saved, nil)

	expected := saved
	received, err := service.Create(offer)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestGet(t *testing.T) {
	offerType := "type"
	offer := &values.Offer{}

	offerRepository.On("Get", offerType).Return(offer, nil)

	expected := offer
	received, err := service.Get(offerType)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestUpdate(t *testing.T) {
	offer := values.Offer{Type: "testOffer", Product: "Webpayments"}
	saved := &values.Offer{}

	offerRepository.On("Update", offer).Return(saved, nil)

	expected := saved
	received, err := service.Update(offer)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestDelete(t *testing.T) {
	offerType := "test"

	offerRepository.On("Delete", offerType).Return(nil)

	err := service.Delete(offerType)

	assert.Nil(t, err)
}

func TestList(t *testing.T) {
	offers := []values.Offer{}

	offerRepository.On("List").Return(offers, nil)

	expected := offers
	received, err := service.List()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}
