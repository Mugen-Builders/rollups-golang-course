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

var setToDoRepositoryDependency = wire.NewSet(
	repository.NewToDoRepositorySQLite,
	wire.Bind(new(domain.ToDoRepository), new(*repository.ToDoRepositorySQLite)),
)

var setAdvanceHandlers = wire.NewSet(
	advance_handler.NewToDoAdvanceHandlers,
)

var setInspectHandlers = wire.NewSet(
	inspect_handler.NewToDoInspectHandlers,
)

func NewAdvanceHandlers(gormDB *gorm.DB) (*AdvanceHandlers, error) {
	wire.Build(
		setToDoRepositoryDependency,
		setAdvanceHandlers,
		wire.Struct(new(AdvanceHandlers), "*"),
	)
	return nil, nil
}

func NewInspectHandlers(gormDB *gorm.DB) (*InspectHandlers, error) {
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
