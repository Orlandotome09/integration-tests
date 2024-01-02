package entity

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DocumentsAtLeastOneHasFile(t *testing.T) {
	documentID1 := uuid.New()
	documentID2 := uuid.New()
	documentID3 := uuid.New()

	documents := Documents{
		{DocumentID: documentID1}, {DocumentID: documentID2}, {DocumentID: documentID3},
	}

	documentFiles := []DocumentFile{
		{DocumentID: documentID1},
	}

	atLeastOneHasFile := documents.AtLeastOneHasFile(documentFiles)

	assert.True(t, atLeastOneHasFile)
}

func Test_DocumentsDoNotHaveAnyFile(t *testing.T) {
	documentID1 := uuid.New()
	documentID2 := uuid.New()
	documentID3 := uuid.New()

	documents := Documents{
		{DocumentID: documentID1}, {DocumentID: documentID2}, {DocumentID: documentID3},
	}

	documentFiles := []DocumentFile{
		{DocumentID: uuid.New()},
	}

	atLeastOneHasFile := documents.AtLeastOneHasFile(documentFiles)

	assert.False(t, atLeastOneHasFile)
}
