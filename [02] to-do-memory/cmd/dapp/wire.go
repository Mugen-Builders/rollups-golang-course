//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Mugen-Builders/to-do-memory/configs"
	"github.com/Mugen-Builders/to-do-memory/internal/domain"
	"github.com/Mugen-Builders/to-do-memory/internal/infra/cartesi/advance_handler"
	"github.com/Mugen-Builders/to-do-memory/internal/infra/cartesi/inspect_handler"
	"github.com/Mugen-Builders/to-do-memory/internal/infra/repository"
	"github.com/google/wire"
)

var setToDoRepositoryDependency = wire.NewSet(
	repository.NewToDoRepositoryInMemory,
	wire.Bind(new(domain.ToDoRepository), new(*repository.ToDoRepositoryInMemory)),
)


var setAdvanceHandlers = wire.NewSet(
	advance_handler.NewToDoAdvanceHandlers,
)

var setInspectHandlers = wire.NewSet(
	inspect_handler.NewToDoInspectHandlers,
)

func NewAdvanceHandlers(db *configs.InMemoryDB) (*AdvanceHandlers, error) {
	wire.Build(
		setToDoRepositoryDependency,
		setAdvanceHandlers,
		wire.Struct(new(AdvanceHandlers), "*"),
	)
	return nil, nil
}

func NewInspectHandlers(db *configs.InMemoryDB) (*InspectHandlers, error) {
	wire.Build(
		setToDoRepositoryDependency,
		setInspectHandlers,
		wire.Struct(new(InspectHandlers), "*"),
	)
	return nil, nil
}

type AdvanceHandlers struct {
	ToDoAdvanceHandlers *advance_handler.ToDoAdvanceHandlers
}

type InspectHandlers struct {
	ToDoInspectHandlers *inspect_handler.ToDoInspectHandlers
}