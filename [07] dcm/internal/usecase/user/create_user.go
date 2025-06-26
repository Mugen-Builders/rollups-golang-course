package user

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CreateUserInputDTO struct {
	Role    string  `json:"role" validate:"required"`
	Address Address `json:"address" validate:"required"`
}

type CreateUserOutputDTO struct {
	Id              uint                    `json:"id"`
	Role            string                  `json:"role"`
	Address         Address                 `json:"address"`
	InvestmentLimit *uint256.Int            `json:"investment_limit,omitempty" gorm:"type:bigint"`
	CreatedAt       int64                   `json:"created_at"`
}

type CreateUserUseCase struct {
	UserRepository repository.UserRepository
}

func NewCreateUserUseCase(userRepository repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		UserRepository: userRepository,
	}
}

func (u *CreateUserUseCase) Execute(ctx context.Context, input *CreateUserInputDTO, metadata rollmelette.Metadata) (*CreateUserOutputDTO, error) {
	user, err := entity.NewUser(input.Role, input.Address, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}

	res, err := u.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &CreateUserOutputDTO{
		Id:             res.Id,
		Role:           string(res.Role),
		Address:        res.Address,
		CreatedAt:      res.CreatedAt,
	}, nil
}
