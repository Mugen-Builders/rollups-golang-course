package advance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/usecase/auction"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type AuctionAdvanceHandlers struct {
	OrderRepository   repository.OrderRepository
	AuctionRepository repository.AuctionRepository
}

func NewAuctionAdvanceHandlers(
	orderRepository repository.OrderRepository,
	auctionRepository repository.AuctionRepository,
) *AuctionAdvanceHandlers {
	return &AuctionAdvanceHandlers{
		OrderRepository:   orderRepository,
		AuctionRepository: auctionRepository,
	}
}

func (h *AuctionAdvanceHandlers) CreateAuction(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input auction.CreateAuctionInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	createAuction := auction.NewCreateAuctionUseCase(
		h.AuctionRepository,
	)

	res, err := createAuction.Execute(ctx, &input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to create auction: %w", err)
	}

	erc20Deposit := deposit.(*rollmelette.ERC20Deposit)
	if err := env.ERC20Transfer(
		erc20Deposit.Token,
		erc20Deposit.Sender,
		env.AppAddress(),
		erc20Deposit.Amount,
	); err != nil {
		return fmt.Errorf("failed to transfer ERC20: %w", err)
	}

	auction, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("auction created - "), auction...))
	return nil
}

func (h *AuctionAdvanceHandlers) CloseAuction(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input auction.CloseAuctionInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	closeAuction := auction.NewCloseAuctionUseCase(h.AuctionRepository, h.OrderRepository)
	res, err := closeAuction.Execute(ctx, &input, metadata)
	if err != nil && res == nil {
		return fmt.Errorf("failed to close auction: %w", err)
	}

	stablecoinAddr := common.Address(res.Token)

	// Process orders
	for _, order := range res.Orders {
		if order.State == domain.OrderStateRejected {
			if err = env.ERC20Transfer(
				stablecoinAddr,
				env.AppAddress(),
				common.Address(order.Investor),
				order.Amount.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer rejected order: %w", err)
			}
		}
	}

	// Transfer amount raised to creator
	if err = env.ERC20Transfer(
		stablecoinAddr,
		env.AppAddress(),
		common.Address(res.Creator),
		res.DebtIssued.ToBig(),
	); err != nil {
		return fmt.Errorf("failed to transfer debt issued: %w", err)
	}

	auction, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte(fmt.Sprintf("auction %v - ", res.State)), auction...))
	return nil
}

func (h *AuctionAdvanceHandlers) SettleAuction(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input auction.SettleAuctionInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	settleAuction := auction.NewSettleAuctionUseCase(
		h.AuctionRepository,
		h.OrderRepository,
	)

	res, err := settleAuction.Execute(ctx, &input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to settle auction: %w", err)
	}

	// Reuse variables for calculations
	interest := new(uint256.Int)
	contractAddr := common.Address(res.Token)
	creatorAddr := common.Address(res.Creator)

	// Process settled orders
	for _, order := range res.Orders {
		if order.State == domain.OrderStateSettled {
			// Calculate interest
			interest.Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100))

			// Calculate total payment
			totalPayment := new(uint256.Int).Add(order.Amount, interest)

			if err := env.ERC20Transfer(
				contractAddr,
				creatorAddr,
				common.Address(order.Investor),
				totalPayment.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer settled order: %w", err)
			}
		}
	}

	auction, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("auction settled - "), auction...))
	return nil
}

func (h *AuctionAdvanceHandlers) ExecuteAuctionCollateral(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input auction.ExecuteAuctionCollateralInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	executeAuctionCollateral := auction.NewExecuteAuctionCollateralUseCase(h.AuctionRepository, h.OrderRepository)
	res, err := executeAuctionCollateral.Execute(ctx, &input, metadata)
	if err != nil {
		return fmt.Errorf("failed to execute auction collateral: %w", err)
	}

	totalFinalValue := uint256.NewInt(0)
	orderFinalValues := make(map[*domain.Order]*uint256.Int)
	for _, order := range res.Orders {
		if order.State == domain.OrderStateAccepted || order.State == domain.OrderStatePartiallyAccepted {
			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100))
			finalValue := new(uint256.Int).Add(order.Amount, interest)
			orderFinalValues[order] = finalValue
			totalFinalValue.Add(totalFinalValue, finalValue)
		}
	}

	for _, order := range res.Orders {
		if order.State == domain.OrderStateAccepted || order.State == domain.OrderStatePartiallyAccepted {
			finalValue := orderFinalValues[order]
			orderShare := new(uint256.Int).Mul(finalValue, res.CollateralAmount)
			orderShare.Div(orderShare, totalFinalValue)

			if err = env.ERC20Transfer(
				common.Address(res.CollateralAddress),
				env.AppAddress(),
				common.Address(order.Investor),
				orderShare.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer collateral to investor: %w", err)
			}
		}
	}

	auction, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("auction collateral executed - "), auction...))
	return nil
}
