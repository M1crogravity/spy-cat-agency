package main

import (
	"bytes"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode    int
	body          bytes.Buffer
	headerWritten bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.headerWritten {
		rw.statusCode = code
		rw.headerWritten = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.headerWritten {
		rw.statusCode = http.StatusOK
		rw.headerWritten = true
	}
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}
