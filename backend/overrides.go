package main

import (
	"net/http"
)

// paramsToHeaders is used to pass all params from a websocket request into headers in an http request
func paramsToHeaders(incoming *http.Request, outgoing *http.Request) *http.Request {
	for k, v := range incoming.URL.Query() {
		outgoing.Header.Set(k, v[0])
	}
	return outgoing
}

// matchAllHeaders can be used on a grpc-gateway mux to pass all headers into a context
func matchAllHeaders(key string) (string, bool) {
	return key, true
}

