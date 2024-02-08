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

// UserRPCRequest1 - decoded RPC argument accompanied by a logical sender
type UserRPCRequest1[P any] struct {
	Sender *gethcommon.Address
	Param1 *P
}

// UserRPCRequest2 - decoded RPC arguments accompanied by a logical sender (2 arguments)
type UserRPCRequest2[P1 any, P2 any] struct {
	Sender *gethcommon.Address
	Param1 *P1
	Param2 *P2
}

// UserResponse - the result of executing the Request against the services. Paired with a validation error that must be returned to the user.
type UserResponse[R any] struct {
	Val *R
	Err error // the error will be encrypted
}

// WithVKEncryption1- handles the VK management, authentication and encryption
// P represents the single request parameter
// R represents the response which will be encrypted
func WithVKEncryption1[P any, R any](
	encManager *EncryptionManager,
	chainID int64,
	encReq []byte, // encrypted request that contains a signed viewing key
	extractFromAndParams func([]any, *EncryptionManager) (*UserRPCRequest1[P], error), // extract the arguments and the logical sender from the plaintext request. Make sure to not return any information from the db in the error.
	executeCall func(*UserRPCRequest1[P], *EncryptionManager) (*UserResponse[R], error), // execute the user call against the authenticated request.
) (*responses.EnclaveResponse, common.SystemError) {
	return WithVKEncryption2[P, P, R](encManager,
		chainID,
		encReq,
		func(params []any, em *EncryptionManager) (*UserRPCRequest2[P, P], error) {
			res, err := extractFromAndParams(params, em)
			if err != nil {
				return nil, err
			}
			if res == nil {
				return nil, nil
			}
			return &UserRPCRequest2[P, P]{res.Sender, res.Param1, nil}, nil
		},
		func(req *UserRPCRequest2[P, P], em *EncryptionManager) (*UserResponse[R], error) {
			return executeCall(&UserRPCRequest1[P]{req.Sender, req.Param1}, em)
		})
}

// WithVKEncryption2 - similar to WithVKEncryption1, but supports two arguments
func WithVKEncryption2[P1 any, P2 any, R any](
	encManager *EncryptionManager,
	chainID int64,
	encReq []byte, // encrypted request that contains a signed viewing key
	extractFromAndParams func([]any, *EncryptionManager) (*UserRPCRequest2[P1, P2], error), // extract the arguments and the logical sender from the plaintext request. Make sure to not return any information from the db in the error.
	executeCall func(*UserRPCRequest2[P1, P2], *EncryptionManager) (*UserResponse[R], error), // execute the user call. Returns a user error or a system error
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

	// 3. Extract the VK from the first element
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

	// 4. Call the function that knows how to extract request specific params from the request
	decodedParams, err := extractFromAndParams(decodedRequest[1:], encManager)
	if err != nil {
		return responses.AsEncryptedError(fmt.Errorf("unable to decode params - %w", err), vk), nil
	}

	// when all return values are null, by convention this is "Not found", so we just return an empty value
	if decodedParams == nil && err == nil {
		// todo - this must be encrypted
		// return responses.AsEncryptedEmptyResponse(vk), nil
		return responses.AsEmptyResponse(), nil
	}

	// 5. Validate the logical sender
	if decodedParams.Sender == nil {
		return responses.AsEncryptedError(fmt.Errorf("invalid request - `from` field is mandatory"), vk), nil
	}

	// IMPORTANT!: this is where we authenticate the call.
	if decodedParams.Sender.Hex() != vk.AccountAddress.Hex() {
		return responses.AsEncryptedError(fmt.Errorf("failed authentication: account: %s does not match the requester: %s", vk.AccountAddress, decodedParams.Sender), vk), nil
	}

	// 6. Make the backend call and convert the response.
	response, sysErr := executeCall(decodedParams, encManager)
	if sysErr != nil {
		return nil, responses.ToInternalError(sysErr)
	}
	if response.Err != nil {
		return responses.AsEncryptedError(response.Err, vk), nil //nolint:nilerr
	}
	if response.Val == nil {
		return responses.AsEncryptedEmptyResponse(vk), nil
	}
	return responses.AsEncryptedResponse[R](response.Val, vk), nil
}
