package input

import (
	"bytes"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	rawmessage = `{
		"metadata": {
			"batch": 1517607828,
			"created": "2018-02-02T22:43:48.0220047+01:00"
		},
		"content": {
			"CITY": "New York",
			"CUSTNO": "100112",
			"DATE": "2017-04-07",
			"FIRSTNAME": "John",
			"GROSS": "12386.86",
			"LASTNAME": "Doe",
			"NET": "10394.00",
			"STREET": "Example Lane 384",
			"ZIP": "10006"
		}
	}
	`
)

func TestRawHandler(t *testing.T) {
	wq, _ := queue.CreateWQ("memory", map[string]string{"topic": "q"})
	params := map[string]string{}

	_, err := NewRawHandler(nil, params)
	assert.Err(t, err)
	_, err = NewRawHandler(wq, params)
	assert.Err(t, err)

	for _, inv := range []string{"", "banana", "10.234", "-1", "0"} {
		params["maxsize"] = inv
		_, err = NewRawHandler(wq, params)
		assert.Err(t, err)
	}

	params["maxsize"] = "12"
	_, err = NewRawHandler(wq, params)
	assert.NoErr(t, err)
}

func TestTransform(t *testing.T) {
	wq, _ := queue.CreateWQ("memory", map[string]string{"topic": "q"})
	raw, _ := NewRawHandler(wq, map[string]string{"maxsize": "1"})

	w := httptest.NewRecorder()
	raw.Transform(w, httptest.NewRequest("POST", "http://localhost/raw", nil))
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	w = httptest.NewRecorder()
	raw.Transform(w, httptest.NewRequest("POST", "http://localhost/raw", bytes.NewBufferString("invalid content")))
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	w = httptest.NewRecorder()
	raw.Transform(w, httptest.NewRequest("POST", "http://localhost/raw", bytes.NewBufferString(rawmessage)))
	assert.Equal(t, w.Result().StatusCode, http.StatusInternalServerError)

	wq.Open()
	w = httptest.NewRecorder()
	raw.Transform(w, httptest.NewRequest("POST", "http://localhost/raw", bytes.NewBufferString(rawmessage)))
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	w = httptest.NewRecorder()
	req := httptest.NewRequest("POST", "http://localhost/raw", bytes.NewBufferString(rawmessage))
	req.ContentLength = -1
	raw.Transform(w, req)
	assert.Equal(t, w.Result().StatusCode, http.StatusOK)

	buf := make([]byte, 2*1024)
	rand.Read(buf)
	w = httptest.NewRecorder()
	raw.Transform(w, httptest.NewRequest("POST", "http://localhost/raw", bytes.NewBuffer(buf)))
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)

	w = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "http://localhost/raw", bytes.NewBuffer(buf))
	req.ContentLength = -1
	raw.Transform(w, req)
	assert.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
}
