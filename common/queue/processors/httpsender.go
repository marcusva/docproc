package processors

import (
	"bytes"
	"fmt"
	"github.com/marcusva/docproc/common/data"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"net/http"
	"net/url"
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
	address  string
	url      *url.URL
	readFrom string
	timeout  time.Duration
}

// NewHTTPSender creates a new HTTPSender
func NewHTTPSender(params map[string]string) (queue.Processor, error) {
	address, ok := params["address"]
	if !ok {
		return nil, fmt.Errorf("parameter 'address' missing")
	}
	if _, err := url.Parse(address); err != nil {
		return nil, fmt.Errorf("parameter 'address' is invalid: %v", err)
	}

	inputid, ok := params["read.from"]
	if !ok {
		return nil, fmt.Errorf("parameter 'read.from' missing")
	}
	var tm uint64
	timeout, ok := params["timeout"]
	if !ok {
		log.Infof("parameter 'timout' not set, using default of '%v'", defaultTimeout.Seconds())
		tm = uint64(defaultTimeout.Seconds())
	} else {
		var err error
		tm, err = strconv.ParseUint(timeout, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return &HTTPSender{
		address:  address,
		readFrom: inputid,
		timeout:  time.Duration(tm) * time.Second,
	}, nil
}

// Name returns the name to be used in configuration files.
func (sender *HTTPSender) Name() string {
	return httpName
}

// Process processes the passed message, and sends the content identified by
// the HTTPSender's configured readFrom to the address of the HTTPSender.
func (sender *HTTPSender) Process(msg *queue.Message) error {
	buf, ok := msg.Content[sender.readFrom]
	if !ok {
		return fmt.Errorf("message '%s' misses identifier '%s'", msg.Metadata[queue.MetaID], sender.readFrom)
	}
	bytebuf, err := data.Bytes(buf)
	if err != nil {
		return err
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
