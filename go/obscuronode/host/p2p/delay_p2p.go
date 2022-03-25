package p2p

import (
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type delayP2PImpl struct {
	p2p   P2P
	delay func() uint64
}

func NewDelayP2P(p2p P2P, delay func() uint64) P2P {
	return &delayP2PImpl{
		p2p:   p2p,
		delay: delay,
	}
}

func (d *delayP2PImpl) Listen(t chan nodecommon.EncryptedTx, r chan obscurocommon.EncodedRollup) {
	d.p2p.Listen(t, r)
}

func (d *delayP2PImpl) StopListening() {
	d.StopListening()
}

func (d *delayP2PImpl) BroadcastTx(bytes []byte) {
	obscurocommon.Schedule(d.delay()/2, func() {
		d.p2p.BroadcastTx(bytes)
	})
}

func (d *delayP2PImpl) BroadcastRollup(bytes []byte) {
	obscurocommon.Schedule(d.delay()/2, func() {
		d.p2p.BroadcastRollup(bytes)
	})
}
