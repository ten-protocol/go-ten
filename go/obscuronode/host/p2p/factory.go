package p2p

// Factory allows building different types of P2P networks
// this is commonly used for its injection / lazy-instantiation capabilities
type Factory struct {
	p2p func(ourAddress string, allAddresses []string) P2P
}

// NewP2P returns a new P2P instance
func (f *Factory) NewP2P(address string, addresses []string) P2P {
	return f.p2p(address, addresses)
}

// NewP2PFactory creates a new Factory givng the P2P init function
func NewP2PFactory(p2p func(ourAddress string, allAddresses []string) P2P) *Factory {
	return &Factory{p2p: p2p}
}
