package enclave

import "github.com/obscuronet/obscuro-playground/go/obscuronode/db"

type BlockStateDB struct {
	innerDB db.Database
}

func NewBlockStateDB(db db.Database) *BlockStateDB {
	return &BlockStateDB{
		innerDB: db,
	}
}

func (b *BlockStateDB) Get(id []byte) ([]byte, error) {
	// HANDLE ENCODING/DECODING
}

func (b *BlockStateDB) Store(id []byte, val BlockState) error {
	// HANDLE ENCODING/DECODING
}

func (b *BlockStateDB) Delete(id []byte) error {
	// HANDLE ENCODING/DECODING
}

type RollupDB struct {
	innerDB db.Database
}

func NewRollupDB(db db.Database) *RollupDB {
	return &RollupDB{
		innerDB: db,
	}
}

func (b *RollupDB) Get(id []byte) ([]byte, error) {
	// HANDLE ENCODING/DECODING
}

func (b *RollupDB) Store(id []byte, val Rollup) error {
	// HANDLE ENCODING/DECODING
}

func (b *RollupDB) Delete(id []byte) error {
	// HANDLE ENCODING/DECODING
}
