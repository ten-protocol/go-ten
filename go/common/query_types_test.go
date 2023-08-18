package common

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestPublicBlockListing_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   PublicBlock
		want    string
		wantErr bool
	}{
		{
			name: "Non-zero RollupHash",
			input: PublicBlock{BlockHeader: types.Header{
				Number:     big.NewInt(1),
				Difficulty: big.NewInt(2),
			}, RollupHash: common.BytesToHash([]byte("hello"))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshal, err := json.Marshal(tt.input)
			if err != nil {
				return
			}

			if !strings.Contains(string(marshal), "rollupHash") {
				t.Errorf("no rollupHash found")
			}

			pub := PublicBlock{}
			err = json.Unmarshal(marshal, &pub)
			if err != nil {
				t.Error(err)
			}

			require.Equal(t, pub.BlockHeader.Number.Uint64(), tt.input.BlockHeader.Number.Uint64())
			require.Equal(t, pub.BlockHeader.Difficulty.Uint64(), tt.input.BlockHeader.Difficulty.Uint64())
			require.Equal(t, pub.RollupHash.Hex(), tt.input.RollupHash.Hex())
		})
	}
}
