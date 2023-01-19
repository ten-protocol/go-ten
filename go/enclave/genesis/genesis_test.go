package genesis

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
	"github.com/obscuronet/go-obscuro/integration/datagenerator"

	gethlog "github.com/ethereum/go-ethereum/log"
)

func TestDefaultGenesis(t *testing.T) {
	gen, err := New("")
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if len(gen.Accounts) != 3 {
		t.Fatal("unexpected number of accounts")
	}

	backingDB := rawdb.NewMemoryDatabase()
	storageDB := db.NewStorage(backingDB, nil, gethlog.New())
	stateDB, err := gen.applyAllocations(storageDB)
	if err != nil {
		t.Fatalf("unable to apply genesis allocations")
	}

	if TestnetGenesis.Accounts[0].Amount.Cmp(stateDB.GetBalance(TestnetGenesis.Accounts[0].Address)) != 0 {
		t.Fatalf("unexpected balance")
	}
}

func TestCustomGenesis(t *testing.T) {
	addr1 := datagenerator.RandomAddress()
	amt1 := datagenerator.RandomUInt64()
	addr2 := datagenerator.RandomAddress()
	amt2 := datagenerator.RandomUInt64()

	gen, err := New(
		fmt.Sprintf(
			`{"Accounts": [
				{"Address": "%s", "Amount": %d},
				{"Address": "%s", "Amount": %d}	] }
				`,
			addr1.Hex(), amt1, addr2.Hex(), amt2))
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if len(gen.Accounts) != 2 {
		t.Fatal("unexpected number of accounts")
	}

	backingDB := rawdb.NewMemoryDatabase()
	storageDB := db.NewStorage(backingDB, nil, gethlog.New())
	stateDB, err := gen.applyAllocations(storageDB)
	if err != nil {
		t.Fatalf("unable to apply genesis allocations")
	}

	if big.NewInt(int64(amt1)).Cmp(stateDB.GetBalance(addr1)) != 0 {
		t.Fatalf("unexpected balance")
	}
	if big.NewInt(int64(amt2)).Cmp(stateDB.GetBalance(addr2)) != 0 {
		t.Fatalf("unexpected balance")
	}
}
