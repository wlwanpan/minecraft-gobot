package mcs

import (
	"context"
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

type server struct {
	wpr *wrapper
}

func (s *server) listen(p int) error {
	addr := fmt.Sprintf(":%d", p)
	listener, err := net.Listen(DEFAULT_GRPC_TRANSPORT_PROTOCOL, addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	services.RegisterLauncherServiceServer(grpcServer, s)

	return grpcServer.Serve(listener)
}

func (s *server) Status(ctx context.Context, _ *services.EmptyReq) (*services.StatusResp, error) {
	var state wrapperState
	if s.wpr == nil {
		state = WRAPPER_STATE_OFFLINE
	} else {
		state = s.wpr.state
	}

	return &services.StatusResp{
		ServerState: wrapperStateMap[state],
	}, nil
}

func (s *server) Start(ctx context.Context, config *services.StartConfig) (*services.ServiceResp, error) {
	memAlloc := int(config.GetMemAlloc())

	if s.wpr == nil {
		s.wpr = newWrapper()
	}

	if err := s.wpr.start(memAlloc); err != nil {
		return &services.ServiceResp{
			Message: err.Error(),
			Status:  500,
		}, nil
	}
	return &services.ServiceResp{
		Message: "Starting server!",
		Status:  200,
	}, nil
}

func (s *server) Stop(ctx context.Context, _ *services.EmptyReq) (*services.ServiceResp, error) {
	if s.wpr == nil {
		return &services.ServiceResp{
			Message: "Server already offline",
			Status:  500,
		}, nil
	}

	if err := s.wpr.stop(); err != nil {
		return &services.ServiceResp{
			Message: err.Error(),
			Status:  500,
		}, nil
	}

	s.wpr = nil

	return &services.ServiceResp{
		Message: "Stopping server!",
		Status:  200,
	}, nil
}

func Start() error {
	s := &server{}

	log.Printf("Starting mcs on: %d", DEFAULT_TCP_PORT)
	return s.listen(DEFAULT_TCP_PORT)
}
