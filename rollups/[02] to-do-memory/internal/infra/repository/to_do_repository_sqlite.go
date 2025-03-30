package repository

import (
	"sync"

	"github.com/Mugen-Builders/to-do-memory/configs"
	"github.com/Mugen-Builders/to-do-memory/internal/domain"
)

type ToDoRepositoryInMemory struct {
	db     map[uint]*domain.ToDo
	mutex  *sync.RWMutex
	nextID uint
}

func NewToDoRepositoryInMemory(db *configs.InMemoryDB) *ToDoRepositoryInMemory {
	return &ToDoRepositoryInMemory{
		db:     db.ToDos,
		mutex:  db.Lock,
		nextID: 1,
	}
}

func (r *ToDoRepositoryInMemory) CreateToDo(input *domain.ToDo) (*domain.ToDo, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	input.Id = r.nextID
	r.nextID++
	r.db[input.Id] = input
	return input, nil
}

func (r *ToDoRepositoryInMemory) FindAllToDos() ([]*domain.ToDo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var todos []*domain.ToDo
	for _, todo := range r.db {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *ToDoRepositoryInMemory) UpdateToDo(input *domain.ToDo) (*domain.ToDo, error) {
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

func (r *ToDoRepositoryInMemory) DeleteToDo(id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	_, exists := r.db[id]
	if !exists {
		return domain.ErrNotFound
	}

	delete(r.db, id)
	return nil
}
