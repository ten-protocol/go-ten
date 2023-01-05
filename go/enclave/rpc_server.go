package enclave

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/rpc"
	"github.com/obscuronet/go-obscuro/go/common/rpc/generated"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"
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
	errStr := ""
	status, err := s.enclave.Status()
	if err != nil {
		errStr = err.Error()
	}
	return &generated.StatusResponse{Status: int32(status), Error: errStr}, nil
}

func (s *RPCServer) Attestation(context.Context, *generated.AttestationRequest) (*generated.AttestationResponse, error) {
	attestation, err := s.enclave.Attestation()
	if err != nil {
		return nil, err
	}
	msg := rpc.ToAttestationReportMsg(attestation)
	return &generated.AttestationResponse{AttestationReportMsg: &msg}, nil
}

func (s *RPCServer) GenerateSecret(context.Context, *generated.GenerateSecretRequest) (*generated.GenerateSecretResponse, error) {
	secret, err := s.enclave.GenerateSecret()
	if err != nil {
		return nil, err
	}
	return &generated.GenerateSecretResponse{EncryptedSharedEnclaveSecret: secret}, nil
}

func (s *RPCServer) InitEnclave(_ context.Context, request *generated.InitEnclaveRequest) (*generated.InitEnclaveResponse, error) {
	errStr := ""
	if err := s.enclave.InitEnclave(request.EncryptedSharedEnclaveSecret); err != nil {
		errStr = err.Error()
	}
	return &generated.InitEnclaveResponse{Error: errStr}, nil
}

func (s *RPCServer) SubmitL1Block(_ context.Context, request *generated.SubmitBlockRequest) (*generated.SubmitBlockResponse, error) {
	bl := s.decodeBlock(request.EncodedBlock)
	receipts := s.decodeReceipts(request.EncodedReceipts)
	blockSubmissionResponse, err := s.enclave.SubmitL1Block(bl, receipts, request.IsLatest)
	if err != nil {
		var rejErr *common.BlockRejectError
		isReject := errors.As(err, &rejErr)
		if isReject {
			// todo: we should avoid errors in response messages and use the gRPC error objects for this stuff (standardized across all enclave responses)
			msg, err := rpc.ToBlockSubmissionRejectionMsg(rejErr)
			if err == nil {
				// send back reject err response
				return &generated.SubmitBlockResponse{BlockSubmissionResponse: &msg}, nil
			}
			s.logger.Warn("failed to process the BlockRejectError, falling back to original error")
		}
		return nil, err
	}

	msg, err := rpc.ToBlockSubmissionResponseMsg(blockSubmissionResponse)
	if err != nil {
		return nil, err
	}
	return &generated.SubmitBlockResponse{BlockSubmissionResponse: &msg}, nil
}

func (s *RPCServer) SubmitTx(_ context.Context, request *generated.SubmitTxRequest) (*generated.SubmitTxResponse, error) {
	encryptedHash, err := s.enclave.SubmitTx(request.EncryptedTx)
	return &generated.SubmitTxResponse{EncryptedHash: encryptedHash}, err
}

func (s *RPCServer) SubmitBatch(_ context.Context, request *generated.SubmitBatchRequest) (*generated.SubmitBatchResponse, error) {
	batch := rpc.FromExtBatchMsg(request.Batch)
	return &generated.SubmitBatchResponse{}, s.enclave.SubmitBatch(batch)
}

func (s *RPCServer) ExecuteOffChainTransaction(_ context.Context, request *generated.OffChainRequest) (*generated.OffChainResponse, error) {
	result, err := s.enclave.ExecuteOffChainTransaction(request.EncryptedParams)
	if err != nil {
		// handle complex errors from the EVM
		errResponse, processErr := serializeEVMError(err)
		if processErr != nil {
			// unable to serialize the error
			return nil, fmt.Errorf("unable to serialise the EVM error - %w", processErr)
		}
		return &generated.OffChainResponse{Error: errResponse}, nil
	}
	return &generated.OffChainResponse{Result: result}, nil
}

