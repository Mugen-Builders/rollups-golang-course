package user

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
)

type DeleteUserInputDTO struct {
	Address Address `json:"address" validate:"required"`
}

type DeleteUserUseCase struct {
	UserRepository repository.UserRepository
}

func NewDeleteUserUseCase(userRepository repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		UserRepository: userRepository,
	}
}

func (u *DeleteUserUseCase) Execute(ctx context.Context, input *DeleteUserInputDTO) error {
	return u.UserRepository.DeleteUser(ctx, input.Address)
}
