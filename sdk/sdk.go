package sdk

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/protobuf/proto"

	"github.com/htquangg/awasm/internal/protocluster/grains/messages"
)

type ResponseWriter struct {
	header     http.Header
	buffer     *bytes.Buffer
	statusCode int
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
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

	w := &ResponseWriter{
		header:     make(http.Header),
		buffer:     new(bytes.Buffer),
		statusCode: 200,
	}
	r, err := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range req.Header {
		r.Header[k] = v.Fields
	}
	h.ServeHTTP(w, r) // execute the user's handler
	os.Stdout.Write(w.buffer.Bytes())

	js, err := json.Marshal(w.Header())
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(js)

	buf := make([]byte, 1<<8)
	binary.LittleEndian.PutUint32(buf[0:4], uint32(w.statusCode))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(w.buffer.Len()))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(len(js)))
	os.Stdout.Write(buf)
}
