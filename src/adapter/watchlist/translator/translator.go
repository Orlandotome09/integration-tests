package watchlistTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/watchlist/http/dto"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"encoding/json"
	"strings"
)

type WatchListTranslator interface {
	ToDomain(responses []dto.WatchlistResponse) (*entity2.Watchlist, error)
}

type watchListTranslator struct{}

func New() WatchListTranslator {
	return &watchListTranslator{}
}

func (ref *watchListTranslator) ToDomain(responses []dto.WatchlistResponse) (*entity2.Watchlist, error) {

	sources := []string{}
	for _, response := range responses {
		for _, source := range response.Sources {
			sources = append(sources, translateSource(source))
		}
	}

	metadata, err := json.Marshal(responses)
	if err != nil {
		return nil, err
	}

	watchlist := &entity2.Watchlist{Sources: sources, Metadata: metadata}

	return watchlist, nil
}

func translateSource(inputSource string) string {
	source := strings.ToUpper(inputSource)
	source = strings.Replace(source, " ", "_", -1)
	return source
}
