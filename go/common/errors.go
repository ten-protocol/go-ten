package common

import (
	"encoding/json"
	"fmt"
)

type Encryptor func([]byte) ([]byte, error)

type SystemError error

type EncryptedResponse[Data any] struct {
	result  *Data
	err     error
	Encoded []byte
}

func EmptyResponse[T any]() *EncryptedResponse[T] {
	return &EncryptedResponse[T]{
		result: nil,
		err:    nil,
	}
}

func (er *EncryptedResponse[T]) WithResult(val T) *EncryptedResponse[T] {
	er.result = &val
	return er
}

func (er *EncryptedResponse[T]) WithError(err error) *EncryptedResponse[T] {
	er.err = err
	return er
}

type transportMessage[T any] struct {
	Result *T
	Err    *string
}

func (er *EncryptedResponse[T]) Finalize(encrypt Encryptor) (*EncryptedResponse[T], error) {
	encodableStruct := &transportMessage[T]{
		Result: er.result,
	}
	if er.err != nil {
		errStr := er.err.Error()
		encodableStruct.Err = &errStr
	}

	bytes, err := json.Marshal(encodableStruct)
	if err != nil {
		return nil, err
	}
	result := EmptyResponse[T]()
	result.Encoded, err = encrypt(bytes)
	return result, err
}

func (er *EncryptedResponse[T]) Decode(encoded []byte) (*T, error) {
	decoded := &transportMessage[T]{}
	err := json.Unmarshal(encoded, decoded)
	if err != nil {
		return nil, err
	}

	var userError error
	if decoded.Err != nil {
		userError = fmt.Errorf(*decoded.Err)
	}

	return decoded.Result, userError
}

func DecodeEncryptedBytes[T any](encoded []byte) (*T, error) {
	decoded := &transportMessage[T]{}
	err := json.Unmarshal(encoded, decoded)
	if err != nil {
		return nil, err
	}

	var userError error
	if decoded.Err != nil {
		userError = fmt.Errorf(*decoded.Err)
	}

	return decoded.Result, userError
}

type EncryptedBytesResponse[Data any] []byte

func (ebr *EncryptedBytesResponse[T]) AsResult(result T) *EncryptedBytesResponse[T] {
	response := transportMessage[T]{
		Result: &result,
	}

	encodedBytes, _ := json.Marshal(response)
	*ebr = encodedBytes
	return ebr
}

func (ebr *EncryptedBytesResponse[T]) AsError(err error) *EncryptedBytesResponse[T] {
	errorStr := err.Error()
	response := transportMessage[T]{
		Err: &errorStr,
	}

	encodedBytes, _ := json.Marshal(response)
	*ebr = encodedBytes
	return ebr
}

func (ebr *EncryptedBytesResponse[T]) Encrypt(encrypt Encryptor) (EncryptedBytesResponse[T], error) {
	encryptedBytes, err := encrypt(*ebr)
	*ebr = encryptedBytes
	return *ebr, err
}
