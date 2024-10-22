package common

type GWUserDB struct {
	ID         string        `json:"id"` // Required by CosmosDB
	UserId     []byte        `json:"userId"`
	PrivateKey []byte        `json:"privateKey"`
	Accounts   []GWAccountDB `json:"accounts"` // List of Accounts
}

type GWAccountDB struct {
	AccountAddress []byte `json:"accountAddress"`
	Signature      []byte `json:"signature"`
	SignatureType  int    `json:"signatureType"`
}
