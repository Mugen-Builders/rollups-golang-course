package usecase

import (
	"github.com/Mugen-Builders/to-do-sqlite/internal/domain"
	"github.com/Mugen-Builders/to-do-sqlite/pkg/rollups"
)

type CreateToDoInputDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateToDoOutputDTO struct {
	Id          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   uint64 `json:"created_at"`
}

type CreateToDoUseCase struct {
	ToDoRepository domain.ToDoRepository
}

func NewCreateToDoUseCase(todoRepository domain.ToDoRepository) *CreateToDoUseCase {
	return &CreateToDoUseCase{
		ToDoRepository: todoRepository,
	}
}

func (u *CreateToDoUseCase) Execute(input *CreateToDoInputDTO, metadata rollups.Metadata) (*CreateToDoOutputDTO, error) {
	res, err := domain.NewToDo(input.Title, input.Description, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}

	res, err = u.ToDoRepository.CreateToDo(res)
	if err != nil {
		return nil, err
	}

	return &CreateToDoOutputDTO{
		Id:          res.Id,
		Title:       res.Title,
		Description: res.Description,
		Completed:   res.Completed,
		CreatedAt:   res.CreatedAt,
	}, nil
}
