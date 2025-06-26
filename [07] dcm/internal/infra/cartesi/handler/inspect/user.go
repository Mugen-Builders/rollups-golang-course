package inspect

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/usecase/user"
	"github.com/rollmelette/rollmelette"
)

type UserInspectHandlers struct {
	UserRepository repository.UserRepository
}

func NewUserInspectHandlers(userRepository repository.UserRepository) *UserInspectHandlers {
	return &UserInspectHandlers{
		UserRepository: userRepository,
	}
}

func (h *UserInspectHandlers) FindUserByAddress(env rollmelette.EnvInspector, payload []byte) error {
	var input user.FindUserByAddressInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findUserByAddress := user.NewFindUserByAddressUseCase(h.UserRepository)
	res, err := findUserByAddress.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find User: %w", err)
	}
	User, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal User: %w", err)
	}
	env.Report(User)
	return nil
}

func (h *UserInspectHandlers) FindAllUsers(env rollmelette.EnvInspector, payload []byte) error {
	ctx := context.Background()
	findAllUsers := user.NewFindAllUsersUseCase(h.UserRepository)
	res, err := findAllUsers.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all Users: %w", err)
	}
	allUsers, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all Users: %w", err)
	}
	env.Report(allUsers)
	return nil
}

func (h *UserInspectHandlers) ERC20BalanceOf(env rollmelette.EnvInspector, payload []byte) error {
	var input user.BalanceOfInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findUserByAddress := user.NewFindUserByAddressUseCase(h.UserRepository)
	res, err := findUserByAddress.Execute(ctx, &user.FindUserByAddressInputDTO{
		Address: input.Address,
	})
	if err != nil {
		return fmt.Errorf("failed to find User: %w", err)
	}

	balance := env.ERC20BalanceOf(
		common.Address(input.Token),
		common.Address(res.Address),
	).String()

	balanceBytes, err := json.Marshal(balance)
	if err != nil {
		return fmt.Errorf("failed to marshal balance: %w", err)
	}

	env.Report(balanceBytes)
	return nil
}
