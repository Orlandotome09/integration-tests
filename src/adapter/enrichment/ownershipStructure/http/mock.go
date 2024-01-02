package ownershipStructureClient

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockOwnershipStructureClient struct {
	OwnershipStructureClient
	mock.Mock
}

func (ref *MockOwnershipStructureClient) Get(id, offerType, partnerID string) (*entity.OwnershipStructure, error) {
	args := ref.Called(id, offerType, partnerID)
	return args.Get(0).(*entity.OwnershipStructure), args.Error(1)
}
