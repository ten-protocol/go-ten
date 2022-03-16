package rpc

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"google.golang.org/grpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"
)

// Receives RPC calls to the enclave process and relays them to the enclave.Enclave.
type server struct {
	UnimplementedEnclaveProtoServer
	enclave enclave.Enclave
}

// StartServer starts a server on the given port on a separate thread. It creates an enclave.Enclave for the provided nodeID,
// and uses it to respond to incoming RPC messages from the host.
func StartServer(port uint64, nodeID common.Address, collector enclave.StatsCollector) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return &grpc.Server{}, fmt.Errorf("enclave RPC server could not listen on port: %w", err)
	}
	grpcServer := grpc.NewServer()
	enclaveServer := server{enclave: enclave.NewEnclave(nodeID, true, collector)}
	RegisterEnclaveProtoServer(grpcServer, &enclaveServer)
	go func(lis net.Listener) {
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Log(fmt.Sprintf("enclave RPC server could not serve: %s", err))
		}
	}(lis)

	return grpcServer, nil
}

// IsReady returns a nil error to indicate that the server is ready.
func (s *server) IsReady(context.Context, *IsReadyRequest) (*IsReadyResponse, error) {
	return &IsReadyResponse{}, nil
}

func (s *server) Attestation(context.Context, *AttestationRequest) (*AttestationResponse, error) {
	msg := AttestationReportMsg{Owner: s.enclave.Attestation().Owner.Bytes()}
	return &AttestationResponse{AttestationReportMsg: &msg}, nil
}

func (s *server) GenerateSecret(context.Context, *GenerateSecretRequest) (*GenerateSecretResponse, error) {
	secret := s.enclave.GenerateSecret()
	return &GenerateSecretResponse{EncryptedSharedEnclaveSecret: secret}, nil
}

func (s *server) FetchSecret(_ context.Context, request *FetchSecretRequest) (*FetchSecretResponse, error) {
	attestationReport := fromAttestationReportMsg(request.AttestationReportMsg)
	secret := s.enclave.FetchSecret(attestationReport)
	return &FetchSecretResponse{EncryptedSharedEnclaveSecret: secret}, nil
}

func (s *server) InitEnclave(_ context.Context, request *InitEnclaveRequest) (*InitEnclaveResponse, error) {
	s.enclave.InitEnclave(request.EncryptedSharedEnclaveSecret)
	return &InitEnclaveResponse{}, nil
}

func (s *server) IsInitialised(context.Context, *IsInitialisedRequest) (*IsInitialisedResponse, error) {
	isInitialised := s.enclave.IsInitialised()
	return &IsInitialisedResponse{IsInitialised: isInitialised}, nil
}

func (s *server) ProduceGenesis(context.Context, *ProduceGenesisRequest) (*ProduceGenesisResponse, error) {
	blockSubmissionResponse := toBlockSubmissionResponseMsg(s.enclave.ProduceGenesis())
	return &ProduceGenesisResponse{BlockSubmissionResponse: &blockSubmissionResponse}, nil
}

func (s *server) IngestBlocks(_ context.Context, request *IngestBlocksRequest) (*IngestBlocksResponse, error) {
	blocks := make([]*types.Block, 0)
	for _, encodedBlock := range request.EncodedBlocks {
		bl := decodeBlock(encodedBlock)
		blocks = append(blocks, &bl)
	}

	s.enclave.IngestBlocks(blocks)
	return &IngestBlocksResponse{}, nil
}

func (s *server) Start(_ context.Context, request *StartRequest) (*StartResponse, error) {
	bl := decodeBlock(request.EncodedBlock)
	go s.enclave.Start(bl)
	return &StartResponse{}, nil
}

func (s *server) SubmitBlock(_ context.Context, request *SubmitBlockRequest) (*SubmitBlockResponse, error) {
	bl := decodeBlock(request.EncodedBlock)
	blockSubmissionResponse := s.enclave.SubmitBlock(bl)

	msg := toBlockSubmissionResponseMsg(blockSubmissionResponse)
	return &SubmitBlockResponse{BlockSubmissionResponse: &msg}, nil
}

func (s *server) SubmitRollup(_ context.Context, request *SubmitRollupRequest) (*SubmitRollupResponse, error) {
	extRollup := fromExtRollupMsg(request.ExtRollup)
	s.enclave.SubmitRollup(extRollup)
	return &SubmitRollupResponse{}, nil
}

func (s *server) SubmitTx(_ context.Context, request *SubmitTxRequest) (*SubmitTxResponse, error) {
	err := s.enclave.SubmitTx(request.EncryptedTx)
	return &SubmitTxResponse{}, err
}

func (s *server) Balance(_ context.Context, request *BalanceRequest) (*BalanceResponse, error) {
	balance := s.enclave.Balance(common.BytesToAddress(request.Address))
	return &BalanceResponse{Balance: balance}, nil
}

func (s *server) RoundWinner(_ context.Context, request *RoundWinnerRequest) (*RoundWinnerResponse, error) {
	extRollup, winner := s.enclave.RoundWinner(common.BytesToHash(request.Parent))
	extRollupMsg := toExtRollupMsg(&extRollup)
	return &RoundWinnerResponse{Winner: winner, ExtRollup: &extRollupMsg}, nil
}

func (s *server) Stop(context.Context, *StopRequest) (*StopResponse, error) {
	s.enclave.Stop()
	return &StopResponse{}, nil
}

func (s *server) GetTransaction(_ context.Context, request *GetTransactionRequest) (*GetTransactionResponse, error) {
	tx := s.enclave.GetTransaction(common.BytesToHash(request.TxHash))
	if tx == nil {
		return &GetTransactionResponse{Known: false, EncodedTransaction: []byte{}}, nil
	}

	var buffer bytes.Buffer
	if err := tx.EncodeRLP(&buffer); err != nil {
		log.Log(fmt.Sprintf("failed to decode transaction sent to enclave: %v", err))
	}
	return &GetTransactionResponse{Known: true, EncodedTransaction: buffer.Bytes()}, nil
}

func decodeBlock(encodedBlock []byte) types.Block {
	block := types.Block{}
	err := rlp.DecodeBytes(encodedBlock, &block)
	if err != nil {
		log.Log(fmt.Sprintf("failed to decode block sent to enclave: %v", err))
	}
	return block
}
