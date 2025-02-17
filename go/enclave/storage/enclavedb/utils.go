package enclavedb

import (
	"strings"
)

// utility for sql query building
func repeat(token string, sep string, count int) string {
	elems := make([]string, count)
	for i := 0; i < count; i++ {
		elems[i] = token
	}
	return strings.Join(elems, sep)
}
