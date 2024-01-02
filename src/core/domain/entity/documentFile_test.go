package entity

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDocumentFilesHaveDocument(t *testing.T) {
	documentID := uuid.New()
	documentFiles := DocumentFiles{
		{DocumentID: uuid.New()}, {DocumentID: uuid.New()}, {DocumentID: documentID},
	}

	haveDocument := documentFiles.HaveDocument(documentID)

	assert.True(t, haveDocument)
}

func TestDocumentFilesDoNotHaveDocument(t *testing.T) {
	documentID := uuid.New()
	documentFiles := DocumentFiles{
		{DocumentID: uuid.New()}, {DocumentID: uuid.New()}, {DocumentID: uuid.New()},
	}

	haveDocument := documentFiles.HaveDocument(documentID)

	assert.False(t, haveDocument)
}
