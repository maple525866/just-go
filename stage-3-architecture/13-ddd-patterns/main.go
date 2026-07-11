package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"just-go/stage-3-architecture/13-ddd-patterns/application"
	"just-go/stage-3-architecture/13-ddd-patterns/domain/order"
	"just-go/stage-3-architecture/13-ddd-patterns/infrastructure/events"
	"just-go/stage-3-architecture/13-ddd-patterns/infrastructure/memory"
)

func main() {
	if err := run(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(output io.Writer) error {
	repository := memory.NewOrderRepository()
	bus := events.NewBus()
	inventory := events.NewInventoryProjection()
	if err := bus.Register(order.OrderConfirmedName, inventory.Handle); err != nil {
		return err
	}
	service, err := application.NewService(repository, bus, order.NoDiscount{}, time.Now)
	if err != nil {
		return err
	}
	address, err := order.NewAddress("Alice", "1 Go Road", "Beijing", "100000", "CN")
	if err != nil {
		return err
	}
	price, err := order.NewMoney(5000, "CNY")
	if err != nil {
		return err
	}

	ctx := context.Background()
	if _, err = service.Create(ctx, application.CreateOrder{
		OrderID: "order-2026", CustomerID: "customer-1", Address: address,
	}); err != nil {
		return err
	}
	if _, err = service.AddLine(ctx, application.AddOrderLine{
		OrderID: "order-2026", LineID: "line-1", ProductID: "go-book",
		Name: "Domain-Driven Go", UnitPrice: price, Quantity: 2,
	}); err != nil {
		return err
	}
	confirmed, err := service.Confirm(ctx, "order-2026")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(output, "order=%s status=%s total=%d %s inventory_reserved=%d\n",
		confirmed.ID(), confirmed.Status(), confirmed.Total().Minor(), confirmed.Total().Currency(), inventory.Reserved("go-book"))
	return err
}
