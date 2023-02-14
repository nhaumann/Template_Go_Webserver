package healthcheck_client

import (
	"context"
	"log"
	"strconv"

	pb "webserver/pkg/healthcheck"

	hcsrver "webserver/pkg/healthcheck/healthcheck_server"

	grpc "google.golang.org/grpc"
	ins "google.golang.org/grpc/credentials/insecure"
)

func Connect(config hcsrver.HealthCheckServerConfig) {
	if config.WebPort == "" {
		log.Fatal("WebPort is empty")
	}

	conn, err := grpc.Dial(config.WebPort, grpc.WithTransportCredentials(ins.NewCredentials()))
	if err != nil {
		log.Fatalf("can not connect: %v", err)

	}
	defer conn.Close()
	c := pb.NewHealthCheckClient(conn)
	r, err := c.IsAlive(context.Background(), &pb.HealthCheckRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println("HealthCheckClient:" + strconv.FormatBool(r.Alive))

}
