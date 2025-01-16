package inspect_handler

import (
	"encoding/json"

	"github.com/Mugen-Builders/to-do-sqlite/internal/domain"
	"github.com/Mugen-Builders/to-do-sqlite/internal/usecase"
	"github.com/Mugen-Builders/to-do-sqlite/pkg/rollups"
)

type TodoInspectHandlers struct {
	TodoRepository domain.TodoRepository
}

func NewTodoInspectHandlers(todoRepository domain.TodoRepository) *TodoInspectHandlers {
	return &TodoInspectHandlers{
		TodoRepository: todoRepository,
	}
}

func (h *TodoInspectHandlers) FindAllTodosHandler() error {
	findAllTodos := usecase.NewFindAllTodosUseCase(h.TodoRepository)
	res, err := findAllTodos.Execute()
	if err != nil {
		return err
	}
	todos, err := json.Marshal(res)
	if err != nil {
		return err
	}
	rollups.SendReport(&rollups.ReportRequest{
		Payload: rollups.Str2Hex(string(todos)),
	})
	return nil
}
