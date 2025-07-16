package repository

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/domain/entity"
	. "github.com/henriquemarlon/cartesi-golang-series/dcm/pkg/custom_type"
)

type CampaignRepository interface {
	CreateCampaign(ctx context.Context, Campaign *entity.Campaign) (*entity.Campaign, error)
	FindCampaignsByDebtor(ctx context.Context, debtor Address) ([]*entity.Campaign, error)
	FindCampaignsByInvestor(ctx context.Context, investor Address) ([]*entity.Campaign, error)
	FindCampaignById(ctx context.Context, id uint) (*entity.Campaign, error)
	FindAllCampaigns(ctx context.Context) ([]*entity.Campaign, error)
	UpdateCampaign(ctx context.Context, Campaign *entity.Campaign) (*entity.Campaign, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	FindOrderById(ctx context.Context, id uint) (*entity.Order, error)
	FindOrdersByCampaignId(ctx context.Context, id uint) ([]*entity.Order, error)
	FindOrdersByState(ctx context.Context, CampaignId uint, state string) ([]*entity.Order, error)
	FindOrdersByInvestor(ctx context.Context, investor Address) ([]*entity.Order, error)
	FindAllOrders(ctx context.Context) ([]*entity.Order, error)
	UpdateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	DeleteOrder(ctx context.Context, id uint) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, User *entity.User) (*entity.User, error)
	FindUsersByRole(ctx context.Context, role string) ([]*entity.User, error)
	FindUserByAddress(ctx context.Context, address Address) (*entity.User, error)
	FindAllUsers(ctx context.Context) ([]*entity.User, error)
	DeleteUser(ctx context.Context, address Address) error
}

type Repository interface {
	CampaignRepository
	OrderRepository
	UserRepository
	Close() error
}
