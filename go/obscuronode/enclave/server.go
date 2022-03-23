package enclave

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon/rpc"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon/rpc/generated"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"google.golang.org/grpc"

	"github.com/ethereum/go-ethereum/common"
)

// Receives RPC calls to the enclave process and relays them to the enclave.Enclave.
type server struct {
	generated.UnimplementedEnclaveProtoServer
	enclave nodecommon.Enclave
}

// StartServer starts a server on the given port on a separate thread. It creates an enclave.Enclave for the provided nodeID,
// and uses it to respond to incoming RPC messages from the host.
func StartServer(port uint64, nodeID common.Address, collector StatsCollector) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return &grpc.Server{}, fmt.Errorf("enclave RPC server could not listen on port: %w", err)
	}
	grpcServer := grpc.NewServer()
	enclaveServer := server{enclave: NewEnclave(nodeID, true, collector)}
	generated.RegisterEnclaveProtoServer(grpcServer, &enclaveServer)
	go func(lis net.Listener) {
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Log(fmt.Sprintf("enclave RPC server could not serve: %s", err))
		}
	}(lis)

	return grpcServer, nil
}

// IsReady returns a nil error to indicate that the server is ready.
func (s *server) IsReady(context.Context, *generated.IsReadyRequest) (*generated.IsReadyResponse, error) {
	return &generated.IsReadyResponse{}, nil
}

func (s *server) Attestation(context.Context, *generated.AttestationRequest) (*generated.AttestationResponse, error) {
	attestation := s.enclave.Attestation()
	msg := generated.AttestationReportMsg{Owner: attestation.Owner.Bytes()}
	return &generated.AttestationResponse{AttestationReportMsg: &msg}, nil
}

func (s *server) GenerateSecret(context.Context, *generated.GenerateSecretRequest) (*generated.GenerateSecretResponse, error) {
	secret := s.enclave.GenerateSecret()
	return &generated.GenerateSecretResponse{EncryptedSharedEnclaveSecret: secret}, nil
}

func (s *server) FetchSecret(_ context.Context, request *generated.FetchSecretRequest) (*generated.FetchSecretResponse, error) {
	attestationReport := rpc.FromAttestationReportMsg(request.AttestationReportMsg)
	secret := s.enclave.FetchSecret(attestationReport)
	return &generated.FetchSecretResponse{EncryptedSharedEnclaveSecret: secret}, nil
}

func (s *server) InitEnclave(_ context.Context, request *generated.InitEnclaveRequest) (*generated.InitEnclaveResponse, error) {
	s.enclave.InitEnclave(request.EncryptedSharedEnclaveSecret)
	return &generated.InitEnclaveResponse{}, nil
}

func (s *server) IsInitialised(context.Context, *generated.IsInitialisedRequest) (*generated.IsInitialisedResponse, error) {
	isInitialised := s.enclave.IsInitialised()
	return &generated.IsInitialisedResponse{IsInitialised: isInitialised}, nil
}

func (s *server) ProduceGenesis(context.Context, *generated.ProduceGenesisRequest) (*generated.ProduceGenesisResponse, error) {
	genesisRollup := s.enclave.ProduceGenesis()
	blockSubmissionResponse := rpc.ToBlockSubmissionResponseMsg(genesisRollup)
	return &generated.ProduceGenesisResponse{BlockSubmissionResponse: &blockSubmissionResponse}, nil
}

func (s *server) IngestBlocks(_ context.Context, request *generated.IngestBlocksRequest) (*generated.IngestBlocksResponse, error) {
	blocks := make([]*types.Block, 0)
	for _, encodedBlock := range request.EncodedBlocks {
		bl := decodeBlock(encodedBlock)
		blocks = append(blocks, &bl)
	}

	s.enclave.IngestBlocks(blocks)
	return &generated.IngestBlocksResponse{}, nil
}

func (s *server) Start(_ context.Context, request *generated.StartRequest) (*generated.StartResponse, error) {
	bl := decodeBlock(request.EncodedBlock)
	s.enclave.Start(bl)
	return &generated.StartResponse{}, nil
}

func (s *server) SubmitBlock(_ context.Context, request *generated.SubmitBlockRequest) (*generated.SubmitBlockResponse, error) {
	bl := decodeBlock(request.EncodedBlock)
	blockSubmissionResponse := s.enclave.SubmitBlock(bl)

	msg := rpc.ToBlockSubmissionResponseMsg(blockSubmissionResponse)
	return &generated.SubmitBlockResponse{BlockSubmissionResponse: &msg}, nil
}

func (s *server) SubmitRollup(_ context.Context, request *generated.SubmitRollupRequest) (*generated.SubmitRollupResponse, error) {
	extRollup := rpc.FromExtRollupMsg(request.ExtRollup)
	s.enclave.SubmitRollup(extRollup)
	return &generated.SubmitRollupResponse{}, nil
}

func (s *server) SubmitTx(_ context.Context, request *generated.SubmitTxRequest) (*generated.SubmitTxResponse, error) {
	err := s.enclave.SubmitTx(request.EncryptedTx)
	return &generated.SubmitTxResponse{}, err
}

func (s *server) Balance(_ context.Context, request *generated.BalanceRequest) (*generated.BalanceResponse, error) {
	balance := s.enclave.Balance(common.BytesToAddress(request.Address))
	return &generated.BalanceResponse{Balance: balance}, nil
}

func (s *server) RoundWinner(_ context.Context, request *generated.RoundWinnerRequest) (*generated.RoundWinnerResponse, error) {
	extRollup, winner := s.enclave.RoundWinner(common.BytesToHash(request.Parent))
	extRollupMsg := rpc.ToExtRollupMsg(&extRollup)
	return &generated.RoundWinnerResponse{Winner: winner, ExtRollup: &extRollupMsg}, nil
}

func (s *server) Stop(context.Context, *generated.StopRequest) (*generated.StopResponse, error) {
	s.enclave.Stop()
	return &generated.StopResponse{}, nil
}

func (s *server) GetTransaction(_ context.Context, request *generated.GetTransactionRequest) (*generated.GetTransactionResponse, error) {
	tx := s.enclave.GetTransaction(common.BytesToHash(request.TxHash))
	if tx == nil {
		return &generated.GetTransactionResponse{Known: false, EncodedTransaction: []byte{}}, nil
	}

	var buffer bytes.Buffer
	if err := tx.EncodeRLP(&buffer); err != nil {
		log.Log(fmt.Sprintf("failed to decode transaction sent to enclave: %v", err))
	}
	return &generated.GetTransactionResponse{Known: true, EncodedTransaction: buffer.Bytes()}, nil
}

func decodeBlock(encodedBlock []byte) types.Block {
	block := types.Block{}
	err := rlp.DecodeBytes(encodedBlock, &block)
	if err != nil {
		log.Log(fmt.Sprintf("failed to decode block sent to enclave: %v", err))
	}
	return block
}
