package sqlite

import (
	"context"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateCampaign(ctx context.Context, input *entity.Campaign) (*entity.Campaign, error) {
	if err := r.Db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}
	return input, nil
}

func (r *SQLiteRepository) FindCampaignById(ctx context.Context, id uint) (*entity.Campaign, error) {
	var Campaign entity.Campaign
	if err := r.Db.WithContext(ctx).
		Preload("Orders").
		First(&Campaign, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrCampaignNotFound
		}
		return nil, fmt.Errorf("failed to find campaign by id: %w", err)
	}
	return &Campaign, nil
}

func (r *SQLiteRepository) FindAllCampaigns(ctx context.Context) ([]*entity.Campaign, error) {
	var Campaigns []*entity.Campaign
	if err := r.Db.WithContext(ctx).
		Preload("Orders").
		Find(&Campaigns).Error; err != nil {
		return nil, fmt.Errorf("failed to find all campaigns: %w", err)
	}
	return Campaigns, nil
}

func (r *SQLiteRepository) FindCampaignsByInvestor(ctx context.Context, investor Address) ([]*entity.Campaign, error) {
	var Campaigns []*entity.Campaign
	if err := r.Db.WithContext(ctx).
		Joins("JOIN orders ON orders.campaign_id = campaigns.id").
		Where("orders.investor = ?", investor).
		Preload("Orders").
		Find(&Campaigns).Error; err != nil {
		return nil, fmt.Errorf("failed to find Campaigns by investor: %w", err)
	}
	return Campaigns, nil
}

func (r *SQLiteRepository) FindCampaignsByDebtor(ctx context.Context, debtor Address) ([]*entity.Campaign, error) {
	var Campaigns []*entity.Campaign
	if err := r.Db.WithContext(ctx).
		Where("debtor = ?", debtor).
		Preload("Orders").
		Find(&Campaigns).Error; err != nil {
		return nil, fmt.Errorf("failed to find campaigns by debtor: %w", err)
	}
	return Campaigns, nil
}

func (r *SQLiteRepository) UpdateCampaign(ctx context.Context, input *entity.Campaign) (*entity.Campaign, error) {
	if err := r.Db.WithContext(ctx).Updates(&input).Error; err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}
	Campaign, err := r.FindCampaignById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	return Campaign, nil
}
