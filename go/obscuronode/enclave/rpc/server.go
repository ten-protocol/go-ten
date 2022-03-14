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
	UnimplementedRouteGuideServer
}

func GetFeature(point Point) Point {
	fmt.Println(point)
	return Point{Latitude: 1, Longitude: 2}
}

func Start() {
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
