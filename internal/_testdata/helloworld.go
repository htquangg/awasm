package main

import (
	"net/http"

	"github.com/htquangg/a-wasm/sdk"
)

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", "a=b") // example of multiple headers
	w.Header().Add("Set-Cookie", "c=d")
	w.Header().Set("Date", "Tue, 15 Nov 1994 08:12:31 GMT")
	// w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"hello": "world2"}`)) // nolint
}

func main() {
	sdk.Handle(http.HandlerFunc(handle))
}
