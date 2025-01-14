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
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/rpc"
	"github.com/ten-protocol/go-ten/go/common/rpc/generated"
	"github.com/ten-protocol/go-ten/go/common/tracers"
	"google.golang.org/grpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
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
		enclave: enclave,
		grpcServer: grpc.NewServer(
			grpc.MaxRecvMsgSize(1024 * 1024 * 50),
		),
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
func (s *RPCServer) Status(ctx context.Context, _ *generated.StatusRequest) (*generated.StatusResponse, error) {
	status, sysError := s.enclave.Status(ctx)
	if sysError != nil {
		s.logger.Error("Enclave error on Status", log.ErrKey, sysError)
	}
	var l2Head []byte
	if status.L2Head != nil {
		l2Head = status.L2Head.Bytes()
	}
	return &generated.StatusResponse{
		StatusCode:  int32(status.StatusCode),
		L1Head:      status.L1Head.Bytes(),
		L2Head:      l2Head,
		EnclaveID:   status.EnclaveID.Bytes(),
		SystemError: toRPCError(sysError),
	}, nil
}

func (s *RPCServer) AddSequencer(ctx context.Context, in *generated.AddSequencerRequest) (*generated.AddSequencerResponse, error) {
	rec := types.Receipt{}
	err := rec.UnmarshalBinary(in.Proof)
	if err != nil {
		return &generated.AddSequencerResponse{SystemError: toRPCError(err)}, nil
	}
	return &generated.AddSequencerResponse{SystemError: toRPCError(s.enclave.AddSequencer(gethcommon.BytesToAddress(in.EnclaveId), rec))}, nil
}

func (s *RPCServer) MakeActive(ctx context.Context, in *generated.MakeActiveRequest) (*generated.MakeActiveResponse, error) {
	return &generated.MakeActiveResponse{SystemError: toRPCError(s.enclave.MakeActive())}, nil
}

func (s *RPCServer) Attestation(ctx context.Context, _ *generated.AttestationRequest) (*generated.AttestationResponse, error) {
	attestation, sysError := s.enclave.Attestation(ctx)
	if sysError != nil {
		s.logger.Error("Error getting attestation", log.ErrKey, sysError)
		return &generated.AttestationResponse{SystemError: toRPCError(sysError)}, nil
	}
	msg := rpc.ToAttestationReportMsg(attestation)
	return &generated.AttestationResponse{AttestationReportMsg: &msg}, nil
}

