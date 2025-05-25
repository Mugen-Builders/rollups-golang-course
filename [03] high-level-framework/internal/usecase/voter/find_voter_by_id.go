package voter

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type FindVoterByIDInputDTO struct {
	Id int `json:"id"`
}

type FindVoterByIDOutputDTO struct {
	Id      int    `json:"id"`
	Address string `json:"address"`
	Name    string `json:"name"`
}

type FindVoterByIDUseCase struct {
	VoterRepository repository.VoterRepository
}

func NewFindVoterByIDUseCase(voterRepository repository.VoterRepository) *FindVoterByIDUseCase {
	return &FindVoterByIDUseCase{VoterRepository: voterRepository}
}

func (uc *FindVoterByIDUseCase) Execute(ctx context.Context, input *FindVoterByIDInputDTO) (*FindVoterByIDOutputDTO, error) {
	voter, err := uc.VoterRepository.FindVoterByID(input.Id)
	if err != nil {
		return nil, err
	}
	return &FindVoterByIDOutputDTO{
		Id:      voter.ID,
		Address: voter.Address,
		Name:    voter.Name,
	}, nil
}
