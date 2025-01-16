package configs

import (
	"sync"

	"github.com/Mugen-Builders/to-do-memory/internal/domain"
)

type InMemoryDB struct {
	ToDos map[uint]*domain.ToDo
	Lock  *sync.RWMutex
}

func SetupInMemoryDB() (*InMemoryDB, error) {
	return &InMemoryDB{
		ToDos: make(map[uint]*domain.ToDo),
		Lock:  &sync.RWMutex{},
	}, nil
}