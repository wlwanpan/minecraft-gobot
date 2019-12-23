package bot

import (
	"context"
	"errors"
	"log"

	"github.com/wlwanpan/minecraft-gobot/config"
	"github.com/wlwanpan/minecraft-gobot/services"
	"google.golang.org/grpc"
)

var (
	ErrMcsNotResponding = errors.New("mcs not responding")

	ErrPingAttemptFailed = errors.New("mcs ping attemp failed")
)

type mcsClient struct {
	grpcConn *grpc.ClientConn
	client   services.McsServiceClient
}

func (c *mcsClient) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &services.PingReq{})
	return err
}

func (c *mcsClient) Start(ctx context.Context, in *services.StartConfig) (*services.ServiceResp, error) {
	return c.client.Start(ctx, in)
}

func (c *mcsClient) Stop(ctx context.Context) (*services.ServiceResp, error) {
	return c.client.Stop(ctx, &services.EmptyReq{})
}

func (c *mcsClient) Status(ctx context.Context) (*services.StatusResp, error) {
	return c.client.Status(ctx, &services.EmptyReq{})
}

func (c *mcsClient) checkConn(ctx context.Context) error {
	if c.grpcConn != nil && c.client != nil {
		return c.Ping(ctx)
	}
	if err := c.initConn(); err != nil {
		return err
	}
	return c.Ping(ctx)
}

func (c *mcsClient) initConn() error {
	conn, err := grpc.Dial(config.Cfg.Bot.McsAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("error dialing mcs: %s", err)
		return ErrMcsNotResponding
	}
	c.grpcConn = conn
	c.client = services.NewMcsServiceClient(conn)

	log.Printf("successfully connected to mcs: %s", config.Cfg.Bot.McsAddr)
	return nil
}

func (c *mcsClient) closeConn() error {
	if c.grpcConn != nil {
		if err := c.grpcConn.Close(); err != nil {
			return err
		}
	}

	c.grpcConn = nil
	c.client = nil
	return nil
}
