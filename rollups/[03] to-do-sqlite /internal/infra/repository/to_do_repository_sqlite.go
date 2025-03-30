package repository

import (
	"fmt"

	"github.com/Mugen-Builders/to-do-sqlite/internal/domain"
	"gorm.io/gorm"
)

type ToDoRepositorySQLite struct {
	Db *gorm.DB
}

func NewToDoRepositorySQLite(db *gorm.DB) *ToDoRepositorySQLite {
	return &ToDoRepositorySQLite{
		Db: db,
	}
}

func (r *ToDoRepositorySQLite) CreateToDo(input *domain.ToDo) (*domain.ToDo, error) {
	if err := r.Db.Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create to-do: %w", err)
	}
	return input, nil
}

func (r *ToDoRepositorySQLite) FindAllToDos() ([]*domain.ToDo, error) {
	var toDos []*domain.ToDo
	if err := r.Db.Find(&toDos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to find all to-dos: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to find all to-dos: %w", err)
	}
	return toDos, nil
}

func (r *ToDoRepositorySQLite) UpdateToDo(input *domain.ToDo) (*domain.ToDo, error) {
	if err := r.Db.Updates(&input).Error; err != nil {
		return nil, fmt.Errorf("failed to update to-do: %w", err)
	}
	toDo, err := r.findToDoById(input.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to update to-do: %w", err)
	}
	return toDo, nil
}

func (r *ToDoRepositorySQLite) DeleteToDo(id uint) error {
	if err := r.Db.Delete(&domain.ToDo{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to delete to-do: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("failed to delete to-do: %w", err)
	}
	return nil
}

func (r *ToDoRepositorySQLite) findToDoById(id uint) (*domain.ToDo, error) {
	var toDo domain.ToDo
	if err := r.Db.First(&toDo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to find to-do by id: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to find to-do by id: %w", err)
	}
	return &toDo, nil
}
