package responses

import (
	"encoding/json"
	"fmt"
)

type EncryptedUserResponse []byte

type EnclaveResponse struct {
	EncUserResponse EncryptedUserResponse
	Err             error
}

func (er *EnclaveResponse) Encode() []byte {
	encoded, err := json.Marshal(er)
	if err != nil {
		panic(err)
	}

	return encoded
}

func AsPlaintextResponse(encResp EncryptedUserResponse) EnclaveResponse {
	return EnclaveResponse{
		EncUserResponse: encResp,
	}
}

func AsPlaintextError(err error) EnclaveResponse {
	return EnclaveResponse{
		Err: err,
	}
}

func AsEncryptedResponse[T any](data *T, encrypt Encryptor) EnclaveResponse {
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

func AsEncryptedError(err error, encrypt Encryptor) EnclaveResponse {
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

func ToEnclaveResponse(encoded []byte) *EnclaveResponse {
	resp := EnclaveResponse{}
	err := json.Unmarshal(encoded, &resp)
	if err != nil {
		panic(err) // Todo change when stable.
	}
	return &resp
}

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
