package user

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
)

type FindUserByRoleInputDTO struct {
	Role string `json:"role" validate:"required"`
}

type FindUserByRoleOutputDTO []*FindUserOutputDTO

type FindUserByRoleUseCase struct {
	userRepository repository.UserRepository
}

func NewFindUserByRoleUseCase(userRepository repository.UserRepository) *FindUserByRoleUseCase {
	return &FindUserByRoleUseCase{
		userRepository: userRepository,
	}
}

func (u *FindUserByRoleUseCase) Execute(ctx context.Context, input *FindUserByRoleInputDTO) ([]*FindUserOutputDTO, error) {
	res, err := u.userRepository.FindUsersByRole(ctx, input.Role)
	if err != nil {
		return nil, err
	}
	output := make(FindUserByRoleOutputDTO, len(res))
	for i, user := range res {
		output[i] = &FindUserOutputDTO{
			Id:             user.Id,
			Role:           string(user.Role),
			Address:        user.Address,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
		}
	}
	return output, nil
}
