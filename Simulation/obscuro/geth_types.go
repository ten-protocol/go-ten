package obscuro

const (
	// HashLength is the expected length of the hash (in bytes)
	HashLength = 32
	// AddressLength is the expected length of the address (in bytes)
	AddressLength = 20
	// BlockNumberLength length of uint64 big endian
	BlockNumberLength = 8
	// IncarnationLength length of uint64 for contract incarnations
	IncarnationLength = 8
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte
