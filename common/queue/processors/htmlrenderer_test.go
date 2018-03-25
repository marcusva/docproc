package processors

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"strings"
	"testing"
)

const (
	htmlmessage = `{
	"metadata": {
		"batch": 1517607828,
		"created": "2018-02-02T22:43:48.0220047+01:00"
	},
	"content": {
		"CITY": "New York",
		"CUSTNO": "100112",
		"DATE": "2017-04-07",
		"DOCTYPE": "INVOICE",
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

func TestHTMLRendererProcess(t *testing.T) {
	params := map[string]string{
		"templates":    "test/html/*.tpl",
		"identifier":   "html",
		"templateroot": "main",
	}
	html, err := NewHTMLRenderer(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(htmlmessage))
	assert.FailOnErr(t, err)

	err = html.Process(msg)
	assert.FailOnErr(t, err)
	data, ok := msg.Content["html"]
	assert.FailIfNot(t, ok, "html output missing")
	assert.FailIfNot(t, strings.Contains(data.(string), "<title>Invoice</title>"), "template error")
}

func TestHTMLRenderName(t *testing.T) {
	params := map[string]string{
		"templates":    "test/html/*.tpl",
		"identifier":   "html",
		"templateroot": "main",
	}
	html, err := NewHTMLRenderer(params)
	assert.FailOnErr(t, err)

	assert.Equal(t, html.Name(), "HTMLRenderer")
}
