package configs

import (
	"sync"

	"github.com/Mugen-Builders/to-do-memory/internal/domain"
)

type InMemoryDB struct {
	Todos map[uint]*domain.Todo
	Lock  sync.RWMutex
}

func SetupInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		Todos: make(map[uint]*domain.Todo),
	}
}
