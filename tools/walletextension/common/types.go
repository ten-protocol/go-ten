package common

type AccountDB struct {
	AccountAddress []byte
	Signature      []byte
}

type UserDB struct {
	UserID     []byte
	PrivateKey []byte
}
