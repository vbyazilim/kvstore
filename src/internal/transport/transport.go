package transport

import "net/http"

// KVStoreHTTPHandler defines /store/ http handler behaviours.
type KVStoreHTTPHandler interface {
	Set(http.ResponseWriter, *http.Request)
}
