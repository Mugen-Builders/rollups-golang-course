package usecase

import "github.com/Mugen-Builders/to-do-sqlite/internal/domain"

type DeleteToDoInputDTO struct {
	Id uint `json:"id"`
}

type DeleteToDoUseCase struct {
	ToDoRepository domain.ToDoRepository
}

func NewDeleteToDoUseCase(todoRepository domain.ToDoRepository) *DeleteToDoUseCase {
	return &DeleteToDoUseCase{
		ToDoRepository: todoRepository,
	}
}

func (u *DeleteToDoUseCase) Execute(input *DeleteToDoInputDTO) error {
	return u.ToDoRepository.DeleteToDo(input.Id)
}
