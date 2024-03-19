package common

type AccountDB struct {
	AccountAddress []byte
	Signature      []byte
	SignatureType  int
}

type UserDB struct {
	UserID     []byte
	PrivateKey []byte
}