func (s *RPCServer) GenerateSecret(ctx context.Context, _ *generated.GenerateSecretRequest) (*generated.GenerateSecretResponse, error) {
	secret, sysError := s.enclave.GenerateSecret(ctx)
	if sysError != nil {
		s.logger.Error("Error generating secret", log.ErrKey, sysError)
		return &generated.GenerateSecretResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.GenerateSecretResponse{EncryptedSharedEnclaveSecret: secret}, nil
}

func (s *RPCServer) InitEnclave(ctx context.Context, request *generated.InitEnclaveRequest) (*generated.InitEnclaveResponse, error) {
	sysError := s.enclave.InitEnclave(ctx, request.EncryptedSharedEnclaveSecret)
	if sysError != nil {
		s.logger.Error("Error initialising the enclave", log.ErrKey, sysError)
	}
	return &generated.InitEnclaveResponse{SystemError: toRPCError(sysError)}, nil
}

func (s *RPCServer) EnclaveID(ctx context.Context, _ *generated.EnclaveIDRequest) (*generated.EnclaveIDResponse, error) {
	id, sysError := s.enclave.EnclaveID(ctx)
	if sysError != nil {
		s.logger.Error("Error getting enclave ID", log.ErrKey, sysError)
		return &generated.EnclaveIDResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.EnclaveIDResponse{EnclaveID: id.Bytes()}, nil
}

func (s *RPCServer) RPCEncryptionKey(ctx context.Context, _ *generated.RPCEncryptionKeyRequest) (*generated.RPCEncryptionKeyResponse, error) {
	key, sysError := s.enclave.RPCEncryptionKey(ctx)
	if sysError != nil {
		s.logger.Error("Error getting enclave ID", log.ErrKey, sysError)
		return &generated.RPCEncryptionKeyResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.RPCEncryptionKeyResponse{RpcPubKey: key}, nil
}

func (s *RPCServer) SubmitL1Block(ctx context.Context, request *generated.SubmitBlockRequest) (*generated.SubmitBlockResponse, error) {
	processedData, err := s.decodeProcessedData(request.EncodedProcessedData)
	if err != nil {
		s.logger.Error("Error decoding receipts", log.ErrKey, err)
		return nil, err
	}

	blockSubmissionResponse, err := s.enclave.SubmitL1Block(ctx, processedData)
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
		s.logger.Error("Unexpected error submitting the L1 block", log.ErrKey, err)
		return nil, err
	}

	msg, err := rpc.ToBlockSubmissionResponseMsg(blockSubmissionResponse)
	if err != nil {
		return nil, err
	}
	return &generated.SubmitBlockResponse{BlockSubmissionResponse: msg}, nil
}

func (s *RPCServer) SubmitBatch(ctx context.Context, request *generated.SubmitBatchRequest) (*generated.SubmitBatchResponse, error) {
	batch := rpc.FromExtBatchMsg(request.Batch)
	sysError := s.enclave.SubmitBatch(ctx, batch)
	if sysError != nil {
		s.logger.Error("Error submitting batch", log.ErrKey, sysError)
	}
	return &generated.SubmitBatchResponse{SystemError: toRPCError(sysError)}, nil
}

func (s *RPCServer) Stop(context.Context, *generated.StopRequest) (*generated.StopResponse, error) {
	// stop the grpcServer on its own goroutine to avoid killing the existing connection
	go s.grpcServer.GracefulStop()
	return &generated.StopResponse{SystemError: toRPCError(s.enclave.Stop())}, nil
}

func (s *RPCServer) GetCode(ctx context.Context, request *generated.GetCodeRequest) (*generated.GetCodeResponse, error) {
	address := gethcommon.BytesToAddress(request.Address)

	blockNrOrHash := &gethrpc.BlockNumberOrHash{}
	err := blockNrOrHash.UnmarshalJSON(request.BlockNrOrHash)
	if err != nil {
		s.logger.Error("Error unmarshalling block nr or hash", log.ErrKey, err)
		return &generated.GetCodeResponse{SystemError: toRPCError(err)}, nil
	}

	code, sysError := s.enclave.GetCode(ctx, address, *blockNrOrHash)
	if sysError != nil {
		s.logger.Error("Error getting code", log.ErrKey, sysError)
		return &generated.GetCodeResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.GetCodeResponse{Code: code}, nil
}

func (s *RPCServer) Subscribe(ctx context.Context, req *generated.SubscribeRequest) (*generated.SubscribeResponse, error) {
	sysError := s.enclave.Subscribe(ctx, gethrpc.ID(req.Id), req.EncryptedSubscription)
	if sysError != nil {
		s.logger.Error("Error subscribing", log.ErrKey, sysError)
	}
	return &generated.SubscribeResponse{SystemError: toRPCError(sysError)}, nil
}

func (s *RPCServer) Unsubscribe(_ context.Context, req *generated.UnsubscribeRequest) (*generated.UnsubscribeResponse, error) {
	sysError := s.enclave.Unsubscribe(gethrpc.ID(req.Id))
	if sysError != nil {
		s.logger.Error("Error unsubscribing", log.ErrKey, sysError)
	}
	return &generated.UnsubscribeResponse{SystemError: toRPCError(sysError)}, nil
}

func (s *RPCServer) HealthCheck(ctx context.Context, _ *generated.EmptyArgs) (*generated.HealthCheckResponse, error) {
	healthy, sysError := s.enclave.HealthCheck(ctx)
	if sysError != nil {
		return &generated.HealthCheckResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.HealthCheckResponse{Status: healthy}, nil
}

func (s *RPCServer) CreateRollup(ctx context.Context, req *generated.CreateRollupRequest) (*generated.CreateRollupResponse, error) {
	var fromSeqNo uint64 = 1
	if req.FromSequenceNumber != nil && *req.FromSequenceNumber > common.L2GenesisSeqNo {
		fromSeqNo = *req.FromSequenceNumber
	}

	rollup, metadata, sysError := s.enclave.CreateRollup(ctx, fromSeqNo)
	if sysError != nil {
		s.logger.Error("Error creating rollup", log.ErrKey, sysError)
	}

	msg := rpc.ToExtRollupMsg(rollup)

	genMetadata := &generated.ExtRollupMetadataResponseMsg{
		CrossChainTree: []byte{},
	}
	if metadata != nil {
		genMetadata.CrossChainTree = metadata.CrossChainTree
	}

	return &generated.CreateRollupResponse{
		Msg:         &msg,
		Metadata:    genMetadata,
		SystemError: toRPCError(sysError),
	}, nil
}

func (s *RPCServer) CreateBatch(ctx context.Context, r *generated.CreateBatchRequest) (*generated.CreateBatchResponse, error) {
	sysError := s.enclave.CreateBatch(ctx, r.SkipIfEmpty)
	if sysError != nil {
		s.logger.Error("Error creating batch", log.ErrKey, sysError)
	}
	return &generated.CreateBatchResponse{}, sysError
}

func (s *RPCServer) DebugTraceTransaction(ctx context.Context, req *generated.DebugTraceTransactionRequest) (*generated.DebugTraceTransactionResponse, error) {
	txHash := gethcommon.BytesToHash(req.TxHash)
	var config tracers.TraceConfig

	sysError := json.Unmarshal(req.Config, &config)
	if sysError != nil {
		s.logger.Error("Error calling debug tx", log.ErrKey, sysError)

		return &generated.DebugTraceTransactionResponse{
			SystemError: toRPCError(fmt.Errorf("unable to unmarshall config - %w", sysError)),
		}, nil
	}

	traceTx, sysError := s.enclave.DebugTraceTransaction(ctx, txHash, &config)
	return &generated.DebugTraceTransactionResponse{Msg: string(traceTx), SystemError: toRPCError(sysError)}, nil
}

func (s *RPCServer) GetBatch(ctx context.Context, request *generated.GetBatchRequest) (*generated.GetBatchResponse, error) {
	batch, err := s.enclave.GetBatch(ctx, gethcommon.BytesToHash(request.KnownHead))
	if err != nil {
		s.logger.Error("Error getting batch", log.ErrKey, err)
		// todo  do we want to exit here or return the usual response
		return nil, err
	}

	encodedBatch, encodingErr := batch.Encoded()
	var sysErr *generated.SystemError
	if encodingErr != nil {
		sysErr = &generated.SystemError{
			ErrorCode:   2,
			ErrorString: encodingErr.Error(),
		}
	}
	return &generated.GetBatchResponse{
		Batch:       encodedBatch,
		SystemError: sysErr,
	}, err
}

func (s *RPCServer) GetBatchBySeqNo(ctx context.Context, request *generated.GetBatchBySeqNoRequest) (*generated.GetBatchResponse, error) {
	batch, err := s.enclave.GetBatchBySeqNo(ctx, request.SeqNo)
	if err != nil {
		s.logger.Error("Error getting batch by seq", log.ErrKey, err)
		// todo  do we want to exit here or return the usual response
		return nil, err
	}

	encodedBatch, encodingErr := batch.Encoded()
	var sysErr *generated.SystemError
	if encodingErr != nil {
		sysErr = &generated.SystemError{
			ErrorCode:   2,
			ErrorString: encodingErr.Error(),
		}
	}
	return &generated.GetBatchResponse{
		Batch:       encodedBatch,
		SystemError: sysErr,
	}, err
}

func (s *RPCServer) ExportCrossChainData(ctx context.Context, request *generated.ExportCrossChainDataRequest) (*generated.ExportCrossChainDataResponse, error) {
	bundle, err := s.enclave.ExportCrossChainData(ctx, request.FromSeqNo, request.ToSeqNo)
	if err != nil {
		return nil, err
	}

	encodedBundle, err := rlp.EncodeToBytes(bundle)
	if err != nil {
		return nil, err
	}

	return &generated.ExportCrossChainDataResponse{
		Msg: encodedBundle,
	}, nil
}

func (s *RPCServer) GetRollupData(ctx context.Context, request *generated.GetRollupDataRequest) (*generated.GetRollupDataResponse, error) {
	rollupMetadata, sysError := s.enclave.GetRollupData(ctx, gethcommon.BytesToHash(request.Hash))
	if sysError != nil {
		s.logger.Error("Error fetching rollup metadata", log.ErrKey, sysError)
		return nil, sysError
	}

	msg := rpc.ToRollupDataMsg(rollupMetadata)

	return &generated.GetRollupDataResponse{
		Msg:         &msg,
		SystemError: toRPCError(sysError),
	}, nil
}

func (s *RPCServer) StreamL2Updates(_ *generated.StreamL2UpdatesRequest, stream generated.EnclaveProto_StreamL2UpdatesServer) error {
	batchChan, stop := s.enclave.StreamL2Updates()
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

func (s *RPCServer) GetTotalContractCount(ctx context.Context, _ *generated.GetTotalContractCountRequest) (*generated.GetTotalContractCountResponse, error) {
	count, sysError := s.enclave.GetTotalContractCount(ctx)
	if sysError != nil {
		s.logger.Error("Error GetTotalContractCount", log.ErrKey, sysError)
	}

	if count == nil {
		count = big.NewInt(0)
	}

	return &generated.GetTotalContractCountResponse{
		Count:       count.Int64(),
		SystemError: toRPCError(sysError),
	}, nil
}

func (s *RPCServer) EncryptedRPC(ctx context.Context, req *generated.EncCallRequest) (*generated.EncCallResponse, error) {
	enclaveResp, sysError := s.enclave.EncryptedRPC(ctx, req.EncryptedParams)
	if sysError != nil {
		s.logger.Error("Error getting receipt", log.ErrKey, sysError)
		return &generated.EncCallResponse{SystemError: toRPCError(sysError)}, nil
	}
	return &generated.EncCallResponse{EncodedEnclaveResponse: enclaveResp.Encode()}, nil
}

func (s *RPCServer) EnclavePublicConfig(ctx context.Context, _ *generated.EnclavePublicConfigRequest) (*generated.EnclavePublicConfigResponse, error) {
	enclaveCfg, sysError := s.enclave.EnclavePublicConfig(ctx)
	if sysError != nil {
		s.logger.Error("Error getting message bus address", log.ErrKey, sysError)
		return &generated.EnclavePublicConfigResponse{SystemError: toRPCError(sysError)}, nil
	}

	publicContracts := make(map[string][]byte, len(enclaveCfg.PublicSystemContracts))
	for k, v := range enclaveCfg.PublicSystemContracts {
		publicContracts[k] = v.Bytes()
	}

	return &generated.EnclavePublicConfigResponse{
		L2MessageBusAddress:             enclaveCfg.L2MessageBusAddress.Bytes(),
		TransactionPostProcessorAddress: enclaveCfg.TransactionPostProcessorAddress.Bytes(),
		PublicSystemContracts:           publicContracts,
	}, nil
}

// decodeProcessedData - converts the rlp encoded bytes to processed if possible.
func (s *RPCServer) decodeProcessedData(encodedData []byte) (*common.ProcessedL1Data, error) {
	var processed common.ProcessedL1Data

	err := rlp.DecodeBytes(encodedData, &processed)
	if err != nil {
		return nil, fmt.Errorf("unable to decode receipts, bytes=%x, err=%w", encodedData, err)
	}

	return &processed, nil
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
