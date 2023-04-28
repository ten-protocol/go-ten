package responses

import (
	"encoding/json"
	"fmt"
)

// InternalErrMsg is the common response returned to the user when an InternalError occurs
var InternalErrMsg = "internal system error"

// This is the encoded & encrypted form of a UserResponse[Type]
type EncryptedUserResponse []byte

// The response that the enclave returns for sensitive API calls
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

// AsEncryptedResponse - wraps the data passed into the proper format, serializes it and encrypts it.
// It is then encoded in a plaintext response.
func AsEncryptedResponse[T any](data *T, encrypt ViewingKeyEncryptor) *EnclaveResponse {
	userResp := UserResponse[T]{
		Result: data,
	}

	encoded, err := json.Marshal(userResp)
	if err != nil {
		return AsPlaintextError(err)
	}

	encrypted, err := encrypt(encoded)
	if err != nil {
		return AsPlaintextError(err)
	}

	return AsPlaintextResponse(encrypted)
}

// AsEncryptedError - Encodes and encrypts an error to be returned for a concrete user.
func AsEncryptedError(err error, encrypt ViewingKeyEncryptor) *EnclaveResponse {
	errStr := err.Error()
	userResp := UserResponse[string]{
		ErrStr: &errStr,
	}

	encoded, err := json.Marshal(userResp)
	if err != nil {
		return AsPlaintextError(err)
	}

	encrypted, err := encrypt(encoded)
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

// DecodeResponse - Extracts the user response from a decrypted bytes field and returns the
// result or nil and optional error.
func DecodeResponse[T any](encoded []byte) (*T, error) {
	resp := UserResponse[T]{}
	err := json.Unmarshal(encoded, &resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrStr != nil {
		return nil, fmt.Errorf(*resp.ErrStr)
	}

	return resp.Result, nil
}
