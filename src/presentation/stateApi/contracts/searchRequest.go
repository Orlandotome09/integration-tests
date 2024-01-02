package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/presentation"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

const (
	minLimitAllowed = 0
	maxLimitAllowed = 100
)

type SearchRequest struct {
	presentation.Paging
	ResultStatus []string        `form:"result_status[]"`
	Filter       string          `form:"filter"`
	PartnerIDs   []string        `form:"partner_ids[]"`
	ParentIDs    []string        `form:"parent_ids[]"`
	OfferTypes []string         `form:"offer_types[]"`
	RuleName   values2.RuleName `form:"rule_name"`
	Engine     string           `form:"engine"`
}

func (ref *SearchRequest) ToDomain() (*entity.SearchProfileStateRequest, error) {
	if ref.Limit <= 0 {
		message := fmt.Sprintf("Limit can't be less or equal %d", minLimitAllowed)
		return nil, errors.New(message)
	}

	if ref.Limit > maxLimitAllowed {
		message := fmt.Sprintf("Limit can't be more than %d", maxLimitAllowed)
		return nil, errors.New(message)
	}

	results := []values2.Result{}
	for _, status := range ref.ResultStatus {
		result := values2.Result(status)
		if result != "" {
			err := result.Validate()
			if err != nil {
				return nil, err
			}
		}
		results = append(results, result)
	}

	if ref.Paging.Limit == 0 {
		ref.Paging.Limit = -1
	}

	if ref.Paging.OffSet == 0 {
		ref.Paging.OffSet = -1
	}

	ref.SortBy = ref.normalizeSortBy(ref.Paging.SortBy)
	ref.OrderBy = ref.normalizeOrderBy(ref.Paging.OrderBy)
	searchProfileRequest := &entity.SearchProfileStateRequest{
		PartnerIDs:   ref.PartnerIDs,
		OfferTypes:   ref.OfferTypes,
		RuleName:     ref.RuleName,
		ResultStatus: results,
		Paging:       ref.Paging.ToDomain(),
	}

	if ref.ParentIDs != nil {
		for _, parentID := range ref.ParentIDs {
			parsedParentID, err := uuid.Parse(parentID)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			searchProfileRequest.ParentIDs = append(searchProfileRequest.ParentIDs, parsedParentID)
		}
	}

	if ref.Filter == "" {
		return searchProfileRequest, nil
	}

	normalizedDocument := ref.normalizeDocument(ref.Filter)

	if profileId, err := uuid.Parse(ref.Filter); err == nil {
		searchProfileRequest.ProfileID = &profileId
	} else if isMatch, _ := regexp.MatchString("([0-9]{14})|([0-9]{11})", normalizedDocument); isMatch {
		searchProfileRequest.DocumentNumber = normalizedDocument
	} else {
		searchProfileRequest.Name = ref.Filter
	}

	return searchProfileRequest, nil
}

func (ref *SearchRequest) normalizeOrderBy(orderBy string) string {
	if orderBy != "" && strings.ToLower(orderBy) != "asc" && strings.ToLower(orderBy) != "desc" {
		return "asc"
	}

	return ref.Paging.OrderBy
}

func (ref *SearchRequest) normalizeDocument(document string) string {
	document = strings.ReplaceAll(document, ".", "")
	document = strings.ReplaceAll(document, "-", "")
	document = strings.ReplaceAll(document, "/", "")
	return document
}

func (ref *SearchRequest) normalizeSortBy(sortBy string) string {
	sortByMap := make(map[string]string)
	sortByMap["name"] = "profiles.name"
	sortByMap["date"] = "profile_states.updated_at"

	if sortBy == "" || sortByMap[sortBy] == "" {
		return sortByMap["date"]
	}

	return sortByMap[sortBy]
}
