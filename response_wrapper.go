package murlog

import "net/http"

type responseWriter struct {
	Status      int

	http.ResponseWriter
	wroteHeader bool
}

func WrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{Status: 200,  ResponseWriter: w}
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.Status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}
