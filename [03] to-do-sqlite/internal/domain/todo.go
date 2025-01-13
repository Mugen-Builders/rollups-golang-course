package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidTodo = errors.New("invalid todo")
	ErrNotFound    = errors.New("todo not found")
)

type TodoRepository interface {
	CreateTodo(todo *Todo) (*Todo, error)
	FindAllTodos() ([]*Todo, error)
	UpdateTodo(todo *Todo) (*Todo, error)
	DeleteTodo(id uint) error
}

type Todo struct {
	Id          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"type:text;not null"`
	Description string `json:"description" gorm:"type:text;not null"`
	Completed   bool   `json:"completed" gorm:"default:false"`
	CreatedAt   uint64 `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt   uint64 `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewTodo(title string, description string, createdAt uint64) (*Todo, error) {
	todo := &Todo{
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
	}
	if err := todo.Validate(); err != nil {
		return nil, err
	}
	return todo, nil
}

func (t *Todo) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("%w: title cannot be empty", ErrInvalidTodo)
	}
	if t.Description == "" {
		return fmt.Errorf("%w: description cannot be empty", ErrInvalidTodo)
	}
	return nil
}
