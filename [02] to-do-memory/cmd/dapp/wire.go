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

var setTodoRepositoryDependency = wire.NewSet(
	repository.NewTodoRepositoryInMemory,
	wire.Bind(new(domain.TodoRepository), new(*repository.TodoRepositoryInMemory)),
)

var setAdvanceHandlers = wire.NewSet(
	advance_handler.NewTodoAdvanceHandlers,
)

var setInspectHandlers = wire.NewSet(
	inspect_handler.NewTodoInspectHandlers,
)

func NewAdvanceHandlers(db *configs.InMemoryDB) (*AdvanceHandlers, error) {
	wire.Build(
		setTodoRepositoryDependency,
		setAdvanceHandlers,
		wire.Struct(new(AdvanceHandlers), "*"),
	)
	return nil, nil
}

func NewInspectHandlers(db *configs.InMemoryDB) (*InspectHandlers, error) {
	wire.Build(
		setTodoRepositoryDependency,
		setInspectHandlers,
		wire.Struct(new(InspectHandlers), "*"),
	)
	return nil, nil
}

type AdvanceHandlers struct {
	TodoAdvanceHandlers *advance_handler.TodoAdvanceHandlers
}

type InspectHandlers struct {
	TodoInspectHandlers *inspect_handler.TodoInspectHandlers
}
