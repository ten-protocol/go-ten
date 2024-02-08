package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/vkhandler"
	"github.com/ten-protocol/go-ten/go/responses"
)

// ResourceStatus used as Status for the UserRPCRequests
type ResourceStatus int

const (
	NotSet        ResourceStatus = iota // after initialisation
	Found                               // the parameters were parsed correctly and a From found
	NotAuthorised                       // not allowed to access the resource
	NotFound                            // resource not found
)

// RpcCallBuilder1 - used during processing of an RPC request, which is a multi-step process
type RpcCallBuilder1[P any, R any] struct {
	Param         *P                                 // value calculated during phase 1 to be used during the execution phase
	VK            *vkhandler.AuthenticatedViewingKey // the vk accompanying the request
	From          *gethcommon.Address                // extracted from the request
	ResourceOwner *gethcommon.Address                // extracted from the database Not applicable for all requests. E.g. For a tx, the owner is the original tx sender
	Status        ResourceStatus                     //
	ReturnValue   *R                                 // the value to be returned to the user
	Err           error                              // encrypted error to be returned to the user
}

// RpcCallBuilder2 - same as RpcCallBuilder1 but with 2 typed intermediate arguments
type RpcCallBuilder2[P1 any, P2 any, R any] struct {
	Param1        *P1
	Param2        *P2
	VK            *vkhandler.AuthenticatedViewingKey
	From          *gethcommon.Address
	ResourceOwner *gethcommon.Address
	Status        ResourceStatus
	ReturnValue   *R
	Err           error
}

// WithVKEncryption1 - handles the VK management, authentication and encryption
// P represents the single request parameter
// R represents the response which will be encrypted
// note - this is a thin wrapper over WithVKEncryption2
func WithVKEncryption1[P any, R any](
	encManager *EncryptionManager,
	chainID int64,
	encReq []byte, // encrypted request that contains a signed viewing key
	extractFromAndParams func([]any, *RpcCallBuilder1[P, R], *EncryptionManager) error, // extract the arguments and the logical sender from the plaintext request. Make sure to not return any information from the db in the error.
	executeCall func(*RpcCallBuilder1[P, R], *EncryptionManager) error, // execute the user call against the authenticated request.
) (*responses.EnclaveResponse, common.SystemError) {
	return WithVKEncryption2[P, P, R](encManager,
		chainID,
		encReq,
		func(params []any, rpcBuilder *RpcCallBuilder2[P, P, R], em *EncryptionManager) error {
			temp := toOneParam(rpcBuilder)
			err := extractFromAndParams(params, temp, em)
			if err != nil {
				return err
			}
			copyBuilder(rpcBuilder, temp)
			return nil
		},
		func(rpcBuilder *RpcCallBuilder2[P, P, R], em *EncryptionManager) error {
			temp := toOneParam(rpcBuilder)
			err := executeCall(temp, em)
			if err != nil {
				return err
			}
			copyBuilder(rpcBuilder, temp)
			return nil
		})
}

// WithVKEncryption2 - similar to WithVKEncryption1, but supports two arguments
func WithVKEncryption2[P1 any, P2 any, R any](
	encManager *EncryptionManager,
	chainID int64,
	encReq []byte, // encrypted request that contains a signed viewing key
	extractFromAndParams func([]any, *RpcCallBuilder2[P1, P2, R], *EncryptionManager) error, // extract the arguments and the logical sender from the plaintext request. Make sure to not return any information from the db in the error.
	executeCall func(*RpcCallBuilder2[P1, P2, R], *EncryptionManager) error, // execute the user call. Returns a user error or a system error
) (*responses.EnclaveResponse, common.SystemError) {
	// 1. Decrypt request
	plaintextRequest, err := encManager.DecryptBytes(encReq)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("could not decrypt params - %w", err)), nil
	}

	// 2. Unmarshall into a generic []any array
	var decodedRequest []any
	if err := json.Unmarshal(plaintextRequest, &decodedRequest); err != nil {
		return responses.AsPlaintextError(fmt.Errorf("could not unmarshal params - %w", err)), nil
	}

	// 3. Extract the VK from the first element and verify it
	if len(decodedRequest) < 1 {
		return responses.AsPlaintextError(fmt.Errorf("invalid request. viewing key is missing")), nil
	}
	rpcVK, ok := decodedRequest[0].(viewingkey.RPCSignedViewingKey)
	if !ok {
		return responses.AsPlaintextError(fmt.Errorf("invalid request. viewing key encoded incorrectly")), nil
	}
	vk, err := vkhandler.VerifyViewingKey(rpcVK, chainID)
	if err != nil {
		return responses.AsPlaintextError(fmt.Errorf("invalid viewing key - %w", err)), nil
	}

	// 4. Call the function that knows how to validate the request
	builder := &RpcCallBuilder2[P1, P2, R]{Status: NotSet, VK: vk}

	err = extractFromAndParams(decodedRequest[1:], builder, encManager)
	if err != nil {
		return nil, responses.ToInternalError(err)
	}
	if builder == nil {
		return nil, responses.ToInternalError(fmt.Errorf("should not happen"))
	}
	if builder.Err != nil {
		return responses.AsEncryptedError(fmt.Errorf("invalid request - %w", builder.Err), vk), nil
	}

	// 5. IMPORTANT!: authenticate the call.
	if builder.From != nil && builder.From.Hex() != vk.AccountAddress.Hex() {
		return responses.AsEncryptedError(fmt.Errorf("failed authentication. Account: %s does not match the from: %s", vk.AccountAddress, builder.From), vk), nil
	}

	// 6. Make the backend call and convert the response.
	err = executeCall(builder, encManager)
	if err != nil {
		return nil, responses.ToInternalError(err)
	}
	if builder.Err != nil {
		return responses.AsEncryptedError(fmt.Errorf("invalid request - %w", builder.Err), vk), nil
	}
	if builder.Status == NotFound || builder.Status == NotAuthorised {
		// if the requested resource was not found, return an empty response
		// todo - this must be encrypted
		// return responses.AsEncryptedEmptyResponse(vk), nil
		return responses.AsEmptyResponse(), nil
	}

	// double check authorisation
	if builder.ResourceOwner != nil && builder.ResourceOwner.Hex() != vk.AccountAddress.Hex() {
		return responses.AsEmptyResponse(), nil
	}

	return responses.AsEncryptedResponse[R](builder.ReturnValue, vk), nil
}

func copyBuilder[P any, R any](rpcBuilder *RpcCallBuilder2[P, P, R], temp *RpcCallBuilder1[P, R]) {
	rpcBuilder.Param1 = temp.Param
	rpcBuilder.Param2 = nil
	rpcBuilder.From = temp.From
	rpcBuilder.ReturnValue = temp.ReturnValue
	rpcBuilder.ResourceOwner = temp.ResourceOwner
	rpcBuilder.VK = temp.VK
	rpcBuilder.Status = temp.Status
	rpcBuilder.Err = temp.Err
}

func toOneParam[P any, R any](rpcBuilder *RpcCallBuilder2[P, P, R]) *RpcCallBuilder1[P, R] {
	return &RpcCallBuilder1[P, R]{
		Param:         rpcBuilder.Param1,
		From:          rpcBuilder.From,
		ReturnValue:   rpcBuilder.ReturnValue,
		ResourceOwner: rpcBuilder.ResourceOwner,
		VK:            rpcBuilder.VK,
		Status:        rpcBuilder.Status,
		Err:           rpcBuilder.Err,
	}
}
