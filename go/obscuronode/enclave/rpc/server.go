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

// TODO - Joel - Return errors as needed.
// TODO - Joel - Establish whether some gRPC methods can be declared without an '(x, error)' return type.

type server struct {
	UnimplementedEnclaveInternalServer
	enclave enclave.Enclave
}

func StartServer(nodeID common.Address, port uint64, collector enclave.StatsCollector) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Log(fmt.Sprintf("enclave RPC server failed to listen: %v", err))
		return
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	enclaveServer := server{enclave: enclave.NewEnclave(nodeID, true, collector)}
	RegisterEnclaveInternalServer(grpcServer, &enclaveServer)
	go func(lis net.Listener) {
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Log(fmt.Sprintf("enclave RPC server could not serve: %s", err))
		}
	}(lis)
}

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
	blocks := make([]*types.Block, 0)
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
	go s.enclave.Start(bl)
	return &StartResponse{}, nil
}

func (s *server) SubmitBlock(ctx context.Context, request *SubmitBlockRequest) (*SubmitBlockResponse, error) {
	bl := types.Block{}
	rlp.DecodeBytes(request.EncodedBlock, &bl)
	blockSubmissionResponse := s.enclave.SubmitBlock(bl)

	msg := toBlockSubmissionResponseMsg(blockSubmissionResponse)
	return &SubmitBlockResponse{BlockSubmissionResponse: &msg}, nil
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
	tx, known := s.enclave.GetTransaction(common.BytesToHash(request.TxHash))
	var buffer bytes.Buffer
	tx.EncodeRLP(&buffer)
	return &GetTransactionResponse{Known: known, EncodedTransaction: buffer.Bytes()}, nil
}
