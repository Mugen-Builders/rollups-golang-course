package user

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
)

type FindUserByAddressInputDTO struct {
	Address Address `json:"address" validate:"required"`
}

type FindUserByAddressUseCase struct {
	UserRepository repository.UserRepository
}

func NewFindUserByAddressUseCase(userRepository repository.UserRepository) *FindUserByAddressUseCase {
	return &FindUserByAddressUseCase{
		UserRepository: userRepository,
	}
}

func (u *FindUserByAddressUseCase) Execute(ctx context.Context, input *FindUserByAddressInputDTO) (*FindUserOutputDTO, error) {
	res, err := u.UserRepository.FindUserByAddress(ctx, input.Address)
	if err != nil {
		return nil, err
	}
	return &FindUserOutputDTO{
		Id:             res.Id,
		Role:           string(res.Role),
		Address:        res.Address,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
	}, nil
}
