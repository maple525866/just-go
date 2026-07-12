package gateway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	inventoryv1 "just-go/stage-3-architecture/14-microservices/api/inventory/v1"
	productv1 "just-go/stage-3-architecture/14-microservices/api/product/v1"
)

var (
	ErrInvalidAddress    = errors.New("invalid grpc address")
	ErrConnectionsClosed = errors.New("grpc connections closed")
)

type DialFunc func(context.Context, string) (*grpc.ClientConn, error)

type Connections struct {
	mu     sync.Mutex
	dial   DialFunc
	conns  map[string]*grpc.ClientConn
	closed bool
}

func NewConnections(dial DialFunc) *Connections {
	if dial == nil {
		dial = DefaultDial
	}
	return &Connections{dial: dial, conns: make(map[string]*grpc.ClientConn)}
}

func DefaultDial(ctx context.Context, address string) (*grpc.ClientConn, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (c *Connections) Product(ctx context.Context, address string) (productv1.ProductServiceClient, error) {
	conn, err := c.conn(ctx, address)
	if err != nil {
		return nil, err
	}
	return productv1.NewProductServiceClient(conn), nil
}

func (c *Connections) Inventory(ctx context.Context, address string) (inventoryv1.InventoryServiceClient, error) {
	conn, err := c.conn(ctx, address)
	if err != nil {
		return nil, err
	}
	return inventoryv1.NewInventoryServiceClient(conn), nil
}

func (c *Connections) Close() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	owned := c.conns
	c.conns = make(map[string]*grpc.ClientConn)
	c.mu.Unlock()

	var closeErrors []error
	for address, conn := range owned {
		if err := conn.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("close %s: %w", address, err))
		}
	}
	return errors.Join(closeErrors...)
}

func (c *Connections) conn(ctx context.Context, address string) (*grpc.ClientConn, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: context is required", ErrInvalidAddress)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	address = strings.TrimSpace(address)
	if _, _, err := net.SplitHostPort(address); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidAddress, err)
	}

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil, ErrConnectionsClosed
	}
	if conn := c.conns[address]; conn != nil {
		c.mu.Unlock()
		return conn, nil
	}
	c.mu.Unlock()

	conn, err := c.dial(ctx, address)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		_ = conn.Close()
		return nil, ErrConnectionsClosed
	}
	if existing := c.conns[address]; existing != nil {
		_ = conn.Close()
		return existing, nil
	}
	c.conns[address] = conn
	return conn, nil
}
