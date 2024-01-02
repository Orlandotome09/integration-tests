package mutex

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"sync"
)

type mutex struct {
	mu   sync.Mutex
	list map[string]bool
}

func New() _interfaces.Mutex {
	return &mutex{}
}

func (ref *mutex) Lock(id string) (success bool) {
	ref.mu.Lock()
	defer ref.mu.Unlock()

	if ref.list == nil {
		ref.list = make(map[string]bool, 0)
	}

	if _, exists := ref.list[id]; !exists {
		ref.list[id] = true
		return true
	}

	return false
}

func (ref *mutex) Release(processingKey string) {
	ref.mu.Lock()
	defer ref.mu.Unlock()

	delete(ref.list, processingKey)
}
