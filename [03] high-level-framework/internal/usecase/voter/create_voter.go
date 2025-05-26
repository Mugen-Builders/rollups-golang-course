package voter

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/pkg/custom_type"
)

type CreateVoterInputDTO struct {
	Address Address `json:"address"`
}

type CreateVoterOutputDTO struct {
	Id      int     `json:"id"`
	Address Address `json:"address"`
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
	}
	err := uc.VoterRepository.CreateVoter(voter)
	if err != nil {
		return nil, err
	}
	return &CreateVoterOutputDTO{
		Id:      voter.ID,
		Address: voter.Address,
	}, nil
}
