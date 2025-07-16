package advance

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/henriquemarlon/cartesi-golang-series/to-do/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/to-do/internal/usecase"
	"github.com/henriquemarlon/cartesi-golang-series/to-do/pkg/rollups"
)

type ToDoAdvanceHandlers struct {
	ToDoRepository repository.ToDoRepository
}

func NewToDoAdvanceHandlers(toDoRepository repository.ToDoRepository) *ToDoAdvanceHandlers {
	return &ToDoAdvanceHandlers{
		ToDoRepository: toDoRepository,
	}
}

func (h *ToDoAdvanceHandlers) CreateToDoHandler(payload []byte, metadata rollups.Metadata) error {
	var input usecase.CreateToDoInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
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

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
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

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
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
