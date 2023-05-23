package crypto

import (
	"bytes"
	"io"

	"github.com/andybalholm/brotli"
)

type DataCompressionService interface {
	Compress(blob []byte) ([]byte, error)
	Decompress(blob []byte) ([]byte, error)
}

func NewBrotliDataCompressionService() DataCompressionService {
	return &brotliDataCompressionService{}
}

type brotliDataCompressionService struct{}

func (cs *brotliDataCompressionService) Compress(in []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := brotli.NewWriterLevel(&buf, brotli.BestCompression)
	_, err := writer.Write(in)
	if closeErr := writer.Close(); err == nil {
		err = closeErr
	}
	return buf.Bytes(), err
}

func (cs *brotliDataCompressionService) Decompress(in []byte) ([]byte, error) {
	r := brotli.NewReader(bytes.NewReader(in))
	return io.ReadAll(r)
}

/*
// commented for now,  and to remove once we run some comparative testing
type gzipDataCompressionService struct{}

func (cs *gzipDataCompressionService) Compress(in []byte) ([]byte, error) {
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

func (cs *gzipDataCompressionService) Decompress(in []byte) ([]byte, error) {
	reader := bytes.NewReader(in)
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	return io.ReadAll(gz)
}
*/
