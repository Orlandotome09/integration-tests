package restrictivelists

import (
	restrictiveListsHttpClient "bitbucket.org/bexstech/temis-compliance/src/adapter/restrictiveLists/http"
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
	"time"
)

type restrictiveListsAdapter struct {
	httpClient restrictiveListsHttpClient.RestrictiveListsHttpClient
}

func NewRestrictiveListsAdapter(httpClient restrictiveListsHttpClient.RestrictiveListsHttpClient) _interfaces.RestrictiveListsAdapter {
	return &restrictiveListsAdapter{
		httpClient: httpClient,
	}
}

func (ref *restrictiveListsAdapter) OccurrenceInBlackList(document string, name string) (*entity.BlacklistStatus, error) {
	internalList, err := ref.httpClient.SearchInternalList(document, name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if len(internalList) == 0 {
		return nil, nil
	}

	return &entity.BlacklistStatus{
		Justification: entity.Justification{
			AddedAt:  internalList[0].CreatedAt,
			Author:   internalList[0].Author,
			Comments: []string{internalList[0].Justification},
		},
	}, nil
}

func (ref *restrictiveListsAdapter) OccurrenceInPepList(documentNumber string) (*entity.PepInformation, error) {
	pepInformation, err := ref.httpClient.SearchPepList(documentNumber)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if pepInformation == nil {
		return nil, nil
	}

	startDate := time.Time{}
	endDate := time.Time{}

	if pepInformation.StartDate != "" {
		result, err := time.Parse("2006-01-02", pepInformation.StartDate)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		startDate = result
	}

	if pepInformation.EndDate != "" {
		result, err := time.Parse("2006-01-02", pepInformation.EndDate)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		endDate = result
	}

	return &entity.PepInformation{
		DocumentNumber: pepInformation.DocumentNumber,
		Name:           pepInformation.Name,
		Role:           pepInformation.Role,
		Institution:    pepInformation.Institution,
		StartDate:      startDate,
		EndDate:        endDate,
		Source:         pepInformation.Source,
		CreatedAt:      pepInformation.CreatedAt,
	}, nil
}
