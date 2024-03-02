package main

import (
	"net/http"

	"github.com/htquangg/a-wasm/sdk"
)

func handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello Alex HUYNH!!!"))
}

func main() {
	sdk.Handle(http.HandlerFunc(handle))
}
