package db

type Database interface {
	Get(id []byte) ([]byte, error)
	Store(id []byte, val []byte) error
	Delete(id []byte) error
}
