package processors

import (
	"bytes"
	"fmt"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"net/http"
	"strconv"
	"time"
)

const (
	httpName       = "HTTPSender"
	defaultTimeout = 5 * time.Second
)

func init() {
	Register(httpName, NewHTTPSender)
}

// HTTPSender sends particular message contents to a HTTP receiver
type HTTPSender struct {
	address    string
	identifier string
	timeout    time.Duration
}

// NewHTTPSender creates a new HTTPSender
func NewHTTPSender(params map[string]string) (queue.Processor, error) {
	address, ok := params["address"]
	if !ok {
		return nil, fmt.Errorf("parameter 'address' missing")
	}
	inputid, ok := params["identifier"]
	if !ok {
		return nil, fmt.Errorf("parameter 'identifier' missing")
	}
	timeout, ok := params["timeout"]
	if !ok {
		log.Infof("parameter 'timout' not set, using default of '%v'", defaultTimeout.Seconds())
		return nil, fmt.Errorf("parameter 'input' missing")
	}
	tm, err := strconv.Atoi(timeout)
	if err != nil {
		log.Errorf("cannot convert 'timeout' to integer: %v", err)
		return nil, err
	}
	return &HTTPSender{
		address:    address,
		identifier: inputid,
		timeout:    time.Duration(tm) * time.Second,
	}, nil
}

// Name returns the name to be used in configuration files.
func (sender *HTTPSender) Name() string {
	return httpName
}

// Process processes the passed message, and sends the content identified by
// the HTTPSender's configured identifier to the address of the HTTPSender.
func (sender *HTTPSender) Process(msg *queue.Message) error {
	buf, ok := msg.Content[sender.identifier]
	if !ok {
		return fmt.Errorf("message '%s' misses identifier '%s'", msg.Metadata[queue.MetaID], sender.identifier)
	}
	var bytebuf []byte
	switch buf.(type) {
	case []byte:
		bytebuf = buf.([]byte)
	case string:
		bytebuf = []byte(buf.(string))
	default:
		log.Infof("content '%s' is not a string or byte buffer, using standard conversion", sender.identifier)
		bytebuf = []byte(fmt.Sprintf("%v", buf))
	}

	request, err := http.NewRequest("POST", sender.address, bytes.NewBuffer(bytebuf))
	if err != nil {
		log.Errorf("could not create HTTP request: %v", err)
		return err
	}
	request.Header.Add("Content-Length", strconv.Itoa(len(bytebuf)))
	// TODO: configurable header types
	request.Header.Add("Content-Type", "text/plain")

	client := &http.Client{}
	client.Timeout = sender.timeout

	result, err := client.Do(request)
	if err != nil {
		log.Errorf("could not send HTTP request: %v", err)
		return err
	}

	// TODO: handle responses properly
	switch result.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusAccepted:
		return nil
	default:
		log.Errorf("invalid status: '%s' on response %v", result.Status, result)
		return fmt.Errorf("invalid status: '%s' on response %v", result.Status, result)
	}
}
