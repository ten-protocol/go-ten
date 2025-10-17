package cosmosdb

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	gethcommon "github.com/ethereum/go-ethereum/common"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

const (
	sessionKeyActivityContainerName = "session_key_activities"
	// shard doc id format: sk_shard_<index>
	skShardPrefix = "sk_shard_"
	// default shard count (mirrors metrics DEFAULT_SHARD_COUNT)
	DEFAULT_SK_SHARD_COUNT = 50
	// 2MB CosmosDB item limit safeguard
	twoMBLimitBytes = 2 * 1024 * 1024
)

type sessionKeyActivityStorageCosmosDB struct {
	client     *azcosmos.Client
	container  *azcosmos.ContainerClient
	shardCount int
}

type sessionKeyActivityDTO struct {
	ID          string                      `json:"id"`
	Items       []sessionKeyActivityItemDTO `json:"items"`
	ShardIndex  int                         `json:"shardIndex"`
	LastUpdated string                      `json:"lastUpdated"`
}

type sessionKeyActivityItemDTO struct {
	Addr       string `json:"addr"`
	UserIDHex  string `json:"userId"`
	LastActive string `json:"lastActive"`
}

func NewSessionKeyActivityStorage(connectionString string) (*sessionKeyActivityStorageCosmosDB, error) {
	client, err := azcosmos.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create CosmosDB client: %w", err)
	}

	ctx := context.Background()
	_, err = client.CreateDatabase(ctx, azcosmos.DatabaseProperties{ID: DATABASE_NAME}, nil)
	if err != nil && !strings.Contains(err.Error(), "Conflict") {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	container, err := client.NewContainer(DATABASE_NAME, sessionKeyActivityContainerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get session key activities container: %w", err)
	}

	return &sessionKeyActivityStorageCosmosDB{client: client, container: container, shardCount: DEFAULT_SK_SHARD_COUNT}, nil
}

func (s *sessionKeyActivityStorageCosmosDB) Load() ([]wecommon.SessionKeyActivity, error) {
	ctx := context.Background()
	result := make([]wecommon.SessionKeyActivity, 0)

	for i := 0; i < s.shardCount; i++ {
		shardID := s.getShardDocumentIDByIndex(i)
		pk := azcosmos.NewPartitionKeyString(shardID)
		resp, err := s.container.ReadItem(ctx, pk, shardID, nil)
		if err != nil {
			if strings.Contains(err.Error(), "NotFound") {
				continue
			}
			return nil, err
		}
		var dto sessionKeyActivityDTO
		if err := json.Unmarshal(resp.Value, &dto); err != nil {
			return nil, err
		}
		for _, it := range dto.Items {
			addr := gethcommon.HexToAddress(it.Addr)
			userID, _ := hex.DecodeString(strings.TrimPrefix(it.UserIDHex, "0x"))
			t, err := time.Parse(time.RFC3339, it.LastActive)
			if err != nil {
				continue
			}
			result = append(result, wecommon.SessionKeyActivity{Addr: addr, UserID: userID, LastActive: t})
		}
	}
	return result, nil
}

func (s *sessionKeyActivityStorageCosmosDB) Save(items []wecommon.SessionKeyActivity) error {
	ctx := context.Background()

	// group by shard index
	byShard := make(map[int][]wecommon.SessionKeyActivity)
	for _, it := range items {
		idx := s.shardIndexForAddress(it.Addr)
		byShard[idx] = append(byShard[idx], it)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	for idx, shardItems := range byShard {
		shardID := s.getShardDocumentIDByIndex(idx)
		dto := sessionKeyActivityDTO{
			ID:          shardID,
			Items:       make([]sessionKeyActivityItemDTO, 0, len(shardItems)),
			ShardIndex:  idx,
			LastUpdated: now,
		}
		for _, it := range shardItems {
			dto.Items = append(dto.Items, sessionKeyActivityItemDTO{
				Addr:       it.Addr.Hex(),
				UserIDHex:  hex.EncodeToString(it.UserID),
				LastActive: it.LastActive.UTC().Format(time.RFC3339),
			})
		}
		b, err := json.Marshal(dto)
		if err != nil {
			return err
		}
		if len(b) > twoMBLimitBytes {
			return fmt.Errorf("session key activity shard %d document exceeds 2MB limit (%d bytes)", idx, len(b))
		}
		pk := azcosmos.NewPartitionKeyString(shardID)
		if _, err := s.container.UpsertItem(ctx, pk, b, nil); err != nil {
			return err
		}
	}
	return nil
}

// shardIndexForAddress computes shard index using FNV-32a hash of address bytes
func (s *sessionKeyActivityStorageCosmosDB) shardIndexForAddress(addr gethcommon.Address) int {
	h := fnv.New32a()
	_, _ = h.Write(addr.Bytes())
	return int(h.Sum32()) % s.shardCount
}

func (s *sessionKeyActivityStorageCosmosDB) getShardDocumentIDByIndex(index int) string {
	return fmt.Sprintf("%s%d", skShardPrefix, index)
}

func (s *sessionKeyActivityStorageCosmosDB) getShardIndexFromID(id string) int {
	var index int
	fmt.Sscanf(id, skShardPrefix+"%d", &index)
	return index
}