func (s *RPCServer) GetTransactionCount(_ context.Context, request *generated.GetTransactionCountRequest) (*generated.GetTransactionCountResponse, error) {
	result, err := s.enclave.GetTransactionCount(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetTransactionCountResponse{Result: result}, nil
}

func (s *RPCServer) Stop(context.Context, *generated.StopRequest) (*generated.StopResponse, error) {
	defer s.grpcServer.GracefulStop()
	err := s.enclave.Stop()
	return &generated.StopResponse{}, err
}

func (s *RPCServer) GetTransaction(_ context.Context, request *generated.GetTransactionRequest) (*generated.GetTransactionResponse, error) {
	encryptedTx, err := s.enclave.GetTransaction(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetTransactionResponse{EncryptedTx: encryptedTx}, nil
}

func (s *RPCServer) GetTransactionReceipt(_ context.Context, request *generated.GetTransactionReceiptRequest) (*generated.GetTransactionReceiptResponse, error) {
	encryptedTxReceipt, err := s.enclave.GetTransactionReceipt(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetTransactionReceiptResponse{EncryptedTxReceipt: encryptedTxReceipt}, nil
}

func (s *RPCServer) AddViewingKey(_ context.Context, request *generated.AddViewingKeyRequest) (*generated.AddViewingKeyResponse, error) {
	err := s.enclave.AddViewingKey(request.ViewingKey, request.Signature)
	if err != nil {
		return nil, err
	}
	return &generated.AddViewingKeyResponse{}, nil
}

func (s *RPCServer) GetBalance(_ context.Context, request *generated.GetBalanceRequest) (*generated.GetBalanceResponse, error) {
	encryptedBalance, err := s.enclave.GetBalance(request.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetBalanceResponse{EncryptedBalance: encryptedBalance}, nil
}

func (s *RPCServer) GetCode(_ context.Context, request *generated.GetCodeRequest) (*generated.GetCodeResponse, error) {
	address := gethcommon.BytesToAddress(request.Address)
	rollupHash := gethcommon.BytesToHash(request.RollupHash)

	code, err := s.enclave.GetCode(address, &rollupHash)
	if err != nil {
		return nil, err
	}
	return &generated.GetCodeResponse{Code: code}, nil
}

func (s *RPCServer) Subscribe(_ context.Context, req *generated.SubscribeRequest) (*generated.SubscribeResponse, error) {
	err := s.enclave.Subscribe(gethrpc.ID(req.Id), req.EncryptedSubscription)
	return &generated.SubscribeResponse{}, err
}

func (s *RPCServer) Unsubscribe(_ context.Context, req *generated.UnsubscribeRequest) (*generated.UnsubscribeResponse, error) {
	err := s.enclave.Unsubscribe(gethrpc.ID(req.Id))
	return &generated.UnsubscribeResponse{}, err
}

func (s *RPCServer) EstimateGas(_ context.Context, req *generated.EstimateGasRequest) (*generated.EstimateGasResponse, error) {
	encryptedBalance, err := s.enclave.EstimateGas(req.EncryptedParams)
	if err != nil {
		// handle complex errors from the EVM
		_, processErr := serializeEVMError(err)
		if processErr != nil {
			// unable to serialize the error
			return nil, fmt.Errorf("unable to serialise the EVM error - %w", processErr)
		}
		return nil, err
	}
	return &generated.EstimateGasResponse{EncryptedResponse: encryptedBalance}, nil
}

func (s *RPCServer) GetLogs(_ context.Context, req *generated.GetLogsRequest) (*generated.GetLogsResponse, error) {
	encryptedLogs, err := s.enclave.GetLogs(req.EncryptedParams)
	if err != nil {
		return nil, err
	}
	return &generated.GetLogsResponse{EncryptedResponse: encryptedLogs}, nil
}

func (s *RPCServer) HealthCheck(_ context.Context, _ *generated.EmptyArgs) (*generated.HealthCheckResponse, error) {
	healthy, err := s.enclave.HealthCheck()
	if err != nil {
		return nil, err
	}
	return &generated.HealthCheckResponse{Status: healthy}, nil
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

// serializeEVMError serialises EVM errors into the RPC response
// always returns a SerialisableError byte slice
func serializeEVMError(err error) ([]byte, error) {
	var errReturn interface{}

	// check if it's a serialized error and handle any error wrapping that might have occurred
	var e *evm.SerialisableError
	if ok := errors.As(err, &e); ok {
		errReturn = e
	} else {
		// it's a generic error, serialise it
		errReturn = &evm.SerialisableError{Err: err.Error()}
	}

	// serialise the error object returned by the evm into a json
	errSerializedBytes, marshallErr := json.Marshal(errReturn)
	if marshallErr != nil {
		return nil, marshallErr
	}
	return errSerializedBytes, nil
}
