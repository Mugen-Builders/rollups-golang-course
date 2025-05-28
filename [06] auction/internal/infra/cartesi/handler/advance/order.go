package advance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/usecase/order"
	"github.com/rollmelette/rollmelette"
)

type OrderAdvanceHandlers struct {
	OrderRepository   repository.OrderRepository
	AuctionRepository repository.AuctionRepository
}

func NewOrderAdvanceHandlers(
	orderRepository repository.OrderRepository,
	auctionRepository repository.AuctionRepository,
) *OrderAdvanceHandlers {
	return &OrderAdvanceHandlers{
		OrderRepository:   orderRepository,
		AuctionRepository: auctionRepository,
	}
}

func (h *OrderAdvanceHandlers) CreateOrder(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input order.CreateOrderInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w, payload: %s", err, string(payload))
	}

	ctx := context.Background()
	createOrder := order.NewCreateOrderUseCase(
		h.OrderRepository,
		h.AuctionRepository,
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
		erc20Deposit.Amount,
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

	ctx := context.Background()
	cancelOrder := order.NewCancelOrderUseCase(
		h.OrderRepository,
		h.AuctionRepository,
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
