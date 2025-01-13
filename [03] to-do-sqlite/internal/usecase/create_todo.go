package usecase

import (
	"github.com/Mugen-Builders/to-do-sqlite/internal/domain"
	"github.com/Mugen-Builders/to-do-sqlite/pkg/rollups"
)

type CreateTodoInputDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

type CreateTodoOutputDTO struct {
	Id          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   uint64 `json:"created_at"`
}

type CreateTodoUseCase struct {
	TodoRepository domain.TodoRepository
}

func NewCreateTodoUseCase(todoRepository domain.TodoRepository) *CreateTodoUseCase {
	return &CreateTodoUseCase{
		TodoRepository: todoRepository,
	}
}

func (u *CreateTodoUseCase) Execute(input *CreateTodoInputDTO, metadata rollups.Metadata) (*CreateTodoOutputDTO, error) {
	res, err := domain.NewTodo(input.Title, input.Description, metadata.Timestamp)
	if err != nil {
		return nil, err
	}

	res, err = u.TodoRepository.CreateTodo(res)
	if err != nil {
		return nil, err
	}

	return &CreateTodoOutputDTO{
		Id:          res.Id,
		Title:       res.Title,
		Description: res.Description,
		Completed:   res.Completed,
		CreatedAt:   res.CreatedAt,
	}, nil
}
