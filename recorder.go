package rest

import (
	"net/http"
)

type recorderResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (self *recorderResponseWriter) WriteHeader(code int) {
	self.Header().Add("X-Powered-By", "go-json-rest")
	self.ResponseWriter.WriteHeader(code)
	self.statusCode = code
	self.wroteHeader = true
}

func (self *recorderResponseWriter) Flush() {
	if !self.wroteHeader {
		self.WriteHeader(http.StatusOK)
	}
	flusher := self.ResponseWriter.(http.Flusher)
	flusher.Flush()
}

func (self *recorderResponseWriter) CloseNotify() <-chan bool {
	notifier := self.ResponseWriter.(http.CloseNotifier)
	return notifier.CloseNotify()
}

func (self *recorderResponseWriter) Write(b []byte) (int, error) {

	if !self.wroteHeader {
		self.WriteHeader(http.StatusOK)
	}

	return self.ResponseWriter.Write(b)
}

func (self *ResourceHandler) recorderWrapper(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		writer := &recorderResponseWriter{w, 0, false}

		// call the handler
		h(writer, r)

		self.env.setVar(r, "statusCode", writer.statusCode)
	}
}
