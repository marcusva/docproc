package processors

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestNewHTTPSender(t *testing.T) {
	_, err := NewHTTPSender(nil)
	assert.Err(t, err)

	params := map[string]string{}
	_, err = NewHTTPSender(params)
	assert.Err(t, err)

	params["address"] = "localhost"
	_, err = NewHTTPSender(params)
	assert.Err(t, err)

	params["read.from"] = "body"
	_, err = NewHTTPSender(params)
	assert.NoErr(t, err)

	params["timeout"] = "-1"
	_, err = NewHTTPSender(params)
	assert.Err(t, err)

	params["timeout"] = "123"
	_, err = NewHTTPSender(params)
	assert.NoErr(t, err)

	params["address"] = "::some##invalid?!!!\\data"
	_, err = NewHTTPSender(params)
	assert.Err(t, err)

}

func TestHTTPSenderCreate(t *testing.T) {
	params := map[string]string{
		"type":      "HTTPSender",
		"address":   "127.0.0.1",
		"read.from": "body",
		"timeout":   "2",
	}
	proc, err := Create(params)
	assert.FailOnErr(t, err)
	assert.Equal(t, proc.Name(), "HTTPSender")
}

func TestHTTPSenderName(t *testing.T) {
	params := map[string]string{
		"address":   "127.0.0.1",
		"read.from": "body",
		"timeout":   "2",
	}
	sender, err := NewHTTPSender(params)
	assert.FailOnErr(t, err)
	assert.Equal(t, sender.Name(), "HTTPSender")
}

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
		"address":   server.URL,
		"read.from": "body",
		"timeout":   "2",
	}
	sender, err := NewHTTPSender(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(httpmessage))
	assert.FailOnErr(t, err)
	assert.FailOnErr(t, sender.Process(msg))
}

func TestHTTPSenderProcessInvalid(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1001 * time.Millisecond)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	params := map[string]string{
		"address":   server.URL,
		"read.from": "noexist",
		"timeout":   "1",
	}
	sender, err := NewHTTPSender(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(httpmessage))
	assert.FailOnErr(t, err)

	assert.Err(t, sender.Process(msg))

	params["read.from"] = "body"
	sender, err = NewHTTPSender(params)
	assert.FailOnErr(t, err)
	assert.Err(t, sender.Process(msg))

}
