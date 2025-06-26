package advance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/usecase/campaign"
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

type CampaignAdvanceHandlers struct {
	OrderRepository    repository.OrderRepository
	UserRepository     repository.UserRepository
	CampaignRepository repository.CampaignRepository
}

func NewCampaignAdvanceHandlers(
	orderRepository repository.OrderRepository,
	userRepository repository.UserRepository,
	campaignRepository repository.CampaignRepository,
) *CampaignAdvanceHandlers {
	return &CampaignAdvanceHandlers{
		OrderRepository:    orderRepository,
		UserRepository:     userRepository,
		CampaignRepository: campaignRepository,
	}
}

func (h *CampaignAdvanceHandlers) CreateCampaign(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input campaign.CreateCampaignInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	createCampaign := campaign.NewCreateCampaignUseCase(
		h.CampaignRepository,
		h.UserRepository,
	)

	res, err := createCampaign.Execute(ctx, &input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to create campaign: %w", err)
	}

	erc20Deposit := deposit.(*rollmelette.ERC20Deposit)
	if err := env.ERC20Transfer(
		erc20Deposit.Token,
		erc20Deposit.Sender,
		env.AppAddress(),
		erc20Deposit.Value,
	); err != nil {
		return fmt.Errorf("failed to transfer ERC20: %w", err)
	}

	campaign, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("campaign created - "), campaign...))
	return nil
}

func (h *CampaignAdvanceHandlers) CloseCampaign(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input campaign.CloseCampaignInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	closeCampaign := campaign.NewCloseCampaignUseCase(h.CampaignRepository, h.OrderRepository)
	res, err := closeCampaign.Execute(ctx, &input, metadata)
	if err != nil && res == nil {
		return fmt.Errorf("failed to close campaign: %w", err)
	}

	token := common.Address(res.Token)

	// Process orders
	for _, order := range res.Orders {
		if order.State == entity.OrderStateRejected {
			if err = env.ERC20Transfer(
				token,
				env.AppAddress(),
				common.Address(order.Investor),
				order.Amount.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer rejected order: %w", err)
			}
		}
	}

	if err := env.ERC20Transfer(token, env.AppAddress(), common.Address(res.Issuer), res.TotalRaised.ToBig()); err != nil {
		return fmt.Errorf("failed to transfer total raised: %w", err)
	}

	campaign, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte(fmt.Sprintf("campaign %v - ", res.State)), campaign...))
	return nil
}

func (h *CampaignAdvanceHandlers) SettleCampaign(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input campaign.SettleCampaignInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	settleCampaign := campaign.NewSettleCampaignUseCase(
		h.CampaignRepository,
		h.OrderRepository,
	)

	res, err := settleCampaign.Execute(ctx, &input, deposit, metadata)
	if err != nil {
		return fmt.Errorf("failed to settle campaign: %w", err)
	}

	contractAddr := common.Address(res.Token)
	issuerAddr := common.Address(res.Issuer)

	// Process settled orders
	for _, order := range res.Orders {
		if order.State == entity.OrderStateSettled {
			// Calculate interest for this order
			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100))

			// Calculate total payment
			totalPayment := new(uint256.Int).Add(order.Amount, interest)

			if err := env.ERC20Transfer(
				contractAddr,
				issuerAddr,
				common.Address(order.Investor),
				totalPayment.ToBig(),
			); err != nil {
				return fmt.Errorf("failed to transfer settled order: %w", err)
			}
		}
	}

	campaign, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("campaign settled - "), campaign...))
	return nil
}

func (h *CampaignAdvanceHandlers) ExecuteCampaignCollateral(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input campaign.ExecuteCampaignCollateralInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	executeCampaignCollateral := campaign.NewExecuteCampaignCollateralUseCase(h.CampaignRepository, h.OrderRepository)
	res, err := executeCampaignCollateral.Execute(ctx, &input, metadata)
	if err != nil {
		return fmt.Errorf("failed to execute campaign collateral: %w", err)
	}

	totalFinalValue := uint256.NewInt(0)
	orderFinalValues := make(map[uint]*uint256.Int)
	for _, order := range res.Orders {
		if order.State == entity.OrderStateSettledByCollateral {
			interest := new(uint256.Int).Mul(order.Amount, order.InterestRate)
			interest.Div(interest, uint256.NewInt(100))
			finalValue := new(uint256.Int).Add(order.Amount, interest)
			orderFinalValues[order.Id] = finalValue
			totalFinalValue.Add(totalFinalValue, finalValue)
		}
	}

	for _, order := range res.Orders {
		if order.State == entity.OrderStateSettledByCollateral {
			finalValue := orderFinalValues[order.Id]
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

	campaign, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	env.Notice(append([]byte("campaign collateral executed - "), campaign...))
	return nil
}
