package cosmosdb

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	gethcommon "github.com/ethereum/go-ethereum/common"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/encryption"
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

// SessionKeyActivityStorage interface defines the session key activity storage operations
type SessionKeyActivityStorage interface {
	Load() ([]wecommon.SessionKeyActivity, error)
	Save([]wecommon.SessionKeyActivity) error
}

type sessionKeyActivityStorageCosmosDB struct {
	client     *azcosmos.Client
	container  *azcosmos.ContainerClient
	shardCount int
	encryptor  encryption.Encryptor
}

type sessionKeyActivityDTO struct {
	ID          string                      `json:"id"`
	Items       []sessionKeyActivityItemDTO `json:"items"`
	ShardIndex  int                         `json:"shardIndex"`
	LastUpdated string                      `json:"lastUpdated"`
}

type sessionKeyActivityItemDTO struct {
	Addr       []byte    `json:"addr"`       // 20 bytes (common.Address.Bytes())
	UserID     []byte    `json:"userId"`     // variable length
	LastActive time.Time `json:"lastActive"` // RFC3339 timestamp
}

func NewSessionKeyActivityStorage(connectionString string, encryptionKey []byte) (SessionKeyActivityStorage, error) {
	// Create encryptor
	encryptor, err := encryption.NewEncryptor(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create encryptor: %w", err)
	}

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

	return &sessionKeyActivityStorageCosmosDB{
		client:     client,
		container:  container,
		shardCount: DEFAULT_SK_SHARD_COUNT,
		encryptor:  *encryptor,
	}, nil
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

		// Unmarshal into EncryptedDocument
		var doc EncryptedDocument
		if err := json.Unmarshal(resp.Value, &doc); err != nil {
			return nil, fmt.Errorf("failed to unmarshal document: %w", err)
		}

		// Decrypt the data
		data, err := s.encryptor.Decrypt(doc.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data: %w", err)
		}

		// Unmarshal decrypted JSON into DTO
		var dto sessionKeyActivityDTO
		if err := json.Unmarshal(data, &dto); err != nil {
			return nil, fmt.Errorf("failed to unmarshal session key activity data: %w", err)
		}

		for _, it := range dto.Items {
			addr := gethcommon.BytesToAddress(it.Addr)
			userID := it.UserID
			result = append(result, wecommon.SessionKeyActivity{Addr: addr, UserID: userID, LastActive: it.LastActive})
		}
	}
	return result, nil
}

// Save performs a full replacement of all stored session key activity data with the provided snapshot.
// This method uses a two-phase approach to ensure data consistency:
//
// Phase 1 (Clear): Writes empty documents to ALL shards to remove any stale or deleted entries.
//
//	This ensures that if a session key was deleted from memory, it's also removed from storage.
//
// Phase 2 (Write): Writes the actual activity data to shards that contain items.
//
//	Items are distributed across shards based on a hash of their address to ensure even distribution.
//
// This full replacement strategy is used because:
//   - The caller provides a complete snapshot of all current activities (from memory)
//   - We need to ensure deleted entries are removed from storage
//   - It's simpler and more reliable than tracking incremental changes
//
// The method groups items by shard index, clears all shards, then writes data only to shards
// that contain items. Each shard document is validated against CosmosDB's 2MB size limit.
func (s *sessionKeyActivityStorageCosmosDB) Save(items []wecommon.SessionKeyActivity) error {
	ctx := context.Background()
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Group items by shard index for efficient batch writing.
	// Each session key address is hashed to determine which shard it belongs to,
	// ensuring even distribution across all shards.
	itemsByShardIndex := make(map[int][]wecommon.SessionKeyActivity)
	for _, item := range items {
		shardIdx := s.shardIndexForAddress(item.Addr)
		itemsByShardIndex[shardIdx] = append(itemsByShardIndex[shardIdx], item)
	}

	// Phase 1: Clear all shards by writing empty documents.
	// This removes any stale data from shards that no longer have active items.
	if err := s.clearAllShards(ctx, timestamp); err != nil {
		return fmt.Errorf("failed to clear shards: %w", err)
	}

	// Phase 2: Write actual data to shards that contain items.
	// Only shards with items are written to, leaving others empty from phase 1.
	for shardIdx, shardItems := range itemsByShardIndex {
		if err := s.writeShardData(ctx, shardIdx, shardItems, timestamp); err != nil {
			return fmt.Errorf("failed to write shard %d: %w", shardIdx, err)
		}
	}

	return nil
}

