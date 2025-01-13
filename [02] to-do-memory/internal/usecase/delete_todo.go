package usecase

import "github.com/Mugen-Builders/to-do-memory/internal/domain"

type DeleteTodoInputDTO struct {
	Id uint `json:"id"`
}

type DeleteTodoUseCase struct {
	TodoRepository domain.TodoRepository
}

func NewDeleteTodoUseCase(todoRepository domain.TodoRepository) *DeleteTodoUseCase {
	return &DeleteTodoUseCase{
		TodoRepository: todoRepository,
	}
}

func (u *DeleteTodoUseCase) Execute(input *DeleteTodoInputDTO) error {
	return u.TodoRepository.DeleteTodo(input.Id)
}
