package node

import (
	"context"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

type httpParamsHandler struct {
	exposedParam string
	next         http.Handler
}

// newTenTokenHandler creates a http.Handler that extracts params
func newHTTPParamsHandler(exposedParam string, next http.Handler) http.Handler {
	return &httpParamsHandler{
		exposedParam: exposedParam,
		next:         next,
	}
}

// ServeHTTP implements http.Handler
func (handler *httpParamsHandler) ServeHTTP(out http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	val := q.Get(handler.exposedParam)
	if len(val) == 0 {
		// Attempt to extract token from path as /v1/<TOKEN>
		// Do not alter routing; only read the segment.
		// Expected segments: ["", "v1", "<TOKEN>", ...]
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) >= 3 && parts[1] == "v1" {
			candidate := parts[2]
			if isLikelyHexToken(candidate) {
				val = candidate
			}
		}
	}
	ctx := context.WithValue(r.Context(), rpc.GWTokenKey{}, val)
	handler.next.ServeHTTP(out, r.WithContext(ctx))
	handler.next.ServeHTTP(out, r)
}

// isLikelyHexToken performs a lightweight validation for a user token
// Accepts both 0x-prefixed (42 chars) and non-prefixed (40 chars) lowercase/uppercase hex
func isLikelyHexToken(s string) bool {
	n := len(s)
	if n != common.MessageUserIDLen && n != common.MessageUserIDLenWithPrefix {
		return false
	}
	if n == common.MessageUserIDLenWithPrefix {
		if !(strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X")) {
			return false
		}
		s = s[2:]
	}
	// Use standard library hex decoder for validation
	if _, err := hex.DecodeString(s); err != nil {
		return false
	}
	return true
}
