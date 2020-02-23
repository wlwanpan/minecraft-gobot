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
	wpr    *wrapper
	backer *backer
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

	var message string
	var state string
	if s.wpr == nil {
		state = WRAPPER_STATE_OFFLINE
	} else {
		state = s.wpr.stateMachine.Current()
		message = s.wpr.lastLogLine
	}

	return &services.StatusResp{
		ServerState: state,
		Message:     message,
	}, nil
}

func (s *server) Backup(ctx context.Context, cfg *services.EmptyReq) (*services.BackupResp, error) {
	log.Printf("Request received: Backup")
	if s.wpr != nil && !s.wpr.isOffline() {
		return &services.BackupResp{
			Status:  services.BackupStatus_FAILED,
			Message: "server must be offline to perform a backup",
		}, nil
	}

	if s.backer == nil {
		log.Println("Backer offline, creating a new one")
		s.backer = newBacker()

		log.Println("Starting backup!")
		s.backer.start()
	}

	var url string
	var message string
	var status services.BackupStatus

	switch s.backer.state {
	case BACKER_STATE_DONE:
		status = services.BackupStatus_DONE
		url = s.backer.lastUrl
		s.backer = nil
	case BACKER_STATE_FAILED:
		status = services.BackupStatus_FAILED
		message = "failed to perform backup"
	case BACKER_STATE_ZIPPING:
		status = services.BackupStatus_ZIPPING
		message = "compressing world"
	case BACKER_STATE_UPLOADING:
		status = services.BackupStatus_UPLOADING
		message = "uploading world"
	}

	return &services.BackupResp{
		Status:  status,
		Message: message,
		LinkUrl: url,
	}, nil
}

func (s *server) Start(ctx context.Context, cfg *services.StartConfig) (*services.ServiceResp, error) {
	memAlloc := int(cfg.GetMemAlloc())
	log.Printf("Request received: Start, mem-alloc=%d", memAlloc)

	if s.wpr == nil {
		s.wpr = newWrapper()
	}

	if err := s.wpr.start(memAlloc); err != nil {
		return &services.ServiceResp{Message: err.Error()}, nil
	}
	return &services.ServiceResp{Message: "Starting server!"}, nil
}

func (s *server) Stop(ctx context.Context, _ *services.EmptyReq) (*services.ServiceResp, error) {
	log.Println("Request received: Stop")

	if s.wpr == nil {
		return &services.ServiceResp{Message: "Server already offline"}, nil
	}

	if err := s.wpr.stop(); err != nil {
		return &services.ServiceResp{Message: err.Error()}, nil
	}

	s.wpr = nil

	return &services.ServiceResp{Message: "Stopping server!"}, nil
}

func Start(port int) error {
	s := &server{}

	log.Printf("Starting mcs on: %d", port)
	return s.listenAndServe(port)
}
