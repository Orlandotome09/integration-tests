package entity

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnrichedInformationFilterCafProvider(t *testing.T) {
	requestID := uuid.New().String()
	enricheInformation := EnrichedInformation{
		Providers: []Provider{
			{
				ProviderName: ProviderCAF,
				RequestID:    requestID,
				Status:       "APPROVED",
			},
		},
	}

	expected := &Provider{
		ProviderName: ProviderCAF,
		RequestID:    requestID,
		Status:       "APPROVED",
	}

	received := enricheInformation.FilterCafProvider()

	assert.Equal(t, expected, received)
}

func TestEnrichedInformationNotFilterCafProvider(t *testing.T) {
	requestID := uuid.New().String()
	enricheInformation := EnrichedInformation{
		Providers: []Provider{
			{
				ProviderName: "BUREAU",
				RequestID:    requestID,
				Status:       "APPROVED",
			},
		},
	}

	var expected *Provider

	received := enricheInformation.FilterCafProvider()

	assert.Equal(t, expected, received)
}
