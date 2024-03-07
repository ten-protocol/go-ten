package node

import (
	"context"
	"net/http"
)

const exposedParams = "exposedParams"

type httpParamsHandler struct {
	exposedParams []string
	next          http.Handler
}

// newTenTokenHandler creates a http.Handler that extracts params
func newHttpParamsHandler(exposedParams []string, next http.Handler) http.Handler {
	return &httpParamsHandler{
		exposedParams: exposedParams,
		next:          next,
	}
}

// ServeHTTP implements http.Handler
func (handler *httpParamsHandler) ServeHTTP(out http.ResponseWriter, r *http.Request) {
	result := make(map[string]string)
	q := r.URL.Query()
	searchParams := []string{"token"}
	// for _, param := range handler.exposedParams {
	// todo wire in
	for _, param := range searchParams {
		val := q.Get(param)
		if len(val) > 0 {
			result[param] = val
		}
	}
	ctx := context.WithValue(r.Context(), exposedParams, result)
	handler.next.ServeHTTP(out, r.WithContext(ctx))
}
