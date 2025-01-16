package advance_handler

import (
	"encoding/json"
	"fmt"

	"github.com/Mugen-Builders/to-do-sqlite/internal/domain"
	"github.com/Mugen-Builders/to-do-sqlite/internal/usecase"
	"github.com/Mugen-Builders/to-do-sqlite/pkg/rollups"
)

type ToDoAdvanceHandlers struct {
	ToDoRepository domain.ToDoRepository
}

func NewToDoAdvanceHandlers(toDoRepository domain.ToDoRepository) *ToDoAdvanceHandlers {
	return &ToDoAdvanceHandlers{
		ToDoRepository: toDoRepository,
	}
}

func (h *ToDoAdvanceHandlers) CreateToDoHandler(payload []byte, metadata rollups.Metadata) error {
	var input usecase.CreateToDoInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	createToDo := usecase.NewCreateToDoUseCase(h.ToDoRepository)
	res, err := createToDo.Execute(&input, metadata)
	if err != nil {
		return err
	}
	toDo, err := json.Marshal(res)
	if err != nil {
		return err
	}
	rollups.SendNotice(&rollups.NoticeRequest{
		Payload: rollups.Str2Hex(fmt.Sprintf("To-Do created - %s", toDo)),
	})
	return nil
}

func (h *ToDoAdvanceHandlers) UpdateToDoHandler(payload []byte, metadata rollups.Metadata) error {
	var input usecase.UpdateToDoInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	updateToDo := usecase.NewUpdateToDoUseCase(h.ToDoRepository)
	res, err := updateToDo.Execute(&input, metadata)
	if err != nil {
		return err
	}
	toDo, err := json.Marshal(res)
	if err != nil {
		return err
	}
	rollups.SendNotice(&rollups.NoticeRequest{
		Payload: rollups.Str2Hex(fmt.Sprintf("To-Do updated - %s", toDo)),
	})
	return nil
}

func (h *ToDoAdvanceHandlers) DeleteToDoHandler(payload []byte, metadata rollups.Metadata) error {
	var input usecase.DeleteToDoInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	deleteToDo := usecase.NewDeleteToDoUseCase(h.ToDoRepository)
	err := deleteToDo.Execute(&input)
	if err != nil {
		return err
	}
	toDo, err := json.Marshal(input)
	if err != nil {
		return err
	}
	rollups.SendNotice(&rollups.NoticeRequest{
		Payload: rollups.Str2Hex(fmt.Sprintf("To-Do deleted - %s", toDo)),
	})
	return nil
}
