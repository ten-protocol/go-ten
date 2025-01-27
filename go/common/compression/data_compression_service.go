package compression

import (
	"bytes"
	"errors"
	"fmt"
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

func NewBrotliDataCompressionService(decompressedSizeLimit int64) DataCompressionService {
	return &brotliDataCompressionService{
		decompressedSizeLimit: decompressedSizeLimit,
	}
}

type brotliDataCompressionService struct {
	decompressedSizeLimit int64
}

func (cs *brotliDataCompressionService) CompressRollup(blob []byte) ([]byte, error) {
	return cs.compress(blob, brotli.BestCompression)
}

func (cs *brotliDataCompressionService) CompressBatch(blob []byte) ([]byte, error) {
	return cs.compress(blob, brotli.DefaultCompression)
}

func (cs *brotliDataCompressionService) Decompress(in []byte) ([]byte, error) {
	if in == nil {
		return nil, errors.New("input is nil")
	}

	r := brotli.NewReader(bytes.NewReader(in))
	// Limit the decompressed size to the decompressedSizeLimit
	limitedReader := io.LimitReader(r, cs.decompressedSizeLimit)
	// Read up to the decompressedSizeLimit
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	if len(data) != int(cs.decompressedSizeLimit) {
		return data, nil
	}

	// Verify that the decompressed data is within the decompressedSizeLimit;
	// if we manage to read again, then the decompressed data is larger than the decompressedSizeLimit
	buf := make([]byte, 1)
	n, readErr := r.Read(buf)
	if readErr != nil && readErr != io.EOF {
		return nil, fmt.Errorf("decompression verification error: %w", readErr)
	}

	if n > 0 {
		return data, errors.New("decompressed size exceeded")
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
