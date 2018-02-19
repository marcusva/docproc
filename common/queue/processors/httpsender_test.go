package processors

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	httpmessage = `{
	"metadata": {
		"batch": 1517607828,
		"created": "2018-02-02T22:43:48.0220047+01:00"
	},
	"content": {
		"body": "some content"
	}
}
`
)

func TestHTTPSenderProcess(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain" {
			w.WriteHeader(500)
			return
		}
		buf, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if string(buf) != "some content" {
			w.WriteHeader(500)
		}
		w.WriteHeader(200)
	})
	server := httptest.NewServer(okHandler)
	defer server.Close()

	params := map[string]string{
		"address":    server.URL,
		"identifier": "body",
		"timeout":    "2",
	}
	sender, err := NewHTTPSender(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(httpmessage))
	assert.FailOnErr(t, err)

	assert.FailOnErr(t, sender.Process(msg))
}
