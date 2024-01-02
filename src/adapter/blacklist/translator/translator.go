package blacklistTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/blacklist/http/dto"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type BlacklistTranslator interface {
	ToDomain(response *dto.BlacklistResponse) *entity.BlacklistStatus
}

type blacklistTranslator struct{}

func New() BlacklistTranslator {
	return &blacklistTranslator{}
}

func (ref *blacklistTranslator) ToDomain(response *dto.BlacklistResponse) *entity.BlacklistStatus {

	comments := make([]string, 0)

	for _, justification := range response.Justifications {
		comments = append(comments, justification.Justification)
	}

	justification := entity.Justification{
		AddedAt: response.CreatedAt,
		// TODO we don't have author
		Author:   "",
		Comments: comments,
	}

	status := &entity.BlacklistStatus{
		Status:        response.Status,
		Justification: justification,
	}

	return status
}
