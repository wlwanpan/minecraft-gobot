package launcher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/wlwanpan/minecraft-gobot/services"
	"google.golang.org/grpc"
)

const (
	DEFAULT_GRPC_TRANSPORT_PROTOCOL = "tcp"

	DEFAULT_TCP_PORT = 7777
)

var (
	ErrLauncherOffline = errors.New("launcher is offline")
)

type grpcService struct {
	launcher *launcher
}

func (s *grpcService) Status(ctx context.Context, _ *services.EmptyReq) (*services.StatusResp, error) {
	var state string
	if s.launcher == nil {
		state = "OFFLINE"
	} else {
		switch s.launcher.currState {
		case LAUNCHER_STATE_INIT:
			state = "STOPPED"
		case LAUNCHER_STATE_READY:
			state = "RUNNING"
		case LAUNCHER_STATE_LOADING:
			state = "LOADING"
		}
	}

	return &services.StatusResp{
		ServerState: state,
		Message:     "",
	}, nil
}

func (s *grpcService) Launch(ctx context.Context, config *services.LaunchConfig) (*services.ServiceResp, error) {
	memAlloc := config.GetMemAlloc()
	s.launcher = newLauncher(int(memAlloc))
	if err := s.launcher.Launch(ctx); err != nil {
		return nil, err
	}
	return &services.ServiceResp{
		Status:  200,
		Message: "Launching server!",
	}, nil
}

func (s *grpcService) Stop(ctx context.Context, _ *services.EmptyReq) (*services.ServiceResp, error) {
	if s.launcher == nil {
		return nil, ErrLauncherOffline
	}
	if err := s.launcher.Stop(ctx); err != nil {
		return nil, err
	}

	return &services.ServiceResp{
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
	services.RegisterLauncherServiceServer(grpcServer, &grpcService{})

	log.Printf("Starting server on %s", addr)

	return grpcServer.Serve(listener)
}
