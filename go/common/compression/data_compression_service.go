package compression

import (
	"bytes"
	"errors"
	"io"

	"github.com/andybalholm/brotli"
)

type DataCompressionService interface {
	// CompressRollup - uses the maximum compression level, because the final size matters when publishing to Ethereum
	CompressRollup(blob []byte) ([]byte, error)
	// CompressBatch - uses the default compression level, because the compression is for the efficiency of the p2p transfer
	CompressBatch(blob []byte) ([]byte, error)
	Decompress(blob []byte) ([]byte, error)
}

func NewBrotliDataCompressionService() DataCompressionService {
	return &brotliDataCompressionService{}
}

type brotliDataCompressionService struct{}

func (cs *brotliDataCompressionService) CompressRollup(blob []byte) ([]byte, error) {
	return cs.compress(blob, brotli.BestCompression)
}

func (cs *brotliDataCompressionService) CompressBatch(blob []byte) ([]byte, error) {
	return cs.compress(blob, brotli.DefaultCompression)
}

func (cs *brotliDataCompressionService) Decompress(in []byte) ([]byte, error) {
	r := brotli.NewReader(bytes.NewReader(in))
	limitedReader := io.LimitReader(r, 1024*1024) // 1MB limit
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1)
	if n, _ := r.Read(buf); n > 0 {
		return nil, errors.New("decompressed size exceeded. Sequencer does not produce more than 1MB of data")
	}

	return data, nil
}

func (cs *brotliDataCompressionService) compress(in []byte, level int) ([]byte, error) {
	var buf bytes.Buffer
	writer := brotli.NewWriterLevel(&buf, level)
	_, err := writer.Write(in)
	if closeErr := writer.Close(); err == nil {
		err = closeErr
	}
	return buf.Bytes(), err
}

/*
// commented for now,  and to remove once we run some comparative testing
type gzipDataCompressionService struct{}

func (cs *gzipDataCompressionService) CompressRollup(in []byte) ([]byte, error) {
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
