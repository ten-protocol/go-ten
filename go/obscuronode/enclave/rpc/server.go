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

// TODO - Joel - Establish whether some gRPC methods can be declared without an '(x, error)' return type.

type server struct {
	UnimplementedEnclaveInternalServer
	enclave enclave.Enclave
}

func (s *server) Init(ctx context.Context, request *InitRequest) (*InitResponse, error) {
	secret := []byte("I took some library books with me when I moved to France")
	s.enclave.Init(secret)
	return &InitResponse{}, nil
}

func (s *server) Start(ctx context.Context, request *StartRequest) (*StartResponse, error) {
	blockAddress := common.BigToAddress(big.NewInt(int64(0)))
	block := obscurocommon.NewBlock(nil, 0, blockAddress, []*obscurocommon.L1Tx{})
	go s.enclave.Start(*block)
	return &StartResponse{}, errors.New("server received request")
}

// TODO - Joel - Have the server start on different ports, not just the same one repeatedly.
// TODO - Joel - Pass in real arguments to create an enclave here. We just use dummies below.
func StartServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	enclaveAddress := common.BigToAddress(big.NewInt(int64(1)))
	enclaveServer := server{enclave: enclave.NewEnclave(enclaveAddress, true, simulation.NewStats(1))}
	RegisterEnclaveInternalServer(grpcServer, &enclaveServer)
	grpcServer.Serve(lis)
}
