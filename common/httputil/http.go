// Package httputil provides convenience functions for HTTP requests
package httputil

import (
	"errors"
	"github.com/marcusva/docproc/common/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Response sends a HTTP response. The passed message will be the response
// body.
func Response(w http.ResponseWriter, code int, message string) error {
	w.Header().Set("Content-Length", strconv.Itoa(len(message)))
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(code)
	_, err := w.Write([]byte(message))
	return err
}

// Error sends a HTTP 500 error response
func Error(w http.ResponseWriter, message string) error {
	return Response(w, http.StatusInternalServerError, message)
}

// BadRequest sends a HTTP 400 bad request error response
func BadRequest(w http.ResponseWriter, message string) error {
	return Response(w, http.StatusBadRequest, message)
}

// InternalServerError sends a simple HTTP 500 internal server error response
func InternalServerError(w http.ResponseWriter) error {
	return Error(w, "Internal Server Error")
}

// NotFound sends a JSON-encoded HTTP 404 resource not found Response to
// the caller.
func NotFound(w http.ResponseWriter) error {
	return Response(w, http.StatusNotFound, "resource not found")
}

// ReadBody reads the request body (if any) and returns its contents as byte
// array. If the body could not be read, an error is sent via
// SendError() or SendResponse().
func ReadBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	if r.Body == nil {
		err := errors.New("Error on reading: no body in request")
		log.Errorf("Error on reading: no body in %v", r)
		BadRequest(w, "no body in request")
		return nil, err
	}
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error on reading body: %v", err)
		Error(w, "error on reading body")
		return nil, err
	}
	return content, nil
}
