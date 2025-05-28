package sqlite

import (
	"context"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateAuction(ctx context.Context, input *domain.Auction) (*domain.Auction, error) {
	if err := r.Db.WithContext(ctx).Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create Auction: %w", err)
	}
	return input, nil
}

func (r *SQLiteRepository) FindAuctionById(ctx context.Context, id uint) (*domain.Auction, error) {
	var Auction domain.Auction
	if err := r.Db.WithContext(ctx).
		Preload("Orders").
		First(&Auction, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrAuctionNotFound
		}
		return nil, fmt.Errorf("failed to find Auction by id: %w", err)
	}
	return &Auction, nil
}

func (r *SQLiteRepository) FindAllAuctions(ctx context.Context) ([]*domain.Auction, error) {
	var Auctions []*domain.Auction
	if err := r.Db.WithContext(ctx).
		Preload("Orders").
		Find(&Auctions).Error; err != nil {
		return nil, fmt.Errorf("failed to find all Auctions: %w", err)
	}
	return Auctions, nil
}

func (r *SQLiteRepository) FindAuctionsByInvestor(ctx context.Context, investor Address) ([]*domain.Auction, error) {
	var Auctions []*domain.Auction
	if err := r.Db.WithContext(ctx).
		Joins("JOIN orders ON orders.Auction_id = Auctions.id").
		Where("orders.investor = ?", investor).
		Preload("Orders").
		Find(&Auctions).Error; err != nil {
		return nil, fmt.Errorf("failed to find Auctions by investor: %w", err)
	}
	return Auctions, nil
}

func (r *SQLiteRepository) FindAuctionsByCreator(ctx context.Context, creator Address) ([]*domain.Auction, error) {
	var Auctions []*domain.Auction
	if err := r.Db.WithContext(ctx).
		Where("creator = ?", creator).
		Preload("Orders").
		Find(&Auctions).Error; err != nil {
		return nil, fmt.Errorf("failed to find Auctions by creator: %w", err)
	}
	return Auctions, nil
}

func (r *SQLiteRepository) UpdateAuction(ctx context.Context, input *domain.Auction) (*domain.Auction, error) {
	if err := r.Db.WithContext(ctx).Updates(&input).Error; err != nil {
		return nil, fmt.Errorf("failed to update Auction: %w", err)
	}
	Auction, err := r.FindAuctionById(ctx, input.Id)
	if err != nil {
		return nil, err
	}
	return Auction, nil
}