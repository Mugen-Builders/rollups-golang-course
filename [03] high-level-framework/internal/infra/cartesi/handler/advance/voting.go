package advance_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/usecase/voting"
	"github.com/rollmelette/rollmelette"
)

type VotingAdvanceHandlers struct {
	VotingRepository repository.VotingRepository
}

func NewVotingAdvanceHandlers(votingRepository repository.VotingRepository) *VotingAdvanceHandlers {
	return &VotingAdvanceHandlers{
		VotingRepository: votingRepository,
	}
}

func (h *VotingAdvanceHandlers) CreateVoting(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input voting.CreateVotingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "msg_sender", metadata.MsgSender)
	createVoting := voting.NewCreateVotingUseCase(h.VotingRepository)
	res, err := createVoting.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to create voting: %w", err)
	}
	votingBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	env.Notice(append([]byte("voting created - "), votingBytes...))
	return nil
}

func (h *VotingAdvanceHandlers) DeleteVoting(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input voting.DeleteVotingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}
	ctx := context.Background()
	deleteVoting := voting.NewDeleteVotingUseCase(h.VotingRepository)
	res, err := deleteVoting.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to delete voting: %w", err)
	}
	votingBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	env.Notice(append([]byte("voting deleted - "), votingBytes...))
	return nil
}
