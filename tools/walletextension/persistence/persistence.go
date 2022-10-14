package persistence

import (
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/rpc"
)

const (
	obscuroDirName           = ".obscuro"
	persistenceFileName      = "wallet_extension_persistence"
	persistenceNumComponents = 4
	persistenceIdxHost       = 0
	persistenceIdxAccount    = 1
	persistenceIdxViewingKey = 2
	persistenceIdxSignedKey  = 3
)

// Persistence handles the persistence of viewing keys.
type Persistence struct {
	filePath string // The path of the file used to store the submitted viewing keys
	hostAddr string // The address of the host the keys are being persisted for
	logger   gethlog.Logger
}

func NewPersistence(hostAddr string, persistenceFilePath string, logger gethlog.Logger) *Persistence {
	// Sets up the persistence file and returns its path. Defaults to the user's home directory if the path is empty.
	if persistenceFilePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic("cannot create persistence file as user's home directory is not defined")
		}
		obscuroDir := filepath.Join(homeDir, obscuroDirName)
		err = os.MkdirAll(obscuroDir, 0o777)
		if err != nil {
			panic(fmt.Sprintf("could not create %s directory in user's home directory", obscuroDirName))
		}

		persistenceFilePath = filepath.Join(obscuroDir, persistenceFileName)
	}

	_, err := os.OpenFile(persistenceFilePath, os.O_CREATE|os.O_RDONLY, 0o644)
	if err != nil {
		logger.Crit("could not create persistence file. ", log.ErrKey, err)
	}

	return &Persistence{
		filePath: persistenceFilePath,
		hostAddr: hostAddr,
		logger:   logger,
	}
}

// PersistViewingKey persists a viewing key to disk.
func (p *Persistence) PersistViewingKey(viewingKey *rpc.ViewingKey) {
	viewingPrivateKeyBytes := crypto.FromECDSA(viewingKey.PrivateKey.ExportECDSA())

	record := []string{
		p.hostAddr,
		viewingKey.Account.Hex(),
		// We encode the bytes as hex to ensure there are no unintentional line breaks to make parsing the file harder.
		hex.EncodeToString(viewingPrivateKeyBytes),
		hex.EncodeToString(viewingKey.SignedKey),
	}

	persistenceFile, err := os.OpenFile(p.filePath, os.O_APPEND|os.O_WRONLY, 0o644)
	defer persistenceFile.Close() //nolint:staticcheck
	if err != nil {
		p.logger.Error("could not open persistence file. ", log.ErrKey, err)
	}

	writer := csv.NewWriter(persistenceFile)
	defer writer.Flush()
	err = writer.Write(record)
	if err != nil {
		p.logger.Error("failed to write viewing key to persistence file. ", log.ErrKey, err)
	}
}

// LoadViewingKeys loads any persisted viewing keys from disk for the given host. Viewing keys for other hosts are ignored.
func (p *Persistence) LoadViewingKeys() map[common.Address]*rpc.ViewingKey {
	viewingKeys := make(map[common.Address]*rpc.ViewingKey)

	persistenceFile, err := os.OpenFile(p.filePath, os.O_RDONLY, 0o644)
	defer persistenceFile.Close() //nolint:staticcheck
	if err != nil {
		p.logger.Error("could not open persistence file. ", log.ErrKey, err)
	}

	reader := csv.NewReader(persistenceFile)
	records, err := reader.ReadAll()
	if err != nil {
		p.logger.Error("could not read records from persistence file. ", log.ErrKey, err)
	}

	for _, record := range records {
		// TODO - Determine strategy for invalid persistence entries - delete? Warn? Shutdown? For now, we log a warning.
		if len(record) != persistenceNumComponents {
			p.logger.Warn(fmt.Sprintf("persistence file entry did not have expected number of components: %s", record))
			continue
		}

		persistedHostAddr := record[persistenceIdxHost]
		if persistedHostAddr != p.hostAddr {
			p.logger.Info(fmt.Sprintf("skipping persistence file entry for another host. Current host is %s, entry was for %s", p.hostAddr, persistedHostAddr))
			continue
		}

		account := common.HexToAddress(record[persistenceIdxAccount])
		viewingKeyPrivateHex := record[persistenceIdxViewingKey]
		viewingKeyPrivateBytes, err := hex.DecodeString(viewingKeyPrivateHex)
		if err != nil {
			p.logger.Warn(fmt.Sprintf("could not decode the following viewing private key from hex in the persistence file: %s", viewingKeyPrivateHex))
			continue
		}
		viewingKeyPrivate, err := crypto.ToECDSA(viewingKeyPrivateBytes)
		if err != nil {
			p.logger.Warn(fmt.Sprintf("could not convert the following viewing private key bytes to ECDSA in the persistence file: %s", viewingKeyPrivateHex))
			continue
		}
		signedKeyHex := record[persistenceIdxSignedKey]
		signedKey, err := hex.DecodeString(signedKeyHex)
		if err != nil {
			p.logger.Warn(fmt.Sprintf("could not decode the following signed key from hex in the persistence file: %s", signedKeyHex))
			continue
		}

		viewingKey := rpc.ViewingKey{
			Account:    &account,
			PrivateKey: ecies.ImportECDSA(viewingKeyPrivate),
			PublicKey:  crypto.CompressPubkey(&viewingKeyPrivate.PublicKey),
			SignedKey:  signedKey,
		}
		viewingKeys[account] = &viewingKey
	}

	p.logReRegisteredViewingKeys(viewingKeys)

	return viewingKeys
}

// Logs and prints the accounts for which we are re-registering viewing keys.
func (p *Persistence) logReRegisteredViewingKeys(viewingKeys map[common.Address]*rpc.ViewingKey) {
	if len(viewingKeys) == 0 {
		return
	}

	var accounts []string //nolint:prealloc
	for account := range viewingKeys {
		accounts = append(accounts, account.Hex())
	}

	msg := fmt.Sprintf("Re-registering persisted viewing keys for the following addresses: %s",
		strings.Join(accounts, ", "))
	p.logger.Info(msg)
	fmt.Println(msg)
}
