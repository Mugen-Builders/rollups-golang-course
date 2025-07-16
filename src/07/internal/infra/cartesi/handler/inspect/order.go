package inspect

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/dcm/internal/infra/repository"
	order "github.com/henriquemarlon/cartesi-golang-series/dcm/internal/usecase/order"
	"github.com/rollmelette/rollmelette"
)

type OrderInspectHandlers struct {
	OrderRepository repository.OrderRepository
}

func NewOrderInspectHandlers(orderRepository repository.OrderRepository) *OrderInspectHandlers {
	return &OrderInspectHandlers{
		OrderRepository: orderRepository,
	}
}

func (h *OrderInspectHandlers) FindOrderById(env rollmelette.EnvInspector, payload []byte) error {
	var input order.FindOrderByIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findOrderById := order.NewFindOrderByIdUseCase(h.OrderRepository)
	res, err := findOrderById.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}
	order, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}
	env.Report(order)
	return nil
}

func (h *OrderInspectHandlers) FindBidsByCampaignId(env rollmelette.EnvInspector, payload []byte) error {
	var input order.FindOrdersByCampaignIdInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findOrdersByCampaignId := order.NewFindOrdersByCampaignIdUseCase(h.OrderRepository)
	res, err := findOrdersByCampaignId.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find orders by campaign id: %v", err)
	}
	orders, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal orders: %w", err)
	}
	env.Report(orders)
	return nil
}

func (h *OrderInspectHandlers) FindAllOrders(env rollmelette.EnvInspector, payload []byte) error {
	ctx := context.Background()
	findAllOrders := order.NewFindAllOrdersUseCase(h.OrderRepository)
	res, err := findAllOrders.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all orders: %w", err)
	}
	allOrders, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal all orders: %w", err)
	}
	env.Report(allOrders)
	return nil
}

func (h *OrderInspectHandlers) FindOrdersByInvestor(env rollmelette.EnvInspector, payload []byte) error {
	var input order.FindOrdersByInvestorInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	ctx := context.Background()
	findOrdersByInvestor := order.NewFindOrdersByInvestorUseCase(h.OrderRepository)
	res, err := findOrdersByInvestor.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find orders by investor: %w", err)
	}
	orders, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal orders: %w", err)
	}
	env.Report(orders)
	return nil
}
