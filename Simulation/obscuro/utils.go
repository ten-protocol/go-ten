package obscuro

import (
	"math"
	"math/rand"
	"time"
)

type Latency func() int
type ScheduledFunc func()

func RndBtw(min int, max int) int {
	return rand.Intn(max-min) + min
}

// least common ancestor of the 2 blocks
func lca(a *Block, b *Block) *Block {
	if a.height == -1 || b.height == -1 {
		return &GenesisBlock
	}
	if a.rootHash == b.rootHash {
		return a
	}
	if a.height > b.height {
		return lca(a.parent, b)
	}
	if b.height > a.height {
		return lca(a, b.parent)
	}
	return lca(a.parent, b.parent)
}

// IsAncestor return true if a is the ancestor of b
func IsAncestor(a *Block, b *Block) bool {
	if a.rootHash == b.rootHash {
		return true
	}
	if a.height >= b.height {
		return false
	}
	return IsAncestor(a, b.parent)
}

// IsRlpAncestor return true if a is the ancestor of b
func IsRlpAncestor(a *Rollup, b *Rollup) bool {
	if a.rootHash == b.rootHash {
		return true
	}
	if a.height >= b.height {
		return false
	}
	return IsRlpAncestor(a, b.parent)
}

// ScheduleInterrupt runs the function after the delay
func ScheduleInterrupt(delay int, doneCh *chan bool, fun ScheduledFunc) {
	ticker := time.NewTicker(Duration(delay))
	go func() {
		executed := false
		select {
		case <-*doneCh:
		case <-ticker.C:
			executed = true
			fun()
		}
		if executed {
			<-*doneCh
		}
		ticker.Stop()
	}()
}

// Schedule runs the function after the delay
func Schedule(delay int, fun ScheduledFunc) {
	ticker := time.NewTicker(Duration(delay))
	go func() {
		select {
		case <-ticker.C:
			ticker.Stop()
			fun()
		}
	}()
}

func Duration(us int) time.Duration {
	return time.Duration(us) * time.Microsecond
}

func generateNonce() Nonce {
	return RndBtw(0, math.MaxInt)
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
