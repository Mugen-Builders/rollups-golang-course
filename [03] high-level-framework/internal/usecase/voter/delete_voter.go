package voter

import (
	"context"
	"errors"

	. "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/pkg/custom_type"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
)

type DeleteVoterInputDTO struct {
	Id int `json:"id"`
}

type DeleteVoterOutputDTO struct {
	Success bool `json:"success"`
}

type DeleteVoterUseCase struct {
	VoterRepository repository.VoterRepository
}

func NewDeleteVoterUseCase(voterRepository repository.VoterRepository) *DeleteVoterUseCase {
	return &DeleteVoterUseCase{VoterRepository: voterRepository}
}

func (uc *DeleteVoterUseCase) Execute(ctx context.Context, input *DeleteVoterInputDTO) (*DeleteVoterOutputDTO, error) {
	voter, err := uc.VoterRepository.FindVoterByID(input.Id)
	if err != nil {
		return &DeleteVoterOutputDTO{Success: false}, err
	}
	msgSender, ok := ctx.Value("msg_sender").(string)
	if !ok {
		return &DeleteVoterOutputDTO{Success: false}, errors.New("error getting msg_sender")
	}
	if voter.Address != HexToAddress(msgSender) {
		return &DeleteVoterOutputDTO{Success: false}, errors.New("unauthorized")
	}
	err = uc.VoterRepository.DeleteVoter(input.Id)
	if err != nil {
		return &DeleteVoterOutputDTO{Success: false}, err
	}
	return &DeleteVoterOutputDTO{Success: true}, nil
}
