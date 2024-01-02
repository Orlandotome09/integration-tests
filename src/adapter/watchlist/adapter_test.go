package watchlistAdapter

import (
	watchlistClient "bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/http/dto"
	watchlistTranslator "bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	entity "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"os"
	"reflect"
	"testing"
)

var (
	translator      *watchlistTranslator.MockWatchlistTranslator
	watchlistclient *watchlistClient.MockWatchListClient
	adapter         interfaces.WatchlistAdapter
)

func TestMain(m *testing.M) {
	translator = &watchlistTranslator.MockWatchlistTranslator{}
	watchlistclient = &watchlistClient.MockWatchListClient{}
	adapter = New(watchlistclient, translator)
	os.Exit(m.Run())
}

func TestSearch(t *testing.T) {
	documentNumber := "111"
	firstName := "Dilma"
	lastName := "Roussef"
	fullName := "Dilma Roussef"
	birthYear := 1950

	responses := []dto.WatchlistResponse{}
	sources := &entity.Watchlist{Sources: []string{"PEP"}}

	watchlistclient.On("SearchIndividual", documentNumber, firstName, lastName, fullName, "BR",
		birthYear).Return(responses, nil)
	translator.On("ToDomain", responses).Return(sources, nil)

	received, err := adapter.SearchIndividual(documentNumber, firstName, lastName, fullName, "BR", birthYear)

	if !reflect.DeepEqual(nil, err) {
		t.Errorf("\nExpected: %v \nGot: %v\n", nil, err)
	}

	if !reflect.DeepEqual(sources, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", sources, received)
	}

}
