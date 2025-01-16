package usecase

import (
    "github.com/Mugen-Builders/to-do-sqlite/internal/domain"
)

type FindTodoOutputDTO struct {
    Id          uint   `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
    CreatedAt   uint64 `json:"created_at"`
    UpdatedAt   uint64 `json:"updated_at"`
}

type FindAllTodosOutputDTO []*FindTodoOutputDTO

type FindAllTodosUseCase struct {
    TodoRepository domain.TodoRepository
}

func NewFindAllTodosUseCase(todoRepository domain.TodoRepository) *FindAllTodosUseCase {
    return &FindAllTodosUseCase{
        TodoRepository: todoRepository,
    }
}

func (u *FindAllTodosUseCase) Execute() (*FindAllTodosOutputDTO, error) {
    res, err := u.TodoRepository.FindAllTodos()
    if err != nil {
        return nil, err
    }
    output := make(FindAllTodosOutputDTO, len(res))
    for i, todo := range res {
        output[i] = &FindTodoOutputDTO{
            Id:          todo.Id,
            Title:       todo.Title,
            Description: todo.Description,
            Completed:   todo.Completed,
            CreatedAt:   todo.CreatedAt,
            UpdatedAt:   todo.UpdatedAt,
        }
    }
    return &output, nil
}
