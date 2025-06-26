package middleware

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/usecase/user"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/router"
	"github.com/rollmelette/rollmelette"
)

type RBACFactory struct {
	userRepository repository.UserRepository
}

func NewRBACFactory(userRepository repository.UserRepository) *RBACFactory {
	return &RBACFactory{
		userRepository: userRepository,
	}
}

func (f *RBACFactory) Create(roles []string) router.Middleware {
	return func(handler any) any {
		switch h := handler.(type) {
		case router.AdvanceHandlerFunc:
			return router.AdvanceHandlerFunc(func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
				var address Address
				ctx := context.Background()

				// Get the sender address from either ERC20 deposit or metadata
				erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
				if ok {
					address = Address(erc20Deposit.Sender)
				} else {
					address = Address(metadata.MsgSender)
				}

				// Find user and check roles
				findUserByAddress := user.NewFindUserByAddressUseCase(f.userRepository)
				user, err := findUserByAddress.Execute(ctx, &user.FindUserByAddressInputDTO{
					Address: address,
				})
				if err != nil {
					return err
				}

				// Check if user has any of the required roles
				var hasRole bool
				for _, role := range roles {
					if user.Role == role {
						hasRole = true
						break
					}
				}
				if !hasRole {
					return fmt.Errorf("user %s lacks required permissions: %v", common.Address(user.Address), roles)
				}

				return h(env, metadata, deposit, payload)
			})
		case router.InspectHandlerFunc:
			return h
		default:
			return handler
		}
	}
}

func (f *RBACFactory) AdminOnly() router.Middleware {
	return f.Create([]string{"admin"})
}

func (f *RBACFactory) InvestorOnly() router.Middleware {
	return f.Create([]string{"investor"})
}

func (f *RBACFactory) IssuerOnly() router.Middleware {
	return f.Create([]string{"issuer"})
}
