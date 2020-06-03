package processors_test

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/queue/processors"
	"github.com/marcusva/docproc/common/testing/assert"
)

const (
	httpmessage = `{
	"metadata": {
		"batch": 1517607828,
		"created": "2018-02-02T22:43:48.0220047+01:00"
	},
	"content": {
		"body": "some content",
		"http-headers": {
			"Content-Type": "application/custom",
			"Authorization": "Basic QWzzzz="
		}
	}
}
`
)

func TestNewHTTPSender(t *testing.T) {
	_, err := processors.NewHTTPSender(nil)
	assert.Err(t, err)

	params := map[string]string{}
	_, err = processors.NewHTTPSender(params)
	assert.Err(t, err)

	params["address"] = "localhost"
	_, err = processors.NewHTTPSender(params)
	assert.Err(t, err)

	params["read.from"] = "body"
	_, err = processors.NewHTTPSender(params)
	assert.NoErr(t, err)

	params["timeout"] = "-1"
	_, err = processors.NewHTTPSender(params)
	assert.Err(t, err)

	params["timeout"] = "123"
	_, err = processors.NewHTTPSender(params)
	assert.NoErr(t, err)

	params["headers"] = "http-headers"
	_, err = processors.NewHTTPSender(params)
	assert.NoErr(t, err)

	params["address"] = "::some##invalid?!!!\\data"
	_, err = processors.NewHTTPSender(params)
	assert.Err(t, err)

}

func TestHTTPSenderCreate(t *testing.T) {
	params := map[string]string{
		"type":      "HTTPSender",
		"address":   "127.0.0.1",
		"read.from": "body",
		"timeout":   "2",
	}
	proc, err := processors.Create(params)
	assert.FailOnErr(t, err)
	assert.Equal(t, proc.Name(), "HTTPSender")
}

func TestHTTPSenderName(t *testing.T) {
	params := map[string]string{
		"address":   "127.0.0.1",
		"read.from": "body",
		"timeout":   "2",
	}
	sender, err := processors.NewHTTPSender(params)
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
	sender, err := processors.NewHTTPSender(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(httpmessage))
	assert.FailOnErr(t, err)
	assert.FailOnErr(t, sender.Process(msg))
}

func TestHTTPSenderProcessHTTPS(t *testing.T) {
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
	server := httptest.NewTLSServer(okHandler)
	defer server.Close()

	cert, _ := x509.ParseCertificate(server.TLS.Certificates[0].Certificate[0])
	certpool := x509.NewCertPool()
	certpool.AddCert(cert)
	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certpool,
		},
	}

	params := map[string]string{
		"address":   server.URL,
		"read.from": "body",
		"timeout":   "2",
	}
	sender, err := processors.NewHTTPSender(params)
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
	sender, err := processors.NewHTTPSender(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(httpmessage))
	assert.FailOnErr(t, err)

	assert.Err(t, sender.Process(msg))

	params["read.from"] = "body"
	sender, err = processors.NewHTTPSender(params)
	assert.FailOnErr(t, err)
	assert.Err(t, sender.Process(msg))
}

func TestHTTPSenderProcessCustomHeaders(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Basic QWzzzz=" {
			w.WriteHeader(500)
			return
		}
		if r.Header.Get("Content-Type") != "application/custom" {
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
		"headers":   "http-headers",
		"timeout":   "2",
	}
	sender, err := processors.NewHTTPSender(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(httpmessage))
	assert.FailOnErr(t, err)
	assert.FailOnErr(t, sender.Process(msg))
}
