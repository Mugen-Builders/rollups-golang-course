package repository

import (
	"sync"

	"github.com/Mugen-Builders/to-do-memory/internal/domain"
)

type TodoRepositoryInMemory struct {
	db     map[uint]*domain.Todo
	mutex  *sync.RWMutex
	nextID uint
}

func NewTodoRepositoryInMemory() *TodoRepositoryInMemory {
	return &TodoRepositoryInMemory{
		db:     make(map[uint]*domain.Todo),
		mutex:  &sync.RWMutex{},
		nextID: 1,
	}
}

func (r *TodoRepositoryInMemory) CreateTodo(input *domain.Todo) (*domain.Todo, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	input.Id = r.nextID
	r.nextID++
	r.db[input.Id] = input
	return input, nil
}

func (r *TodoRepositoryInMemory) FindAllTodos() ([]*domain.Todo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var todos []*domain.Todo
	for _, todo := range r.db {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepositoryInMemory) UpdateTodo(input *domain.Todo) (*domain.Todo, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	todo, exists := r.db[input.Id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	todo.Title = input.Title
	todo.Description = input.Description
	todo.Completed = input.Completed

	r.db[input.Id] = todo
	
	return todo, nil
}

func (r *TodoRepositoryInMemory) DeleteTodo(id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.db[id]
	if !exists {
		return domain.ErrNotFound
	}

	delete(r.db, id)
	return nil
}
