package middleware

import (
	"io"
	"net/http"
)

type CustomResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
