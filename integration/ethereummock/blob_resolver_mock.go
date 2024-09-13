package ethereummock

//import (
//	"context"
//	gethcommon "github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/core/types"
//	"github.com/ethereum/go-ethereum/crypto/kzg4844"
//	gethlog "github.com/ethereum/go-ethereum/log"
//	"github.com/ten-protocol/go-ten/go/host/l1"
//	"os"
//)
//
//type mockBlobResolver struct {
//	beaconPort int
//}
//
//func NewMockBlobResolver(port int) l1.BlobResolver {
//	return &mockBlobResolver{
//		beaconPort: port,
//	}
//}
//
//func (m *mockBlobResolver) IsMock() bool {
//	return true
//}
//
//func (m *mockBlobResolver) FetchBlobs(_ context.Context, b *types.Header, _ []gethcommon.Hash) ([]*kzg4844.Blob, error) {
//	// Create a logger
//	l := gethlog.New()
//
//	// Create a temporary directory for blob storage
//	tempDir, err := os.MkdirTemp("", "mock-blob-resolver-*")
//	if err != nil {
//		l.Crit("failed to create temp directory", "err", err)
//		return nil, err
//	}
//
//	// Initialize and start the BeaconMock
//	beaconApi := NewBeaconMock(l, tempDir, uint64(0), uint64(0))
//	err = beaconApi.Start("localhost", m.beaconPort)
//	if err != nil {
//		l.Crit("failed to start BeaconMock", "err", err)
//		return nil, err
//	}
//	defer beaconApi.Close()
//
//	slot := b.Number.Uint64()
//
//	//FIXME need to store them first
//	bundle, err := beaconApi.LoadBlobsBundle(slot)
//	if err != nil {
//		l.Error("failed to load blobs bundle", "slot", slot, "err", err)
//		return nil, err
//	}
//
//	var blobs []*kzg4844.Blob
//	for _, blobData := range bundle.Blobs {
//		blob := &kzg4844.Blob{}
//		copy(blob[:], blobData)
//		blobs = append(blobs, blob)
//	}
//
//	return blobs, nil
//}
