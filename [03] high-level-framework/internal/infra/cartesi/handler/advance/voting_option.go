package advance_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/usecase/voting_option"
	"github.com/rollmelette/rollmelette"
)

type VotingOptionAdvanceHandlers struct {
	VotingOptionRepository repository.VotingOptionRepository
}

func NewVotingOptionAdvanceHandlers(votingOptionRepository repository.VotingOptionRepository) *VotingOptionAdvanceHandlers {
	return &VotingOptionAdvanceHandlers{
		VotingOptionRepository: votingOptionRepository,
	}
}

func (h *VotingOptionAdvanceHandlers) CreateVotingOption(env rollmelette.Env, payload []byte) error {
	var input voting_option.CreateVotingOptionInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}
	ctx := context.Background()
	createVotingOption := voting_option.NewCreateVotingOptionUseCase(h.VotingOptionRepository)
	res, err := createVotingOption.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to create voting option: %w", err)
	}
	votingOptionBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	env.Notice(append([]byte("voting option created - "), votingOptionBytes...))
	return nil
}

func (h *VotingOptionAdvanceHandlers) DeleteVotingOption(env rollmelette.Env, payload []byte) error {
	var input voting_option.DeleteVotingOptionInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}
	ctx := context.Background()
	deleteVotingOption := voting_option.NewDeleteVotingOptionUseCase(h.VotingOptionRepository)
	res, err := deleteVotingOption.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to delete voting option: %w", err)
	}
	votingOptionBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	env.Notice(append([]byte("voting option deleted - "), votingOptionBytes...))
	return nil
}
