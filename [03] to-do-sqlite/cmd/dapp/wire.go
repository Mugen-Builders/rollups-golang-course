//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Mugen-Builders/to-do-sqlite/internal/domain"
	"github.com/Mugen-Builders/to-do-sqlite/internal/infra/cartesi/advance_handler"
	"github.com/Mugen-Builders/to-do-sqlite/internal/infra/cartesi/inspect_handler"
	"github.com/Mugen-Builders/to-do-sqlite/internal/infra/repository"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var setTodoRepositoryDependency = wire.NewSet(
	repository.NewTodoRepositorySQLite,
	wire.Bind(new(domain.TodoRepository), new(*repository.TodoRepositorySQLite)),
)

var setAdvanceHandlers = wire.NewSet(
	advance_handler.NewTodoAdvanceHandlers,
)

var setInspectHandlers = wire.NewSet(
	inspect_handler.NewTodoInspectHandlers,
)

func NewAdvanceHandlers(gormDB *gorm.DB) (*AdvanceHandlers, error) {
	wire.Build(
		setTodoRepositoryDependency,
		setAdvanceHandlers,
		wire.Struct(new(AdvanceHandlers), "*"),
	)
	return nil, nil
}

func NewInspectHandlers(gormDB *gorm.DB) (*InspectHandlers, error) {
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
