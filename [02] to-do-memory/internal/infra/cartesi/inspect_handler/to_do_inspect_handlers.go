package inspect_handler

import (
	"encoding/json"

	"github.com/Mugen-Builders/to-do-memory/internal/domain"
	"github.com/Mugen-Builders/to-do-memory/internal/usecase"
	"github.com/Mugen-Builders/to-do-memory/pkg/rollups"
)

type ToDoInspectHandlers struct {
	ToDoRepository domain.ToDoRepository
}

func NewToDoInspectHandlers(toDoRepository domain.ToDoRepository) *ToDoInspectHandlers {
	return &ToDoInspectHandlers{
		ToDoRepository: toDoRepository,
	}
}

func (h *ToDoInspectHandlers) FindAllToDosHandler() error {
	findAllToDos := usecase.NewFindAllToDosUseCase(h.ToDoRepository)
	res, err := findAllToDos.Execute()
	if err != nil {
		return err
	}
	toDos, err := json.Marshal(res)
	if err != nil {
		return err
	}
	rollups.SendReport(&rollups.ReportRequest{
		Payload: rollups.Str2Hex(string(toDos)),
	})
	return nil
}
