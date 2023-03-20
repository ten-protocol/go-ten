package responses

import (
	"encoding/json"
	"fmt"
)

// This is the encoded & encrypted form of a UserResponse[Type]
type EncryptedUserResponse []byte

// The response that the enclave returns for sensitive API calls
// The user response is encrypted while the error is in plaintext
type EnclaveResponse struct {
	EncUserResponse EncryptedUserResponse
	Err             error
}

// ToEnclaveResponse - Converts an encoded plaintext into an enclave response
func ToEnclaveResponse(encoded []byte) *EnclaveResponse {
	resp := EnclaveResponse{}
	transportStruct := struct {
		Resp EncryptedUserResponse
		Err  *string
	}{}

	err := json.Unmarshal(encoded, &transportStruct)
	if err != nil {
		panic(err) // Todo change when stable.
	}
	resp.EncUserResponse = transportStruct.Resp
	if transportStruct.Err != nil {
		resp.Err = fmt.Errorf(*transportStruct.Err)
	}

	return &resp
}

// Encode - serializes the enclave response into a json
func (er *EnclaveResponse) Encode() []byte {
	transportStruct := struct {
		Resp EncryptedUserResponse
		Err  *string
	}{
		Resp: er.EncUserResponse,
	}

	if er.Err != nil {
		errStr := er.Err.Error()
		transportStruct.Err = &errStr
	}

	encoded, err := json.Marshal(transportStruct)
	if err != nil {
		panic(err)
	}

	return encoded
}

// AsPlaintextResponse - creates the plaintext part of the enclave response
// It would be visible that there is an enclave response,
// but the bytes in it will still be encrypted
func AsPlaintextResponse(encResp EncryptedUserResponse) EnclaveResponse {
	return EnclaveResponse{
		EncUserResponse: encResp,
	}
}

// AsPlaintextError - generates a plaintext response containing a visible to the host error.
func AsPlaintextError(err error) EnclaveResponse {
	return EnclaveResponse{
		Err: err,
	}
}

// AsEncryptedResponse - wraps the data passed into the proper format, serializes it and encrypts it.
// It is then encoded in a plaintext response.
func AsEncryptedResponse[T any](data *T, encrypt ViewingKeyEncryptor) EnclaveResponse {
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
func AsEncryptedError(err error, encrypt ViewingKeyEncryptor) EnclaveResponse {
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
