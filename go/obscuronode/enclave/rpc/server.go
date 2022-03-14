package rpc

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var Port = flag.Int("port", 50051, "The server port")

type server struct {
	UnimplementedRouteGuideServer
}

func (s *server) GetFeature(ctx context.Context, point *Point) (*Point, error) {
	return &Point{Latitude: 777, Longitude: 777}, nil
}

func StartServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	RegisterRouteGuideServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
