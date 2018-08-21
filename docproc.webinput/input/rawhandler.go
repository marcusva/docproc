package input

import (
	"encoding/json"
	"fmt"
	"github.com/marcusva/docproc/common/httputil"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"net/http"
	"strconv"
)

const (
	rawName = "RawHandler"
)

func init() {
	Register(rawName, NewRawHandler)
}

// RawHandler is a simple queue.Message consumer for message sent via HTTP.
type RawHandler struct {
	*queue.Writer

	// maxSize is the maximum allowed message size in bytes
	maxSize int64
}

// Transform reads a queue.Message as JSON from the request body and places it
// on the WriteQueue of the RawHandler.
func (raw *RawHandler) Transform(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Received HTTP request")

	reader := r.Body
	if r.ContentLength == -1 {
		// Unknown content length, limit it to the configured maximum
		reader = http.MaxBytesReader(w, r.Body, raw.maxSize)
	} else if r.ContentLength > raw.maxSize {
		httputil.BadRequest(w, "message exceeds the allowed file size")
	}

	var tmp queue.Message

	dec := json.NewDecoder(reader)
	if err := dec.Decode(&tmp); err != nil {
		log.Errorf("could not decode the request body: %v", err)
		httputil.BadRequest(w, "invalid request body")
		return
	}

	// We'll throw away all possible Metadata and create it from scratch.
	msg := queue.NewMessage(tmp.Content)

	if err := raw.Consume(msg); err != nil {
		log.Errorf("could not publish message to queue: %v", err)
		httputil.Error(w, "internal server error on processing")
	}
	httputil.Response(w, http.StatusOK, "OK")
}

// NewRawHandler creates a new RawHandler
func NewRawHandler(wq queue.WriteQueue, params map[string]string) (WebHandler, error) {
	if wq == nil {
		return nil, fmt.Errorf("write queue must not be nil")
	}

	size, ok := params["maxsize"]
	if !ok {
		return nil, fmt.Errorf("parameter 'maxsize' missing")
	}
	maxSize, err := strconv.ParseInt(size, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("parameter 'maxsize' is not an integer")
	}
	if maxSize <= 0 {
		return nil, fmt.Errorf("parameter 'maxsize' must be greater than 0")
	}
	maxSize *= 1024
	return &RawHandler{
		Writer:  queue.NewWriter(wq, nil),
		maxSize: maxSize,
	}, nil
}
