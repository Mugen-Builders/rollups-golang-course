package sqlite

import (
	"fmt"

	"github.com/henriquemarlon/to-do/internal/domain"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateToDo(input *domain.ToDo) (*domain.ToDo, error) {
	if err := r.Db.Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create to-do: %w", err)
	}
	return input, nil
}

func (r *SQLiteRepository) FindAllToDos() ([]*domain.ToDo, error) {
	var toDos []*domain.ToDo
	if err := r.Db.Find(&toDos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to find all to-dos: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to find all to-dos: %w", err)
	}
	return toDos, nil
}

func (r *SQLiteRepository) UpdateToDo(input *domain.ToDo) (*domain.ToDo, error) {
	if err := r.Db.Updates(&input).Error; err != nil {
		return nil, fmt.Errorf("failed to update to-do: %w", err)
	}
	toDo, err := r.findToDoById(input.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to update to-do: %w", err)
	}
	return toDo, nil
}

func (r *SQLiteRepository) DeleteToDo(id uint) error {
	if err := r.Db.Delete(&domain.ToDo{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to delete to-do: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("failed to delete to-do: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) findToDoById(id uint) (*domain.ToDo, error) {
	var toDo domain.ToDo
	if err := r.Db.First(&toDo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to find to-do by id: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to find to-do by id: %w", err)
	}
	return &toDo, nil
}
