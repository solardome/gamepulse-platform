package account

import (
	"context"
	"fmt"
	"time"

	accountv1 "github.com/solardome/gamepulse-platform/gen/account/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	rpc     accountv1.AccountServiceClient
	timeout time.Duration
}

func New(target string, timeout time.Duration) (*Client, error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("create grpc client for target %q: %w", target, err)
	}

	return &Client{
		conn:    conn,
		rpc:     accountv1.NewAccountServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Ping(ctx context.Context, origin string) (*accountv1.PingResponse, error) {
	callCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.rpc.Ping(callCtx, &accountv1.PingRequest{
		Origin: origin,
	})
}
