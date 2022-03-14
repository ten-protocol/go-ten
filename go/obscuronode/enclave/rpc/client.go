package rpc

import (
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")

func getGetFeature(client EnclaveClient) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// inputPoint := Point{Latitude: 1, Longitude: 1}
	// point, _ := client.GetFeature(ctx, &inputPoint)
	// fmt.Println(point.Longitude)
}

func StartClient() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := NewEnclaveClient(conn)

	getGetFeature(client)
}
