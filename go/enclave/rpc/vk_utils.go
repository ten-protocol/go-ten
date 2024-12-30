package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/vkhandler"
	"github.com/ten-protocol/go-ten/go/responses"
)

// ResourceStatus used as Status for the UserRPCRequests
type ResourceStatus int

var errInt = errors.New("internal error")

const (
	NotSet        ResourceStatus = iota // after initialisation
	Found                               // the parameters were parsed correctly and a From found
	NotAuthorised                       // not allowed to access the resource
	NotFound                            // resource not found
)

// CallBuilder - builder used during processing of an RPC request, which is a multi-step process
type CallBuilder[P any, R any] struct {
	ctx         context.Context
	Param       *P                                 // value calculated during phase 1 to be used during the execution phase
	VK          *vkhandler.AuthenticatedViewingKey // the vk accompanying the request
	From        *gethcommon.Address                // extracted from the request
	Status      ResourceStatus
	ReturnValue *R    // value to be returned to the user, encrypted
	Err         error // error to be returned to the user, encrypted
}

type (
	ValidateFunc[P any, R any] func([]any, *CallBuilder[P, R], *EncryptionManager) error
	ExecuteFunc[P any, R any]  func(*CallBuilder[P, R], *EncryptionManager) error
)

// HandleEncryptedRPC - handles the decryption, VK, and encryption
// validate - extract and validate the arguments
// execute - execute the user call only after authorising. Make sure to return a default value that makes sense in case of NotAuthorised
// note - authorisation is specific to each call
// e.g. - "getTransaction" or "getBalance" have to perform authorisation
// "Ten_call" , "Estimate_Gas" - have to authenticate the "From" - which will be used by the EVM
func HandleEncryptedRPC(ctx context.Context,
	encManager *EncryptionManager,
	encReq []byte, // encrypted request that contains a signed viewing key
) (*responses.EnclaveResponse, common.SystemError) {
	// 1. Decrypt request
	plaintextRequest, err := encManager.DecryptBytes(encReq)
	if err != nil {
		return responses.AsPlaintextError(common.FailedDecryptErr), nil
	}

	// 2. Unmarshall
	var decodedRequest rpc.RequestWithVk
	if err := json.Unmarshal(plaintextRequest, &decodedRequest); err != nil {
		return responses.AsPlaintextError(fmt.Errorf("could not unmarshal params - %w", err)), nil
	}

	// 3. Verify the VK
	if decodedRequest.VK == nil {
		return responses.AsPlaintextError(fmt.Errorf("invalid request. viewing key is missing")), nil
	}
	vk, err := vkhandler.VerifyViewingKey(decodedRequest.VK, encManager.config.ObscuroChainID)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("invalid viewing key - %w", err)), nil
	}

	// 4. Call the function that knows how to validate the request
	switch decodedRequest.Method {
	case rpc.ERPCCall:
		return withVKEncryption(ctx, encManager, decodedRequest, vk, TenCallValidate, TenCallExecute)
	case rpc.ERPCGetBalance:
		return withVKEncryption(ctx, encManager, decodedRequest, vk, GetBalanceValidate, GetBalanceExecute)
	case rpc.ERPCGetTransactionByHash:
		return withVKEncryption(ctx, encManager, decodedRequest, vk, GetTransactionValidate, GetTransactionExecute)
	case rpc.ERPCGetRawTransactionByHash:
		// todo - implement?
		return withVKEncryption(ctx, encManager, decodedRequest, vk, GetTransactionValidate, GetTransactionExecute)
	case rpc.ERPCGetTransactionCount:
		return withVKEncryption(ctx, encManager, decodedRequest, vk, GetTransactionCountValidate, GetTransactionCountExecute)
	case rpc.ERPCGetTransactionReceipt:
		return withVKEncryption(ctx, encManager, decodedRequest, vk, GetTransactionReceiptValidate, GetTransactionReceiptExecute)
	case rpc.ERPCSendRawTransaction:
		return withVKEncryption(ctx, encManager, decodedRequest, vk, SubmitTxValidate, SubmitTxExecute)
	case rpc.ERPCResend:
		// todo - implement
		return withVKEncryption(ctx, encManager, decodedRequest, vk, SubmitTxValidate, SubmitTxExecute)
	case rpc.ERPCEstimateGas:
		return withVKEncryption(ctx, encManager, decodedRequest, vk, EstimateGasValidate, EstimateGasExecute)
	case rpc.ERPCGetLogs:
		return withVKEncryption(ctx, encManager, decodedRequest, vk, GetLogsValidate, GetLogsExecute)
	case rpc.ERPCGetStorageAt:
		return withVKEncryption(ctx, encManager, decodedRequest, vk, TenStorageReadValidate, TenStorageReadExecute)
	case rpc.ERPCDebugLogs:
		return withVKEncryption(ctx, encManager, decodedRequest, vk, DebugLogsValidate, DebugLogsExecute)
	case rpc.ERPCGetPersonalTransactions:
		return withVKEncryption(ctx, encManager, decodedRequest, vk, GetPersonalTransactionsValidate, GetPersonalTransactionsExecute)
	default:
		panic(fmt.Sprintf("unsupported method %s", decodedRequest.Method))
	}
}

// withVKEncryption
// P - the type of the temporary parameter calculated after phase 1
// R - the type of the result
func withVKEncryption[P any, R any](
	ctx context.Context,
	encManager *EncryptionManager,
	decodedRequest rpc.RequestWithVk,
	vk *vkhandler.AuthenticatedViewingKey,
	validate ValidateFunc[P, R],
	execute ExecuteFunc[P, R],
) (*responses.EnclaveResponse, common.SystemError) {
	// 4. Call the function that knows how to validate the request
	builder := &CallBuilder[P, R]{Status: NotSet, VK: vk, ctx: ctx}

	err := validate(decodedRequest.Params, builder, encManager)
	if err != nil {
		return responses.AsPlaintextError(errInt), responses.ToInternalError(err)
	}
	if builder.Err != nil {
		return responses.AsEncryptedError(builder.Err, vk), nil //nolint:nilerr
	}

	// 5. Execute the authorisation and call
	// Note - it is the responsibility of this function to check that the authenticated address is authorised to view the data
	err = execute(builder, encManager)
	if err != nil {
		return responses.AsPlaintextError(errInt), responses.ToInternalError(err)
	}
	if builder.Err != nil {
		return responses.AsEncryptedError(builder.Err, vk), nil //nolint:nilerr
	}
	if builder.Status == NotFound {
		// if the requested resource was not found, return an empty response
		// todo - this must be encrypted - but we have some logic that expects it unencrypted, which is a bug
		// return responses.AsEncryptedEmptyResponse(vk), nil
		return responses.AsEmptyResponse(), nil
	}
	if builder.Status == NotAuthorised {
		// if the requested resource was not found, return an empty response
		return responses.AsEncryptedError(errors.New("not authorised"), vk), nil
	}

	return responses.AsEncryptedResponse[R](builder.ReturnValue, vk), nil
}

func authenticateFrom(vk *vkhandler.AuthenticatedViewingKey, from *gethcommon.Address) error {
	if from == nil || from.Hex() != vk.AccountAddress.Hex() {
		return fmt.Errorf("failed authentication. Account: %s does not match the from: %s", vk.AccountAddress, from)
	}
	return nil
}
