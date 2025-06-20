package enclavedb

import (
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type FilterTopics [][]gethcommon.Hash

func (ft FilterTopics) IsEmpty() bool {
	return len(ft) == 0
}

func (ft FilterTopics) EventTypes() []gethcommon.Hash {
	if ft.IsEmpty() {
		return nil
	}
	return ft[0]
}

func (ft FilterTopics) HasEventTypes() bool {
	return len(ft.EventTypes()) > 0
}

func (ft FilterTopics) TopicsOnPos(pos int) []gethcommon.Hash {
	if ft.IsEmpty() {
		return nil
	}
	if len(ft) < pos+1 {
		return nil
	}
	return ft[pos]
}

func (ft FilterTopics) HasTopicsOnPos(pos int) bool {
	return len(ft.TopicsOnPos(pos)) > 0
}

// utility for sql query building
func repeat(token string, sep string, count int) string {
	elems := make([]string, count)
	for i := 0; i < count; i++ {
		elems[i] = token
	}
	return strings.Join(elems, sep)
}
