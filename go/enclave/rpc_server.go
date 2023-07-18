package enclave

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/rpc"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"google.golang.org/grpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

// RPCServer receives RPC calls to the enclave process and relays them to the enclave.Enclave.
type RPCServer struct {
	generated.UnimplementedEnclaveProtoServer
	enclave       common.Enclave
	grpcServer    *grpc.Server
	logger        gethlog.Logger
	listenAddress string
}

// NewEnclaveRPCServer prepares an enclave RPCServer (doesn't start listening until `StartServer` is called
func NewEnclaveRPCServer(listenAddress string, enclave common.Enclave, logger gethlog.Logger) *RPCServer {
	return &RPCServer{
		enclave:       enclave,
		grpcServer:    grpc.NewServer(),
		logger:        logger,
		listenAddress: listenAddress,
	}
}

// StartServer starts a RPCServer on the given port on a separate thread. It creates an enclave.Enclave for the provided nodeID,
// and uses it to respond to incoming RPC messages from the host.
func (s *RPCServer) StartServer() error {
	lis, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return fmt.Errorf("RPCServer could not listen on port: %w", err)
	}
	generated.RegisterEnclaveProtoServer(s.grpcServer, s)

	go func(lis net.Listener) {
		s.logger.Info(fmt.Sprintf("RPCServer listening on address %s.", s.listenAddress))
		err = s.grpcServer.Serve(lis)
		if err != nil {
			s.logger.Info("RPCServer could not serve", log.ErrKey, err)
		}
	}(lis)

	return nil
}

// Status returns the current status of the RPCServer as an enum value (see common.Status for details)
func (s *RPCServer) Status(context.Context, *generated.StatusRequest) (*generated.StatusResponse, error) {
	status, sysError := s.enclave.Status()
	return &generated.StatusResponse{
		StatusCode:  int32(status.StatusCode),
		L1Head:      status.L1Head.Bytes(),
		L2Head:      status.L2Head.Bytes(),
		SystemError: toRPCError(sysError),
	}, nil
}

func (s *RPCServer) Attestation(context.Context, *generated.AttestationRequest) (*generated.AttestationResponse, error) {
	attestation, sysError := s.enclave.Attestation()
	if sysError != nil {
		return &generated.AttestationResponse{SystemError: toRPCError(sysError)}, nil
	}
	msg := rpc.ToAttestationReportMsg(attestation)
	return &generated.AttestationResponse{AttestationReportMsg: &msg}, nil
}

