package watchlistClient

import (
	dto2 "bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/http/dto"
	"github.com/stretchr/testify/mock"
)

type MockWatchListClient struct {
	WatchListClient
	mock.Mock
}

func (ref *MockWatchListClient) SearchIndividual(documentNumber, firstName, lastName, fullName, countryCode string,
	birthYear int) ([]dto2.WatchlistResponse, error) {
	args := ref.Called(documentNumber, firstName, lastName, fullName, countryCode, birthYear)
	return args.Get(0).([]dto2.WatchlistResponse), args.Error(1)
}

func (ref *MockWatchListClient) SearchCompany(legalName, countryCode string) ([]dto2.WatchlistResponse, error) {
	args := ref.Called(legalName, countryCode)
	return args.Get(0).([]dto2.WatchlistResponse), args.Error(1)
}
