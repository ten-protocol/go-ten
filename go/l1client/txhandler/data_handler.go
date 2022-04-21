package txhandler

import (
	"bytes"
	"compress/gzip"
	b64 "encoding/base64"
	"io/ioutil"
)

func EncodeToString(bytes []byte) string {
	return b64.StdEncoding.EncodeToString(bytes)
}

func DecodeFromString(in string) []byte {
	bytes, err := b64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return bytes
}

func Compress(in []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(in); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return b.Bytes()
}

func Decompress(in []byte) []byte {
	reader := bytes.NewReader(in)
	gz, err := gzip.NewReader(reader)
	if err != nil {
		panic(err)
	}
	defer gz.Close()

	output, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}

	return output
}
