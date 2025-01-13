package repository

import (
	"fmt"

	"github.com/Mugen-Builders/to-do-sqlite/internal/domain"
	"gorm.io/gorm"
)

type TodoRepositorySQLite struct {
	Db *gorm.DB
}

func NewTodoRepositorySQLite(db *gorm.DB) *TodoRepositorySQLite {
	return &TodoRepositorySQLite{
		Db: db,
	}
}

func (r *TodoRepositorySQLite) CreateTodo(input *domain.Todo) (*domain.Todo, error) {
	if err := r.Db.Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}
	return input, nil
}

func (r *TodoRepositorySQLite) FindAllTodos() ([]*domain.Todo, error) {
	var todos []*domain.Todo
	if err := r.Db.Find(&todos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to find all todos: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to find all todos: %w", err)
	}
	return todos, nil
}

func (r *TodoRepositorySQLite) UpdateTodo(input *domain.Todo) (*domain.Todo, error) {
	if err := r.Db.Updates(&input).Error; err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	todo, err := r.findTodoById(input.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}
	return todo, nil
}

func (r *TodoRepositorySQLite) DeleteTodo(id uint) error {
	if err := r.Db.Delete(&domain.Todo{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to delete todo: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	return nil
}

func (r *TodoRepositorySQLite) findTodoById(id uint) (*domain.Todo, error) {
	var todo domain.Todo
	if err := r.Db.First(&todo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to find todo by id: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to find todo by id: %w", err)
	}
	return &todo, nil
}
