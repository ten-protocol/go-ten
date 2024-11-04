package common

type GWUserDB struct {
	UserId     []byte        `json:"userId"`
	PrivateKey []byte        `json:"privateKey"`
	Accounts   []GWAccountDB `json:"accounts"`
}

type GWAccountDB struct {
	AccountAddress []byte `json:"accountAddress"`
	Signature      []byte `json:"signature"`
	SignatureType  int    `json:"signatureType"`
}
