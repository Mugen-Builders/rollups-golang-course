package advance_handler

import (
	"encoding/json"
	"fmt"

	"github.com/Mugen-Builders/to-do-memory/internal/domain"
	"github.com/Mugen-Builders/to-do-memory/internal/usecase"
	"github.com/Mugen-Builders/to-do-memory/pkg/rollups"
)

type TodoAdvanceHandlers struct {
	TodoRepository domain.TodoRepository
}

func NewTodoAdvanceHandlers(todoRepository domain.TodoRepository) *TodoAdvanceHandlers {
	return &TodoAdvanceHandlers{
		TodoRepository: todoRepository,
	}
}

func (h *TodoAdvanceHandlers) CreateTodoHandler(data *rollups.AdvanceResponse) error {
	var input usecase.CreateTodoInputDTO
	if err := json.Unmarshal(data.Payload, &input); err != nil {
		return err
	}
	createTodo := usecase.NewCreateTodoUseCase(h.TodoRepository)
	res, err := createTodo.Execute(&input, data.Metadata)
	if err != nil {
		return err
	}
	todo, err := json.Marshal(res)
	if err != nil {
		return err
	}
	rollups.SendNotice(&rollups.NoticeRequest{
		Payload: fmt.Sprintf("todo created - %v", todo),
	})
	return nil
}

func (h *TodoAdvanceHandlers) UpdateTodoHandler(data *rollups.AdvanceResponse) error {
	var input usecase.UpdateTodoInputDTO
	if err := json.Unmarshal(data.Payload, &input); err != nil {
		return err
	}
	updateTodo := usecase.NewUpdateTodoUseCase(h.TodoRepository)
	res, err := updateTodo.Execute(&input, data.Metadata)
	if err != nil {
		return err
	}
	todo, err := json.Marshal(res)
	if err != nil {
		return err
	}
	rollups.SendNotice(&rollups.NoticeRequest{
		Payload: fmt.Sprintf("todo updated - %v", todo),
	})
	return nil
}

func (h *TodoAdvanceHandlers) DeleteTodoHandler(data *rollups.AdvanceResponse) error {
	var input usecase.DeleteTodoInputDTO
	if err := json.Unmarshal(data.Payload, &input); err != nil {
		return err
	}
	deleteTodo := usecase.NewDeleteTodoUseCase(h.TodoRepository)
	err := deleteTodo.Execute(&input)
	if err != nil {
		return err
	}
	todo, err := json.Marshal(input)
	if err != nil {
		return err
	}
	rollups.SendNotice(&rollups.NoticeRequest{
		Payload: fmt.Sprintf("todo deleted - %v", todo),
	})
	return nil
}
