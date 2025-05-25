package voter

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type CreateVoterInputDTO struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type CreateVoterOutputDTO struct {
	Id      int    `json:"id"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

type CreateVoterUseCase struct {
	VoterRepository repository.VoterRepository
}

func NewCreateVoterUseCase(voterRepository repository.VoterRepository) *CreateVoterUseCase {
	return &CreateVoterUseCase{VoterRepository: voterRepository}
}

func (uc *CreateVoterUseCase) Execute(ctx context.Context, input *CreateVoterInputDTO) (*CreateVoterOutputDTO, error) {
	voter := &domain.Voter{
		Address: input.Address,
		Name:    input.Name,
	}
	err := uc.VoterRepository.CreateVoter(voter)
	if err != nil {
		return nil, err
	}
	return &CreateVoterOutputDTO{
		Id:      voter.ID,
		Address: voter.Address,
		Name:    voter.Name,
	}, nil
}
