package watchlistTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/http/dto"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockWatchlistTranslator struct {
	WatchListTranslator
	mock.Mock
}

func (ref *MockWatchlistTranslator) ToDomain(response []dto.WatchlistResponse) (*entity2.Watchlist, error) {
	args := ref.Called(response)
	return args.Get(0).(*entity2.Watchlist), args.Error(1)
}
