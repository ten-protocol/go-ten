package node

import (
	"context"
	"net/http"
	"strings"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
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
	if val == "" {
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
	if n != 40 && n != 42 {
		return false
	}
	if n == 42 {
		if !(strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X")) {
			return false
		}
		s = s[2:]
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
