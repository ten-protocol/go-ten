package responses

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ten-protocol/go-ten/go/common/syserr"
)

// InternalErrMsg is the common response returned to the user when an InternalError occurs
var InternalErrMsg = "internal system error"

// EncryptedUserResponse - This is the encoded & encrypted form of a UserResponse[Type]
type EncryptedUserResponse []byte

// EnclaveResponse - The response that the enclave returns for sensitive API calls
// The user response is encrypted while the error is in plaintext
type EnclaveResponse struct {
	EncUserResponse EncryptedUserResponse
	Err             *string
}

// Encode - serializes the enclave response into a json
func (er *EnclaveResponse) Encode() []byte {
	encoded, err := json.Marshal(er)
	if err != nil {
		panic(err)
	}

	return encoded
}

func (er *EnclaveResponse) Error() error {
	if er.Err != nil {
		return fmt.Errorf(*er.Err)
	}
	return nil
}

// AsPlaintextResponse - creates the plaintext part of the enclave response
// It would be visible that there is an enclave response,
// but the bytes in it will still be encrypted
// todo - rename
func AsPlaintextResponse(encResp EncryptedUserResponse) *EnclaveResponse {
	return &EnclaveResponse{
		EncUserResponse: encResp,
	}
}

// AsEmptyResponse - Creates an empty enclave response. Useful for when no error
// encountered but also no result found.
func AsEmptyResponse() *EnclaveResponse {
	return &EnclaveResponse{
		EncUserResponse: nil,
		Err:             nil,
	}
}

// AsSystemErr - generates a plaintext response containing a visible error.
func AsSystemErr() *EnclaveResponse {
	return &EnclaveResponse{
		Err: &InternalErrMsg,
	}
}

// AsPlaintextError - generates a plaintext response containing a visible to the host error.
func AsPlaintextError(err error) *EnclaveResponse {
	errStr := err.Error()
	return &EnclaveResponse{
		Err: &errStr,
	}
}

type Encryptor interface {
	Encrypt(bytes []byte) ([]byte, error)
}

// AsEncryptedResponse - wraps the data passed into the proper format, serializes it and encrypts it.
// It is then encoded in a plaintext response.
func AsEncryptedResponse[T any](data *T, encryptHandler Encryptor) *EnclaveResponse {
	userResp := UserResponse[T]{
		Result: data,
	}

	encoded, err := json.Marshal(userResp)
	if err != nil {
		return AsPlaintextError(err)
	}

	encrypted, err := encryptHandler.Encrypt(encoded)
	if err != nil {
		return AsPlaintextError(err)
	}

	return AsPlaintextResponse(encrypted)
}

// AsEncryptedError - Encodes and encrypts an error to be returned for a concrete user.
func AsEncryptedError(err error, encrypt Encryptor) *EnclaveResponse {
	userResp := UserResponse[string]{
		Err: convertError(err),
	}

	encoded, err := json.Marshal(userResp)
	if err != nil {
		return AsPlaintextError(err)
	}

	encrypted, err := encrypt.Encrypt(encoded)
	if err != nil {
		return AsPlaintextError(err)
	}

	return AsPlaintextResponse(encrypted)
}

// ToEnclaveResponse - Converts an encoded plaintext into an enclave response
func ToEnclaveResponse(encoded []byte) *EnclaveResponse {
	resp := EnclaveResponse{}
	err := json.Unmarshal(encoded, &resp)
	if err != nil {
		panic(err) // Todo change when stable.
	}
	return &resp
}

// ToInternalError - Converts an error to an InternalError
func ToInternalError(err error) error {
	if err == nil {
		return nil
	}

	return syserr.NewInternalError(err)
}

// DecodeResponse - Extracts the user response from a decrypted bytes field and returns the
// result or nil and optional error.
func DecodeResponse[T any](encoded []byte) (*T, error) {
	resp := UserResponse[T]{}
	err := json.Unmarshal(encoded, &resp)
	if err != nil {
		return nil, fmt.Errorf("could not decode response. Cause: %w", err)
	}
	if resp.Err != nil {
		return nil, resp.Err
	}

	return resp.Result, nil
}

func convertError(err error) *errutil.DataError {
	// check if it's a serialized error and handle any error wrapping that might have occurred
	var e *errutil.DataError
	if ok := errors.As(err, &e); ok {
		return e
	}
	return &errutil.DataError{Err: err.Error()}
}
