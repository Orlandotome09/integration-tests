package offer

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type offerService struct {
	validate        *validator.Validate
	offerRepository interfaces.OfferRepository
}

func NewOfferService(validate *validator.Validate,
	offerRepository interfaces.OfferRepository) interfaces.OfferService {
	return &offerService{
		validate:        validate,
		offerRepository: offerRepository,
	}
}

func (ref *offerService) Create(offer values.Offer) (*values.Offer, error) {
	err := ref.validate.Struct(offer)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	created, err := ref.offerRepository.Create(offer)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return created, nil
}
func (ref *offerService) Get(offerType string) (*values.Offer, error) {
	offer, err := ref.offerRepository.Get(offerType)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return offer, err
}

func (ref *offerService) Update(offer values.Offer) (*values.Offer, error) {
	err := ref.validate.Struct(offer)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	updated, err := ref.offerRepository.Update(offer)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return updated, nil
}

func (ref *offerService) Delete(offerType string) error {
	return ref.offerRepository.Delete(offerType)
}

func (ref *offerService) List() ([]values.Offer, error) {
	return ref.offerRepository.List()
}
