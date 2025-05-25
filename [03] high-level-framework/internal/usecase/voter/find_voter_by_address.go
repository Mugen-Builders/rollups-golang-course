package voter

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type FindVoterByAddressInputDTO struct {
	Address string `json:"address"`
}

type FindVoterByAddressOutputDTO struct {
	Id      int    `json:"id"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

type FindVoterByAddressUseCase struct {
	VoterRepository repository.VoterRepository
}

func NewFindVoterByAddressUseCase(voterRepository repository.VoterRepository) *FindVoterByAddressUseCase {
	return &FindVoterByAddressUseCase{VoterRepository: voterRepository}
}

func (uc *FindVoterByAddressUseCase) Execute(ctx context.Context, input *FindVoterByAddressInputDTO) (*FindVoterByAddressOutputDTO, error) {
	voter, err := uc.VoterRepository.FindVoterByAddress(input.Address)
	if err != nil {
		return nil, err
	}
	return &FindVoterByAddressOutputDTO{
		Id:      voter.ID,
		Address: voter.Address,
		Name:    voter.Name,
	}, nil
}
