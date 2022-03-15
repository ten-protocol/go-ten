package rpc

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/obscuronet/obscuro-playground/integration/simulation"
	"log"
	"math/big"
	"net"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
	//"github.com/obscuronet/obscuro-playground/integration/simulation"

	"google.golang.org/grpc"
)

var Port = flag.Int("port", 50051, "The server port")

// TODO - Joel - Return errors as needed.
// TODO - Joel - Establish whether some gRPC methods can be declared without an '(x, error)' return type.

type server struct {
	UnimplementedEnclaveInternalServer
	enclave enclave.Enclave
}

// TODO - Joel - Add all other methods.

func (s *server) Attestation(ctx context.Context, request *AttestationRequest) (*AttestationResponse, error) {
	msg := AttestationReportMsg{Owner: s.enclave.Attestation().Owner.Bytes()}
	return &AttestationResponse{AttestationReportMsg: &msg}, nil
}

func (s *server) GenerateSecret(ctx context.Context, request *GenerateSecretRequest) (*GenerateSecretResponse, error) {
	secret := s.enclave.GenerateSecret()
	return &GenerateSecretResponse{EncryptedSharedEnclaveSecret: secret}, nil
}

func (s *server) FetchSecret(ctx context.Context, request *FetchSecretRequest) (*FetchSecretResponse, error) {
	attestationReport := toAttestationReport(request.AttestationReportMsg)
	secret := s.enclave.FetchSecret(attestationReport)
	return &FetchSecretResponse{EncryptedSharedEnclaveSecret: secret}, nil
}

func (s *server) Init(ctx context.Context, request *InitRequest) (*InitResponse, error) {
	s.enclave.Init(request.EncryptedSharedEnclaveSecret)
	return &InitResponse{}, nil
}

func (s *server) IsInitialised(ctx context.Context, request *IsInitialisedRequest) (*IsInitialisedResponse, error) {
	isInitialised := s.enclave.IsInitialised()
	return &IsInitialisedResponse{IsInitialised: isInitialised}, nil
}

func (s *server) ProduceGenesis(ctx context.Context, request *ProduceGenesisRequest) (*ProduceGenesisResponse, error) {
	blockSubmissionResponse := toBlockSubmissionResponseMsg(s.enclave.ProduceGenesis())
	return &ProduceGenesisResponse{BlockSubmissionResponse: &blockSubmissionResponse}, nil
}

func (s *server) IngestBlocks(ctx context.Context, request *IngestBlocksRequest) (*IngestBlocksResponse, error) {
	var blocks []*types.Block
	for _, encodedBlock := range request.EncodedBlocks {
		bl := types.Block{}
		rlp.DecodeBytes(encodedBlock, &bl)
		blocks = append(blocks, &bl)
	}

	s.enclave.IngestBlocks(blocks)
	return &IngestBlocksResponse{}, nil
}

func (s *server) Start(ctx context.Context, request *StartRequest) (*StartResponse, error) {
	bl := types.Block{}
	rlp.DecodeBytes(request.EncodedBlock, &bl)
	// TODO - Joel - Work out if we want to start the goroutine here. We probably do.
	go s.enclave.Start(bl)
	return &StartResponse{}, nil
}

func (s *server) SubmitBlock(ctx context.Context, request *SubmitBlockRequest) (*SubmitBlockResponse, error) {
	bl := types.Block{}
	rlp.DecodeBytes(request.EncodedBlock, &bl)
	s.enclave.SubmitBlock(bl)
	return &SubmitBlockResponse{}, nil
}

func (s *server) SubmitRollup(ctx context.Context, request *SubmitRollupRequest) (*SubmitRollupResponse, error) {
	extRollup := toExtRollup(request.ExtRollup)
	s.enclave.SubmitRollup(extRollup)
	return &SubmitRollupResponse{}, nil
}

func (s *server) SubmitTx(ctx context.Context, request *SubmitTxRequest) (*SubmitTxResponse, error) {
	err := s.enclave.SubmitTx(request.EncryptedTx)
	return &SubmitTxResponse{}, err
}

func (s *server) Balance(ctx context.Context, request *BalanceRequest) (*BalanceResponse, error) {
	balance := s.enclave.Balance(common.BytesToAddress(request.Address))
	return &BalanceResponse{Balance: balance}, nil
}

func (s *server) RoundWinner(ctx context.Context, request *RoundWinnerRequest) (*RoundWinnerResponse, error) {
	extRollup, winner := s.enclave.RoundWinner(common.BytesToHash(request.Parent))
	extRollupMsg := toExtRollupMsg(&extRollup)
	return &RoundWinnerResponse{Winner: winner, ExtRollup: &extRollupMsg}, nil
}

func (s *server) Stop(ctx context.Context, request *StopRequest) (*StopResponse, error) {
	s.enclave.Stop()
	return &StopResponse{}, nil
}

func (s *server) GetTransaction(ctx context.Context, request *GetTransactionRequest) (*GetTransactionResponse, error) {
	tx, unknown := s.enclave.GetTransaction(common.BytesToHash(request.TxHash))
	var buffer bytes.Buffer
	tx.EncodeRLP(&buffer)
	return &GetTransactionResponse{Unknown: unknown, EncodedTransaction: buffer.Bytes()}, nil
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
