package repository

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/auction/internal/domain"
	. "github.com/henriquemarlon/cartesi-golang-series/auction/pkg/custom_type"
)

type AuctionRepository interface {
	CreateAuction(ctx context.Context, Auction *domain.Auction) (*domain.Auction, error)
	FindAuctionsByCreator(ctx context.Context, creator Address) ([]*domain.Auction, error)
	FindAuctionsByInvestor(ctx context.Context, investor Address) ([]*domain.Auction, error)
	FindAuctionById(ctx context.Context, id uint) (*domain.Auction, error)
	FindAllAuctions(ctx context.Context) ([]*domain.Auction, error)
	UpdateAuction(ctx context.Context, Auction *domain.Auction) (*domain.Auction, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	FindOrderById(ctx context.Context, id uint) (*domain.Order, error)
	FindOrdersByAuctionId(ctx context.Context, id uint) ([]*domain.Order, error)
	FindOrdersByState(ctx context.Context, AuctionId uint, state string) ([]*domain.Order, error)
	FindOrdersByInvestor(ctx context.Context, investor Address) ([]*domain.Order, error)
	FindAllOrders(ctx context.Context) ([]*domain.Order, error)
	UpdateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	DeleteOrder(ctx context.Context, id uint) error
}

type Repository interface {
	AuctionRepository
	OrderRepository
	Close() error
}
