package node

import (
	"context"
	"net/http"

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
	ctx := context.WithValue(r.Context(), rpc.GWTokenKey{}, val)
	handler.next.ServeHTTP(out, r.WithContext(ctx))
	handler.next.ServeHTTP(out, r)
}
