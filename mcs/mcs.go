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
)

type server struct {
	wpr *wrapper
}

func (s *server) listenAndServe(p int) error {
	addr := fmt.Sprintf(":%d", p)
	listener, err := net.Listen(DEFAULT_GRPC_TRANSPORT_PROTOCOL, addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	services.RegisterMcsServiceServer(grpcServer, s)

	return grpcServer.Serve(listener)
}

func (s *server) Ping(ctx context.Context, _ *services.PingReq) (*services.PongResp, error) {
	log.Println("Request received: Ping")
	return &services.PongResp{}, nil
}

func (s *server) Status(ctx context.Context, _ *services.EmptyReq) (*services.StatusResp, error) {
	log.Println("Request received: Status")

	var state wrapperState
	var message string
	if s.wpr == nil {
		state = WRAPPER_STATE_OFFLINE
	} else {
		state = s.wpr.state
		message = s.wpr.lastLogLine
	}

	return &services.StatusResp{
		ServerState: wrapperStateMap[state],
		Message:     message,
	}, nil
}

func (s *server) Backup(ctx context.Context, cfg *services.EmptyReq) (*services.BackupResp, error) {
	return &services.BackupResp{}, nil
}

func (s *server) Start(ctx context.Context, cfg *services.StartConfig) (*services.ServiceResp, error) {
	memAlloc := int(cfg.GetMemAlloc())
	log.Printf("Request received: Start, mem-alloc=%d", memAlloc)

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
	log.Println("Request received: Stop")

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

func Start(port int) error {
	s := &server{}

	log.Printf("Starting mcs on: %d", port)
	return s.listenAndServe(port)
}