// clearAllShards writes empty documents to all shards to remove stale data.
// This is phase 1 of the Save operation, ensuring that deleted session keys
// are removed from storage. Each empty document includes metadata (ID, ShardIndex, LastUpdated)
// but no activity items.
func (s *sessionKeyActivityStorageCosmosDB) clearAllShards(ctx context.Context, timestamp string) error {
	for shardIdx := 0; shardIdx < s.shardCount; shardIdx++ {
		shardID := s.getShardDocumentIDByIndex(shardIdx)
		dto := sessionKeyActivityDTO{
			ID:          shardID,
			Items:       make([]sessionKeyActivityItemDTO, 0),
			ShardIndex:  shardIdx,
			LastUpdated: timestamp,
		}

		// Marshal to JSON
		dtoJSON, err := json.Marshal(dto)
		if err != nil {
			return fmt.Errorf("failed to marshal empty shard %d: %w", shardIdx, err)
		}

		// Encrypt the data
		ciphertext, err := s.encryptor.Encrypt(dtoJSON)
		if err != nil {
			return fmt.Errorf("failed to encrypt empty shard %d: %w", shardIdx, err)
		}

		// Create EncryptedDocument
		doc := EncryptedDocument{
			ID:   shardID,
			Data: ciphertext,
		}

		// Marshal EncryptedDocument to JSON and validate size against CosmosDB's 2MB limit
		b, err := json.Marshal(doc)
		if err != nil {
			return fmt.Errorf("failed to marshal encrypted document for empty shard %d: %w", shardIdx, err)
		}
		if len(b) > twoMBLimitBytes {
			return fmt.Errorf("session key activity shard %d empty doc exceeds 2MB limit (%d bytes)", shardIdx, len(b))
		}

		// Upsert the empty document to clear any existing data in this shard
		pk := azcosmos.NewPartitionKeyString(shardID)
		if _, err := s.container.UpsertItem(ctx, pk, b, nil); err != nil {
			return fmt.Errorf("failed to upsert empty shard %d: %w", shardIdx, err)
		}
	}
	return nil
}

// writeShardData writes activity items to a specific shard.
// This is phase 2 of the Save operation, writing actual data to shards that contain items.
// The method converts session key activities to DTOs and validates the document size
// against CosmosDB's 2MB limit before writing.
func (s *sessionKeyActivityStorageCosmosDB) writeShardData(ctx context.Context, shardIdx int, items []wecommon.SessionKeyActivity, timestamp string) error {
	shardID := s.getShardDocumentIDByIndex(shardIdx)
	dto := sessionKeyActivityDTO{
		ID:          shardID,
		Items:       make([]sessionKeyActivityItemDTO, 0, len(items)),
		ShardIndex:  shardIdx,
		LastUpdated: timestamp,
	}

	// Convert session key activities to DTOs for storage
	for _, item := range items {
		dto.Items = append(dto.Items, sessionKeyActivityItemDTO{
			Addr:       item.Addr.Bytes(),
			UserID:     item.UserID,
			LastActive: item.LastActive,
		})
	}

	// Marshal to JSON
	dtoJSON, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal shard %d: %w", shardIdx, err)
	}

	// Encrypt the data
	ciphertext, err := s.encryptor.Encrypt(dtoJSON)
	if err != nil {
		return fmt.Errorf("failed to encrypt shard %d: %w", shardIdx, err)
	}

	// Create EncryptedDocument
	doc := EncryptedDocument{
		ID:   shardID,
		Data: ciphertext,
	}

	// Marshal EncryptedDocument to JSON and validate size against CosmosDB's 2MB limit
	b, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted document for shard %d: %w", shardIdx, err)
	}
	if len(b) > twoMBLimitBytes {
		return fmt.Errorf("session key activity shard %d document exceeds 2MB limit (%d bytes)", shardIdx, len(b))
	}

	// Upsert the document with actual data
	pk := azcosmos.NewPartitionKeyString(shardID)
	if _, err := s.container.UpsertItem(ctx, pk, b, nil); err != nil {
		return fmt.Errorf("failed to upsert shard %d: %w", shardIdx, err)
	}

	return nil
}

// shardIndexForAddress computes the shard index for a given address using FNV-32a hash.
// The address bytes are hashed and then modulo'd by the shard count to ensure
// even distribution across all available shards. This deterministic hashing ensures
// the same address always maps to the same shard.
func (s *sessionKeyActivityStorageCosmosDB) shardIndexForAddress(addr gethcommon.Address) int {
	h := fnv.New32a()
	_, _ = h.Write(addr.Bytes())
	return int(h.Sum32()) % s.shardCount
}

// getShardDocumentIDByIndex generates the document ID for a shard given its index.
// The format is "sk_shard_<index>" (e.g., "sk_shard_0", "sk_shard_1", etc.).
// This ID is used both as the document ID and partition key in CosmosDB.
func (s *sessionKeyActivityStorageCosmosDB) getShardDocumentIDByIndex(index int) string {
	return fmt.Sprintf("%s%d", skShardPrefix, index)
}
