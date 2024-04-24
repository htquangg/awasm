package main

import (
	"net/http"

	"github.com/htquangg/a-wasm/sdk"
)

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-AWASM-PACKAGE", "awasm-test")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"hello": "world"}`)) // nolint
}

func main() {
	sdk.Handle(http.HandlerFunc(handle))
}
