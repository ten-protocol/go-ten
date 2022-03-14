package rpc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	"github.com/obscuronet/obscuro-playground/integration/simulation"

	"google.golang.org/grpc"
)

var Port = flag.Int("port", 50051, "The server port")

type server struct {
	UnimplementedEnclaveServer
	enclave enclave.Enclave
}

func (s *server) Start(ctx context.Context, request *StartRequest) (*StartResponse, error) {
	blockAddress := common.BigToAddress(big.NewInt(int64(0)))
	block := obscurocommon.NewBlock(nil, 0, blockAddress, []*obscurocommon.L1Tx{})
	go s.enclave.Start(*block)
	return &StartResponse{}, errors.New("server received request")
}

func StartServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	// TODO - Joel - Pass in real arguments here. These are just dummies.
	enclaveAddress := common.BigToAddress(big.NewInt(int64(1)))
	enclaveServer := server{enclave: enclave.NewEnclave(enclaveAddress, true, simulation.NewStats(1))}
	RegisterEnclaveServer(grpcServer, &enclaveServer)
	grpcServer.Serve(lis)
}
