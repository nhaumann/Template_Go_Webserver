package healthcheck_server

import (
	"context"
	"log"
	"net"

	pb "glossary/pkg/healthcheck"

	grpc "google.golang.org/grpc"
)

type HealthCheckServerConfig struct {
	WebPort string
}

type HealthCheckServer struct {
	pb.UnimplementedHealthCheckServer
}

func (s *HealthCheckServer) IsAlive(ctx context.Context, in *pb.HealthCheckRequest) (*pb.Alive, error) {
	return &pb.Alive{Alive: true}, nil
}

func Serve(config HealthCheckServerConfig) {

	if config.WebPort == "" {
		log.Fatal("WebPort is empty")
	}

	//serve webserver in the background, allow user input to continue
	lis, err := net.Listen("tcp", config.WebPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHealthCheckServer(s, &HealthCheckServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
