package rpc

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var Port = flag.Int("port", 50051, "The server port")

type server struct {
	UnimplementedEnclaveServer
}

//func (s *server) GetFeature(ctx context.Context, point *Point) (*Point, error) {
//	return &Point{Latitude: 777, Longitude: 777}, nil
//}

func StartServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	RegisterEnclaveServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}
