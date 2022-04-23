package txhandler

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
)

// EncodeToString encodes a byte array to a string
func EncodeToString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// DecodeFromString decodes a string to a byte array
func DecodeFromString(in string) []byte {
	bytesStr, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return bytesStr
}

// Compress compresses the byte array using gzip
func Compress(in []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(in); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Decompress decompresses the byte array using gzip
func Decompress(in []byte) ([]byte, error) {
	reader := bytes.NewReader(in)
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	return ioutil.ReadAll(gz)
}
