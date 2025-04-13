package response

import (
	"bytes"
	"fmt"

	"github.com/codecrafters-io/http-server-starter-go/app/gcho/headers"
)

type Response struct {
	Headers headers.Headers
	Status  int
	body    *bytes.Buffer
}

func New() Response {
	return Response{
		Headers: headers.New(),
		body:    bytes.NewBuffer(nil),
	}
}

func (r *Response) Header(key string, value string) {
	r.Headers.Set(key, value)
}

func (r *Response) Write(b []byte) (n int, err error) {
	r.Header("Content-Length", fmt.Sprintf("%v", len(b)))
	if r.Status == 0 {
		r.Status = 200
	}
	if r.Headers.Get("Content-Type") == "" {
		r.Header("Content-Type", "text/plain")
	}
	return r.body.Write(b)
}

func (r *Response) Bytes() []byte {
	return r.body.Bytes()
}
