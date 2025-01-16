package advance_handler

import (
	"encoding/json"
	"fmt"

	"github.com/Mugen-Builders/to-do-sqlite/internal/domain"
	"github.com/Mugen-Builders/to-do-sqlite/internal/usecase"
	"github.com/Mugen-Builders/to-do-sqlite/pkg/rollups"
)

type TodoAdvanceHandlers struct {
	TodoRepository domain.TodoRepository
}

func NewTodoAdvanceHandlers(todoRepository domain.TodoRepository) *TodoAdvanceHandlers {
	return &TodoAdvanceHandlers{
		TodoRepository: todoRepository,
	}
}

func (h *TodoAdvanceHandlers) CreateTodoHandler(payload []byte, metadata rollups.Metadata) error {
	var input usecase.CreateTodoInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	createTodo := usecase.NewCreateTodoUseCase(h.TodoRepository)
	res, err := createTodo.Execute(&input, metadata)
	if err != nil {
		return err
	}
	todo, err := json.Marshal(res)
	if err != nil {
		return err
	}
	rollups.SendNotice(&rollups.NoticeRequest{
		Payload: rollups.Str2Hex(fmt.Sprintf("todo created - %s", todo)),
	})
	return nil
}

func (h *TodoAdvanceHandlers) UpdateTodoHandler(payload []byte, metadata rollups.Metadata) error {
	var input usecase.UpdateTodoInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	updateTodo := usecase.NewUpdateTodoUseCase(h.TodoRepository)
	res, err := updateTodo.Execute(&input, metadata)
	if err != nil {
		return err
	}
	todo, err := json.Marshal(res)
	if err != nil {
		return err
	}
	rollups.SendNotice(&rollups.NoticeRequest{
		Payload: rollups.Str2Hex(fmt.Sprintf("todo updated - %s", todo)),
	})
	return nil
}

func (h *TodoAdvanceHandlers) DeleteTodoHandler(payload []byte, metadata rollups.Metadata) error {
	var input usecase.DeleteTodoInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
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
		Payload: rollups.Str2Hex(fmt.Sprintf("todo deleted - %s", todo)),
	})
	return nil
}
