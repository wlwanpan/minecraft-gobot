package launcher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/wlwanpan/minecraft-manager/messages"
	"google.golang.org/grpc"
)

const (
	DEFAULT_GRPC_TRANSPORT_PROTOCOL = "tcp"

	DEFAULT_TCP_PORT = 7777
)

var (
	minecraftLauncher *MinecraftLauncher

	ErrLauncherOffline = errors.New("launcher is offline")
)

type CmdService struct{}

func (cmds *CmdService) Status(ctx context.Context, _ *messages.EmptyReq) (*messages.ServiceResp, error) {
	var message string
	if minecraftLauncher == nil {
		message = "OFFLINE"
	} else {
		switch minecraftLauncher.currState {
		case LAUNCHER_STATE_INIT:
			message = "STOPPED"
		case LAUNCHER_STATE_READY:
			message = "RUNNING"
		case LAUNCHER_STATE_LOADING:
			message = "LOADING"
		}
	}

	return &messages.ServiceResp{
		Status:  200,
		Message: message,
	}, nil
}

func (cmds *CmdService) Launch(ctx context.Context, config *messages.LaunchConfig) (*messages.ServiceResp, error) {
	memAlloc := config.GetMemAlloc()
	minecraftLauncher = NewMinecraftLauncher(int(memAlloc))
	if err := minecraftLauncher.Launch(ctx); err != nil {
		return nil, err
	}
	return &messages.ServiceResp{
		Status:  200,
		Message: "Launching server!",
	}, nil
}

func (cmds *CmdService) Stop(ctx context.Context, _ *messages.EmptyReq) (*messages.ServiceResp, error) {
	if minecraftLauncher == nil {
		return nil, ErrLauncherOffline
	}
	if err := minecraftLauncher.Stop(ctx); err != nil {
		return nil, err
	}

	return &messages.ServiceResp{
		Status:  200,
		Message: "Stopping server!",
	}, nil
}

func Start() error {
	addr := fmt.Sprintf(":%d", DEFAULT_TCP_PORT)
	listener, err := net.Listen(DEFAULT_GRPC_TRANSPORT_PROTOCOL, addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	messages.RegisterCmdServiceServer(grpcServer, &CmdService{})

	log.Printf("Starting server on %s", addr)

	return grpcServer.Serve(listener)
}
