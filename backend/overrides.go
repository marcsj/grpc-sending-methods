package backend

import (
	"net/http"
)

// ParamsToHeaders is used to pass all params from a websocket request into headers in an http request
func ParamsToHeaders(incoming *http.Request, outgoing *http.Request) *http.Request {
	for k, v := range incoming.URL.Query() {
		outgoing.Header.Set(k, v[0])
	}
	return outgoing
}

// MatchAllHeaders can be used on a grpc-gateway mux to pass all headers into a context
func MatchAllHeaders(key string) (string, bool) {
	return key, true
}

