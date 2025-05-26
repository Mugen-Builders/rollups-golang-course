package voting

import (
	"context"
	"errors"

	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/high-level-framework/pkg/custom_type"
)

type DeleteVotingInputDTO struct {
	Id int `json:"id"`
}

type DeleteVotingOutputDTO struct {
	Success bool `json:"success"`
}

type DeleteVotingUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewDeleteVotingUseCase(votingRepository repository.VotingRepository) *DeleteVotingUseCase {
	return &DeleteVotingUseCase{VotingRepository: votingRepository}
}

func (uc *DeleteVotingUseCase) Execute(ctx context.Context, input *DeleteVotingInputDTO) (*DeleteVotingOutputDTO, error) {
	voting, err := uc.VotingRepository.FindVotingByID(input.Id)
	if err != nil {
		return &DeleteVotingOutputDTO{Success: false}, err
	}
	msgSender, ok := ctx.Value("msg_sender").(string)
	if !ok {
		return &DeleteVotingOutputDTO{Success: false}, errors.New("error getting msg_sender")
	}
	if voting.Creator != HexToAddress(msgSender) {
		return &DeleteVotingOutputDTO{Success: false}, errors.New("unauthorized")
	}
	err = uc.VotingRepository.DeleteVoting(input.Id)
	if err != nil {
		return &DeleteVotingOutputDTO{Success: false}, err
	}
	return &DeleteVotingOutputDTO{Success: true}, nil
}
