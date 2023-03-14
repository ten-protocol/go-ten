package node

type Node interface {
	Start() error
	Upgrade() error
}
