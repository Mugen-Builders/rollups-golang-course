package usecase

import (
	"github.com/Mugen-Builders/to-do-memory/internal/domain"
	"github.com/Mugen-Builders/to-do-memory/pkg/rollups"
)

type UpdateTodoInputDTO struct {
	Id          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type UpdateTodoOutputDTO struct {
	Id          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   uint64 `json:"created_at"`
	UpdatedAt   uint64 `json:"updated_at"`
}

type UpdateTodoUseCase struct {
	TodoRepository domain.TodoRepository
}

func NewUpdateTodoUseCase(todoRepository domain.TodoRepository) *UpdateTodoUseCase {
	return &UpdateTodoUseCase{
		TodoRepository: todoRepository,
	}
}

func (u *UpdateTodoUseCase) Execute(input *UpdateTodoInputDTO, metadata rollups.Metadata) (*UpdateTodoOutputDTO, error) {
	res, err := u.TodoRepository.UpdateTodo(&domain.Todo{
		Id:          input.Id,
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
		UpdatedAt:   metadata.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateTodoOutputDTO{
		Id:          res.Id,
		Title:       res.Title,
		Description: res.Description,
		Completed:   res.Completed,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}
