package inspect

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/infra/repository"
	auction_usecase "github.com/henriquemarlon/cartesi-golang-series/auction/internal/usecase/auction"
	"github.com/rollmelette/rollmelette"
)

type AuctionInspectHandlers struct {
	AuctionRepository repository.AuctionRepository
}

func NewAuctionInspectHandlers(auctionRepository repository.AuctionRepository) *AuctionInspectHandlers {
	return &AuctionInspectHandlers{
		AuctionRepository: auctionRepository,
	}
}

func (h *AuctionInspectHandlers) FindAuctionById(env rollmelette.EnvInspector, payload []byte) error {
	var input auction_usecase.FindAuctionByIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findAuctionById := auction_usecase.NewFindAuctionByIdUseCase(h.AuctionRepository)
	res, err := findAuctionById.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find auction: %w", err)
	}
	auction, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal auction: %w", err)
	}
	env.Report(auction)
	return nil
}

func (h *AuctionInspectHandlers) FindAllAuctions(env rollmelette.EnvInspector, payload []byte) error {
	ctx := context.Background()
	findAllAuctionsUseCase := auction_usecase.NewFindAllAuctionsUseCase(h.AuctionRepository)
	res, err := findAllAuctionsUseCase.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all auctions: %w", err)
	}
	allAuctions, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all auctions: %w", err)
	}
	env.Report(allAuctions)
	return nil
}

func (h *AuctionInspectHandlers) FindAuctionsByInvestor(env rollmelette.EnvInspector, payload []byte) error {
	var input auction_usecase.FindAuctionsByInvestorInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findAuctionsByInvestor := auction_usecase.NewFindAuctionsByInvestorUseCase(h.AuctionRepository)
	res, err := findAuctionsByInvestor.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find auctions by investor: %w", err)
	}
	auctions, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal auctions: %w", err)
	}
	env.Report(auctions)
	return nil
}

func (h *AuctionInspectHandlers) FindAuctionsByCreator(env rollmelette.EnvInspector, payload []byte) error {
	var input auction_usecase.FindAuctionsByCreatorInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findAuctionsByCreator := auction_usecase.NewFindAuctionsByCreatorUseCase(h.AuctionRepository)
	res, err := findAuctionsByCreator.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find auctions by creator: %w", err)
	}
	auctions, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal auctions: %w", err)
	}
	env.Report(auctions)
	return nil
}
