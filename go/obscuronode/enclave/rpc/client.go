package rpc

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func getGetFeature(client RouteGuideClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	inputPoint := Point{Latitude: 1, Longitude: 1}
	point, _ := client.GetFeature(ctx, &inputPoint)
	fmt.Println(point.Longitude)
}

func StartClient() {
	fmt.Println("client started")
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := NewRouteGuideClient(conn)

	getGetFeature(client)
}
