package offerRepo

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/repository/errorHandler"
	"bitbucket.org/bexstech/temis-compliance/src/repository/model"
)

type offerRepository struct {
	db           *gorm.DB
	errorHandler errorHandler.ErrorHandler
}

func NewOfferRepository(db *gorm.DB, errorHandler errorHandler.ErrorHandler) interfaces.OfferRepository {
	return &offerRepository{
		db:           db,
		errorHandler: errorHandler,
	}
}

func (ref *offerRepository) Create(offer values.Offer) (*values.Offer, error) {
	record := model.Offer{}.FromDomain(offer)

	result := ref.db.Create(record)
	if result.Error != nil && !ref.errorHandler.IsRecordDuplicated(result) {
		return nil, errors.WithStack(result.Error)
	}

	current, err := ref.get(record.OfferType)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if ref.errorHandler.IsRecordDuplicated(result) {
		return nil, values.NewErrorDuplicated(current.ToDomain())
	}

	return current.ToDomain(), nil
}

func (ref *offerRepository) Get(offerType string) (*values.Offer, error) {
	record, err := ref.get(offerType)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if record == nil {
		return nil, nil
	}

	return record.ToDomain(), nil
}

func (ref *offerRepository) Update(offer values.Offer) (*values.Offer, error) {
	record := model.Offer{}.FromDomain(offer)

	result := ref.db.Model(record).Updates(record)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, values.NewErrorPreconditionNotFound(
			fmt.Sprintf("Offer (type): %v", record.OfferType))
	}

	updated, err := ref.get(record.OfferType)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return updated.ToDomain(), nil
}

func (ref *offerRepository) Delete(offerType string) error {

	result := ref.db.Delete(&model.Offer{OfferType: offerType})
	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	if result.RowsAffected == 0 {
		return values.NewErrorPreconditionNotFound(
			fmt.Sprintf("Offer (type): %v", offerType))
	}

	return nil
}

func (ref *offerRepository) List() ([]values.Offer, error) {
	records := []model.Offer{}

	if result := ref.db.Find(&records); result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	offers := make([]values.Offer, len(records))
	for i, record := range records {
		offer := record.ToDomain()
		offers[i] = *offer
	}

	return offers, nil
}

func (ref *offerRepository) get(offerType string) (*model.Offer, error) {
	record := &model.Offer{}

	result := ref.db.Find(record, "offer_type = ?", offerType)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return record, nil
}
