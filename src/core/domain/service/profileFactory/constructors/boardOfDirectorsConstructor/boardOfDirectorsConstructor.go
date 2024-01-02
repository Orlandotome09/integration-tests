package boardOfDirectorsConstructor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/pkg/errors"
)

type boardOfDirectorsConstructor struct {
	boardOfDirectorsAdapter interfaces.BoardOfDirectorsAdapter
}

func New(boardOfDirectorsService interfaces.BoardOfDirectorsAdapter) interfaces.ProfileConstructor {
	return &boardOfDirectorsConstructor{boardOfDirectorsAdapter: boardOfDirectorsService}
}

func (ref *boardOfDirectorsConstructor) Assemble(profileWrapper *entity.ProfileWrapper) error {

	if !profileWrapper.Profile.ShouldGetBoardOfDirectors() {
		return nil
	}

	directors, err := ref.boardOfDirectorsAdapter.Search(*profileWrapper.Profile.ProfileID)
	if err != nil {
		return errors.WithStack(err)
	}

	profileWrapper.Mutex.Lock()
	defer profileWrapper.Mutex.Unlock()
	profileWrapper.Profile.BoardOfDirectors = directors

	return nil

}
