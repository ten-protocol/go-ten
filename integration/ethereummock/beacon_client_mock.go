package ethereummock

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const kzgBlobSize = 131072

// BeaconMock presents a beacon-node in testing, without leading any chain-building.
// This merely serves a fake beacon API, and holds on to blocks,
// to complement the actual block-building to happen in testing (e.g. through the fake consensus geth module).
type BeaconMock struct {
	log log.Logger

	blobs     map[uint64][]*kzg4844.Blob
	blobsLock sync.Mutex

	beaconSrv         *http.Server
	beaconAPIListener net.Listener

	genesisTime    uint64
	secondsPerSlot uint64
	port           int
}

func NewBeaconMock(log log.Logger, genesisTime uint64, secondsPerSlot uint64, port int) *BeaconMock {
	return &BeaconMock{
		log:            log,
		genesisTime:    genesisTime,
		secondsPerSlot: secondsPerSlot,
		port:           port,
		blobs:          make(map[uint64][]*kzg4844.Blob),
	}
}

func (f *BeaconMock) Start(host string) error {
	address := fmt.Sprintf("%s:%d", host, f.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to open tcp listener for http beacon api server: %w", err)
	}
	f.beaconAPIListener = listener

	mux := new(http.ServeMux)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	mux.HandleFunc("/eth/v1/beacon/genesis", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(&ethadapter.APIGenesisResponse{
			Data: ethadapter.ReducedGenesisData{
				GenesisTime: ethadapter.Uint64String(f.genesisTime),
			},
		})
		if err != nil {
			f.log.Error("genesis handler err", "err", err)
		}
	})
	mux.HandleFunc("/eth/v1/config/spec", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(&ethadapter.APIConfigResponse{
			Data: ethadapter.ReducedConfigData{
				SecondsPerSlot: ethadapter.Uint64String(f.secondsPerSlot),
			},
		})
		if err != nil {
			f.log.Error("config handler err", "err", err)
		}
	})
	mux.HandleFunc("/eth/v1/beacon/blob_sidecars/", func(w http.ResponseWriter, r *http.Request) {
		blockID := strings.TrimPrefix(r.URL.Path, "/eth/v1/beacon/blob_sidecars/")
		slot, err := strconv.ParseUint(blockID, 10, 64)
		if err != nil {
			f.log.Error("could not parse block id from request", "url", r.URL.Path, "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		blobs, err := f.LoadBlobs(slot)
		if errors.Is(err, ethereum.NotFound) {
			f.log.Error("failed to load blobs - not found", "slot", slot, "err", err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			f.log.Error("failed to load blobs", "slot", slot, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		query := r.URL.Query()
		rawIndices := query["indices"]
		indices := make([]uint64, 0, len(blobs))
		if len(rawIndices) == 0 {
			for i := range blobs {
				indices = append(indices, uint64(i))
			}
		} else {
			for _, raw := range rawIndices {
				ix, err := strconv.ParseUint(raw, 10, 64)
				if err != nil {
					f.log.Error("could not parse index from request", "url", r.URL)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				indices = append(indices, ix)
			}
		}

		var mockBeaconBlockRoot [32]byte
		mockBeaconBlockRoot[0] = 42
		binary.LittleEndian.PutUint64(mockBeaconBlockRoot[32-8:], slot)
		sidecars := make([]*ethadapter.APIBlobSidecar, len(indices))
		for i, ix := range indices {
			if ix >= uint64(len(blobs)) {
				f.log.Error("blob index from request is out of range", "url", r.URL)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			blob := blobs[ix]

			commitment, err := kzg4844.BlobToCommitment(blob)
			if err != nil {
				f.log.Error("failed to compute KZG commitment", "blobIndex", ix, "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			proof, err := kzg4844.ComputeBlobProof(blob, commitment)
			if err != nil {
				f.log.Error("failed to compute KZG proof", "blobIndex", ix, "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			sidecars[i] = &ethadapter.APIBlobSidecar{
				Index:         ethadapter.Uint64String(ix),
				Blob:          *blob,
				KZGCommitment: ethadapter.Bytes48(commitment[:]),
				KZGProof:      ethadapter.Bytes48(proof[:]),
				SignedBlockHeader: ethadapter.SignedBeaconBlockHeader{
					Message: ethadapter.BeaconBlockHeader{
						StateRoot: mockBeaconBlockRoot,
						Slot:      ethadapter.Uint64String(slot),
					},
				},
			}
		}
		if err := json.NewEncoder(w).Encode(&ethadapter.APIGetBlobSidecarsResponse{Data: sidecars}); err != nil {
			f.log.Error("blobs handler err", "err", err)
		}
	})
	f.beaconSrv = &http.Server{
		Handler:           mux,
		ReadTimeout:       time.Second * 20,
		ReadHeaderTimeout: time.Second * 20,
		WriteTimeout:      time.Second * 20,
		IdleTimeout:       time.Second * 20,
	}
	go func() {
		if err := f.beaconSrv.Serve(f.beaconAPIListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			f.log.Error("failed to start fake-pos beacon server for blobs testing", "err", err)
		}
	}()
	return nil
}

// StoreBlobs stores the array of blobs against the slot number.
func (f *BeaconMock) StoreBlobs(slot uint64, blobs []*kzg4844.Blob) error {
	for _, blob := range blobs {
		commitment, _ := kzg4844.BlobToCommitment(blob)
		versionedHash := ethadapter.KZGToVersionedHash(commitment)
		println("MockBeacon storing blob hash: ", versionedHash.Hex(), " at slot: ", slot)
	}
	f.blobsLock.Lock()
	defer f.blobsLock.Unlock()

	if len(blobs) > 0 {
		f.blobs[slot] = append(f.blobs[slot], blobs...)
	}
	return nil
}

// LoadBlobs retrieves the array of blobs for a given slot.
func (f *BeaconMock) LoadBlobs(slot uint64) ([]*kzg4844.Blob, error) {
	f.blobsLock.Lock()
	defer f.blobsLock.Unlock()

	blobs, exists := f.blobs[slot]
	if !exists {
		return nil, fmt.Errorf("no blobs found for slot %d: %w", slot, ethereum.NotFound)
	}
	return blobs, nil
}

func (f *BeaconMock) Close() error {
	var out error
	if f.beaconSrv != nil {
		out = errors.Join(out, f.beaconSrv.Close())
	}
	if f.beaconAPIListener != nil {
		out = errors.Join(out, f.beaconAPIListener.Close())
	}
	return out
}

func (f *BeaconMock) BeaconAddr() string {
	return "http://" + f.beaconAPIListener.Addr().String()
}
