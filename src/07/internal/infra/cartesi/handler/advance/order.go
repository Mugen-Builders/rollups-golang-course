package advance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/usecase/order"
	"github.com/rollmelette/rollmelette"
)

type OrderAdvanceHandlers struct {
	OrderRepository    repository.OrderRepository
	UserRepository     repository.UserRepository
	CampaignRepository repository.CampaignRepository
}

func NewOrderAdvanceHandlers(
	orderRepository repository.OrderRepository,
	userRepository repository.UserRepository,
	campaignRepository repository.CampaignRepository,
) *OrderAdvanceHandlers {
	return &OrderAdvanceHandlers{
		OrderRepository:    orderRepository,
		UserRepository:     userRepository,
		CampaignRepository: campaignRepository,
	}
}

func (h *OrderAdvanceHandlers) CreateOrder(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input order.CreateOrderInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w, payload: %s", err, string(payload))
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	createOrder := order.NewCreateOrderUseCase(
		h.OrderRepository,
		h.CampaignRepository,
	)

	res, err := createOrder.Execute(ctx, &input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	erc20Deposit, ok := deposit.(*rollmelette.ERC20Deposit)
	if !ok {
		return fmt.Errorf("invalid deposit custom_type, expected ERC20Deposit")
	}

	if err := env.ERC20Transfer(
		erc20Deposit.Token,
		erc20Deposit.Sender,
		env.AppAddress(),
		erc20Deposit.Value,
	); err != nil {
		return fmt.Errorf("failed to transfer ERC20: %w", err)
	}

	order, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("order created - "), order...))
	return nil
}

func (h *OrderAdvanceHandlers) CancelOrder(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input order.CancelOrderInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	cancelOrder := order.NewCancelOrderUseCase(
		h.OrderRepository,
		h.CampaignRepository,
	)

	res, err := cancelOrder.Execute(ctx, &input, metadata)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	if err := env.ERC20Transfer(
		common.Address(res.Token),
		env.AppAddress(),
		metadata.MsgSender,
		res.Amount.ToBig(),
	); err != nil {
		return fmt.Errorf("failed to transfer ERC20: %w", err)
	}

	order, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("order canceled - "), order...))
	return nil
}
