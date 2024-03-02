package sdk

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/htquangg/a-wasm/internal/protocluster/grains/messages"

	"google.golang.org/protobuf/proto"
)

type ResponseWriter struct {
	buffer     bytes.Buffer
	statusCode int
}

func (*ResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *ResponseWriter) Write(b []byte) (n int, err error) {
	return w.buffer.Write(b)
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.statusCode = status
}

func Handle(h http.Handler) {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	var req messages.HTTPRequest
	if err := proto.Unmarshal(b, &req); err != nil {
		log.Fatal(err)
	}

	w := &ResponseWriter{}
	r, err := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range req.Header {
		r.Header[k] = v.Fields
	}
	h.ServeHTTP(w, r) // execute the user's handler
	os.Stdout.Write(w.buffer.Bytes())

	buf := make([]byte, 1<<3)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(w.statusCode))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(w.buffer.Len()))
	os.Stdout.Write(buf)
}
