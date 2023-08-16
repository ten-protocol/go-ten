package common

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TestEntity struct {
	RollupHash common.Hash `json:"rollupHash"`

	types.Header
}

func (te *TestEntity) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		RollupHash string       `json:"rollupHash"`
		Header     types.Header `json:"header"`
	}{
		RollupHash: te.RollupHash.Hex(),
		Header:     te.Header,
	})
}

func TestPublicBlockListing_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   TestEntity
		want    string
		wantErr bool
	}{
		{
			name:    "Non-zero RollupHash",
			input:   TestEntity{Header: types.Header{Number: big.NewInt(1)}, RollupHash: common.BytesToHash([]byte("hello"))},
			want:    `{"Number":1,"rollupHash":"0x1234567890abcdef"}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.MarshalJSON()
			if !strings.Contains(string(got), "rollupHash") {
				t.Errorf("fck")
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}
