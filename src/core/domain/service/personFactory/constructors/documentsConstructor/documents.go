package documentsConstructor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type documentsPersonConstructor struct {
	documentAdapter     interfaces.DocumentAdapter
	documentFileAdapter interfaces.DocumentFileAdapter
}

func New(documentAdapter interfaces.DocumentAdapter, documentFileAdapter interfaces.DocumentFileAdapter) interfaces.PersonConstructor {
	return &documentsPersonConstructor{
		documentAdapter:     documentAdapter,
		documentFileAdapter: documentFileAdapter,
	}
}

func (ref *documentsPersonConstructor) Assemble(personWrapper *entity.PersonWrapper) error {
	if !personWrapper.Person.ShouldGetDocuments() {
		return nil
	}

	documents, err := ref.documentAdapter.Find(personWrapper.Person.EntityID.String())
	if err != nil {
		return errors.WithStack(err)
	}

	documentFiles := make([]entity.DocumentFile, 0)

	for _, document := range documents {
		files, err := ref.documentFileAdapter.FindByDocumentID(document.DocumentID)
		if err != nil {
			return errors.WithStack(err)
		}
		documentFiles = append(documentFiles, files...)
	}

	personWrapper.Mutex.Lock()
	defer personWrapper.Mutex.Unlock()
	personWrapper.Person.Documents = documents
	personWrapper.Person.DocumentFiles = documentFiles

	return nil
}
