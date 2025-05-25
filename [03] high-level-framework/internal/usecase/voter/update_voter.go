package voter

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type UpdateVoterInputDTO struct {
	Id      int    `json:"id"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

type UpdateVoterOutputDTO struct {
	Id      int    `json:"id"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

type UpdateVoterUseCase struct {
	VoterRepository repository.VoterRepository
}

func NewUpdateVoterUseCase(voterRepository repository.VoterRepository) *UpdateVoterUseCase {
	return &UpdateVoterUseCase{VoterRepository: voterRepository}
}

func (uc *UpdateVoterUseCase) Execute(ctx context.Context, input *UpdateVoterInputDTO) (*UpdateVoterOutputDTO, error) {
	voter := &domain.Voter{
		ID:      input.Id,
		Address: input.Address,
		Name:    input.Name,
	}
	err := uc.VoterRepository.UpdateVoter(voter)
	if err != nil {
		return nil, err
	}
	return &UpdateVoterOutputDTO{
		Id:      voter.ID,
		Address: voter.Address,
		Name:    voter.Name,
	}, nil
}