func (s *RPCServer) GenerateSecret(context.Context, *generated.GenerateSecretRequest) (*generated.GenerateSecretResponse, error) {
	secret, sysError := s.enclave.GenerateSecret()
	if sysError != nil {
		return &generated.GenerateSecretResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.GenerateSecretResponse{EncryptedSharedEnclaveSecret: secret}, nil
}

func (s *RPCServer) InitEnclave(_ context.Context, request *generated.InitEnclaveRequest) (*generated.InitEnclaveResponse, error) {
	sysError := s.enclave.InitEnclave(request.EncryptedSharedEnclaveSecret)
	return &generated.InitEnclaveResponse{SystemError: toRPCError(sysError)}, nil
}

func (s *RPCServer) SubmitL1Block(_ context.Context, request *generated.SubmitBlockRequest) (*generated.SubmitBlockResponse, error) {
	bl := s.decodeBlock(request.EncodedBlock)
	receipts := s.decodeReceipts(request.EncodedReceipts)
	blockSubmissionResponse, err := s.enclave.SubmitL1Block(bl, receipts, request.IsLatest)
	if err != nil {
		var rejErr *errutil.BlockRejectError
		isReject := errors.As(err, &rejErr)
		if isReject {
			return &generated.SubmitBlockResponse{
				BlockSubmissionResponse: &generated.BlockSubmissionResponseMsg{
					Error: &generated.BlockSubmissionErrorMsg{
						Cause:  rejErr.Wrapped.Error(),
						L1Head: rejErr.L1Head.Bytes(),
					},
				},
			}, nil
		}
		return nil, err
	}

	msg, err := rpc.ToBlockSubmissionResponseMsg(blockSubmissionResponse)
	if err != nil {
		return nil, err
	}
	return &generated.SubmitBlockResponse{BlockSubmissionResponse: msg}, nil
}

func (s *RPCServer) SubmitTx(_ context.Context, request *generated.SubmitTxRequest) (*generated.SubmitTxResponse, error) {
	enclaveResponse, sysError := s.enclave.SubmitTx(request.EncryptedTx)
	if sysError != nil {
		return &generated.SubmitTxResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.SubmitTxResponse{EncodedEnclaveResponse: enclaveResponse.Encode()}, nil
}

func (s *RPCServer) SubmitBatch(_ context.Context, request *generated.SubmitBatchRequest) (*generated.SubmitBatchResponse, error) {
	batch := rpc.FromExtBatchMsg(request.Batch)
	sysError := s.enclave.SubmitBatch(batch)
	return &generated.SubmitBatchResponse{SystemError: toRPCError(sysError)}, nil
}

func (s *RPCServer) ObsCall(_ context.Context, request *generated.ObsCallRequest) (*generated.ObsCallResponse, error) {
	enclaveResp, sysError := s.enclave.ObsCall(request.EncryptedParams)
	if sysError != nil {
		return &generated.ObsCallResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.ObsCallResponse{EncodedEnclaveResponse: enclaveResp.Encode()}, nil
}

func (s *RPCServer) GetTransactionCount(_ context.Context, request *generated.GetTransactionCountRequest) (*generated.GetTransactionCountResponse, error) {
	enclaveResp, sysError := s.enclave.GetTransactionCount(request.EncryptedParams)
	if sysError != nil {
		return &generated.GetTransactionCountResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.GetTransactionCountResponse{EncodedEnclaveResponse: enclaveResp.Encode()}, nil
}

func (s *RPCServer) Stop(context.Context, *generated.StopRequest) (*generated.StopResponse, error) {
	// stop the grpcServer on its own goroutine to avoid killing the existing connection
	go s.grpcServer.GracefulStop()
	return &generated.StopResponse{SystemError: toRPCError(s.enclave.Stop())}, nil
}

func (s *RPCServer) GetTransaction(_ context.Context, request *generated.GetTransactionRequest) (*generated.GetTransactionResponse, error) {
	enclaveResp, sysError := s.enclave.GetTransaction(request.EncryptedParams)
	if sysError != nil {
		return &generated.GetTransactionResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.GetTransactionResponse{EncodedEnclaveResponse: enclaveResp.Encode()}, nil
}

func (s *RPCServer) GetTransactionReceipt(_ context.Context, request *generated.GetTransactionReceiptRequest) (*generated.GetTransactionReceiptResponse, error) {
	enclaveResponse, sysError := s.enclave.GetTransactionReceipt(request.EncryptedParams)
	if sysError != nil {
		return &generated.GetTransactionReceiptResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.GetTransactionReceiptResponse{EncodedEnclaveResponse: enclaveResponse.Encode()}, nil
}

func (s *RPCServer) GetBalance(_ context.Context, request *generated.GetBalanceRequest) (*generated.GetBalanceResponse, error) {
	enclaveResp, sysError := s.enclave.GetBalance(request.EncryptedParams)
	if sysError != nil {
		return &generated.GetBalanceResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.GetBalanceResponse{EncodedEnclaveResponse: enclaveResp.Encode()}, nil
}

func (s *RPCServer) GetCode(_ context.Context, request *generated.GetCodeRequest) (*generated.GetCodeResponse, error) {
	address := gethcommon.BytesToAddress(request.Address)
	rollupHash := gethcommon.BytesToHash(request.RollupHash)

	code, sysError := s.enclave.GetCode(address, &rollupHash)
	if sysError != nil {
		return &generated.GetCodeResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.GetCodeResponse{Code: code}, nil
}

func (s *RPCServer) Subscribe(_ context.Context, req *generated.SubscribeRequest) (*generated.SubscribeResponse, error) {
	sysError := s.enclave.Subscribe(gethrpc.ID(req.Id), req.EncryptedSubscription)
	return &generated.SubscribeResponse{SystemError: toRPCError(sysError)}, nil
}

func (s *RPCServer) Unsubscribe(_ context.Context, req *generated.UnsubscribeRequest) (*generated.UnsubscribeResponse, error) {
	sysError := s.enclave.Unsubscribe(gethrpc.ID(req.Id))
	return &generated.UnsubscribeResponse{SystemError: toRPCError(sysError)}, nil
}

func (s *RPCServer) EstimateGas(_ context.Context, req *generated.EstimateGasRequest) (*generated.EstimateGasResponse, error) {
	enclaveResp, sysError := s.enclave.EstimateGas(req.EncryptedParams)
	if sysError != nil {
		return &generated.EstimateGasResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.EstimateGasResponse{EncodedEnclaveResponse: enclaveResp.Encode()}, nil
}

func (s *RPCServer) GetLogs(_ context.Context, req *generated.GetLogsRequest) (*generated.GetLogsResponse, error) {
	enclaveResp, sysError := s.enclave.GetLogs(req.EncryptedParams)
	if sysError != nil {
		return &generated.GetLogsResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.GetLogsResponse{EncodedEnclaveResponse: enclaveResp.Encode()}, nil
}

func (s *RPCServer) HealthCheck(_ context.Context, _ *generated.EmptyArgs) (*generated.HealthCheckResponse, error) {
	healthy, sysError := s.enclave.HealthCheck()
	if sysError != nil {
		return &generated.HealthCheckResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.HealthCheckResponse{Status: healthy}, nil
}

func (s *RPCServer) CreateRollup(_ context.Context, _ *generated.CreateRollupRequest) (*generated.CreateRollupResponse, error) {
	rollup, err := s.enclave.CreateRollup()

	msg := rpc.ToExtRollupMsg(rollup)

	return &generated.CreateRollupResponse{
		Msg:         &msg,
		SystemError: toRPCError(err),
	}, nil
}

func (s *RPCServer) CreateBatch(_ context.Context, _ *generated.CreateBatchRequest) (*generated.CreateBatchResponse, error) {
	err := s.enclave.CreateBatch()
	return &generated.CreateBatchResponse{}, err
}

func (s *RPCServer) DebugTraceTransaction(_ context.Context, req *generated.DebugTraceTransactionRequest) (*generated.DebugTraceTransactionResponse, error) {
	txHash := gethcommon.BytesToHash(req.TxHash)
	var config tracers.TraceConfig

	err := json.Unmarshal(req.Config, &config)
	if err != nil {
		return &generated.DebugTraceTransactionResponse{
			SystemError: toRPCError(fmt.Errorf("unable to unmarshall config - %w", err)),
		}, nil
	}

	traceTx, err := s.enclave.DebugTraceTransaction(txHash, &config)
	return &generated.DebugTraceTransactionResponse{Msg: string(traceTx), SystemError: toRPCError(err)}, nil
}

func (s *RPCServer) GetBatch(_ context.Context, request *generated.GetBatchRequest) (*generated.GetBatchResponse, error) {
	batch, err := s.enclave.GetBatch(gethcommon.BytesToHash(request.KnownHead))
	if err != nil {
		return nil, err
	}

	encodedBatch, encodingErr := batch.Encoded()
	return &generated.GetBatchResponse{
		Batch: encodedBatch,
		SystemError: &generated.SystemError{
			ErrorCode:   2,
			ErrorString: encodingErr.Error(),
		},
	}, err
}

func (s *RPCServer) StreamL2Updates(request *generated.StreamL2UpdatesRequest, stream generated.EnclaveProto_StreamL2UpdatesServer) error {
	var fromHash *common.L2BatchHash
	if request.KnownHead != nil {
		knownHead := gethcommon.BytesToHash(request.KnownHead)
		fromHash = &knownHead
	}

	batchChan, stop := s.enclave.StreamL2Updates(fromHash)
	defer stop()

	for {
		batchResp, ok := <-batchChan
		if !ok {
			s.logger.Info("Enclave closed batch channel.")
			break
		}

		encoded, err := json.Marshal(batchResp)
		if err != nil {
			s.logger.Error("Error marshalling batch response", log.ErrKey, err)
			return nil
		}

		if err := stream.Send(&generated.EncodedUpdateResponse{Batch: encoded}); err != nil {
			s.logger.Info("Failed streaming batch back to client", log.ErrKey, err)
			// not quite sure there is any point to this, we failed to send a batch
			// so error will probably not get sent either.
			return err
		}
	}

	return nil
}

func (s *RPCServer) DebugEventLogRelevancy(_ context.Context, req *generated.DebugEventLogRelevancyRequest) (*generated.DebugEventLogRelevancyResponse, error) {
	txHash := gethcommon.BytesToHash(req.TxHash)

	logs, err := s.enclave.DebugEventLogRelevancy(txHash)

	return &generated.DebugEventLogRelevancyResponse{Msg: string(logs), SystemError: toRPCError(err)}, nil
}

func (s *RPCServer) GetTotalContractCount(_ context.Context, req *generated.GetTotalContractCountRequest) (*generated.GetTotalContractCountResponse, error) {
	count, err := s.enclave.GetTotalContractCount()

	if count == nil {
		count = big.NewInt(0)
	}

	return &generated.GetTotalContractCountResponse{
		Count:       count.Int64(),
		SystemError: toRPCError(err),
	}, nil
}

func (s *RPCServer) decodeBlock(encodedBlock []byte) types.Block {
	block := types.Block{}
	err := rlp.DecodeBytes(encodedBlock, &block)
	if err != nil {
		s.logger.Info("failed to decode block sent to enclave", log.ErrKey, err)
	}
	return block
}

// decodeReceipts - converts the rlp encoded bytes to receipts if possible.
func (s *RPCServer) decodeReceipts(encodedReceipts []byte) types.Receipts {
	receipts := make(types.Receipts, 0)

	err := rlp.DecodeBytes(encodedReceipts, &receipts)
	if err != nil {
		s.logger.Crit("failed to decode receipts sent to enclave", log.ErrKey, err)
	}

	return receipts
}

func toRPCError(err common.SystemError) *generated.SystemError {
	if err == nil {
		return nil
	}
	return &generated.SystemError{
		ErrorCode:   1,
		ErrorString: err.Error(),
	}
}
