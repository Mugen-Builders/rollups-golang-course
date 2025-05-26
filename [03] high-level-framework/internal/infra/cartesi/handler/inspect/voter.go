package inspect_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/usecase/voter"
	"github.com/rollmelette/rollmelette"
)

type VoterInspectHandlers struct {
	VoterRepository repository.VoterRepository
}

func NewVoterInspectHandlers(voterRepository repository.VoterRepository) *VoterInspectHandlers {
	return &VoterInspectHandlers{
		VoterRepository: voterRepository,
	}
}

func (h *VoterInspectHandlers) FindVoterByID(env rollmelette.EnvInspector, payload []byte) error {
	var input voter.FindVoterByIDInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}
	ctx := context.Background()
	findVoterByID := voter.NewFindVoterByIDUseCase(h.VoterRepository)
	voterRes, err := findVoterByID.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find voter by id: %w", err)
	}
	voterBytes, err := json.Marshal(voterRes)
	if err != nil {
		return fmt.Errorf("failed to marshal voter: %w", err)
	}
	env.Report(voterBytes)
	return nil
}

func (h *VoterInspectHandlers) FindVoterByAddress(env rollmelette.EnvInspector, payload []byte) error {
	var input voter.FindVoterByAddressInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}
	ctx := context.Background()
	findVoterByAddress := voter.NewFindVoterByAddressUseCase(h.VoterRepository)
	voterRes, err := findVoterByAddress.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find voter by address: %w", err)
	}
	voterBytes, err := json.Marshal(voterRes)
	if err != nil {
		return fmt.Errorf("failed to marshal voter: %w", err)
	}
	env.Report(voterBytes)
	return nil
}
