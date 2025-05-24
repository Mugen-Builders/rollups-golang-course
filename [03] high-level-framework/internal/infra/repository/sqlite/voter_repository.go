package sqlite

import (
	"github.com/henriquemarlon/cartesi-golang-series/high-level-framework/internal/domain"
)

func (r *SQLiteRepository) CreateVoter(voter *domain.Voter) error {
	return r.db.Create(voter).Error
}

func (r *SQLiteRepository) GetVoterByID(id int) (*domain.Voter, error) {
	var voter domain.Voter
	err := r.db.First(&voter, id).Error
	if err != nil {
		return nil, err
	}
	return &voter, nil
}

func (r *SQLiteRepository) GetVoterByAddress(address string) (*domain.Voter, error) {
	var voter domain.Voter
	err := r.db.Where("address = ?", address).First(&voter).Error
	if err != nil {
		return nil, err
	}
	return &voter, nil
}

func (r *SQLiteRepository) UpdateVoter(voter *domain.Voter) error {
	return r.db.Save(voter).Error
}

func (r *SQLiteRepository) DeleteVoter(id int) error {
	return r.db.Delete(&domain.Voter{}, id).Error
}

func (r *SQLiteRepository) HasVoted(voterID, votingID int) (bool, error) {
	var count int64
	err := r.db.Model(&domain.VotingOption{}).
		Joins("JOIN votes ON votes.option_id = voting_options.id").
		Where("votes.voter_id = ? AND voting_options.voting_id = ?", voterID, votingID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
