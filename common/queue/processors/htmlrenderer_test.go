package processors_test

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/queue/processors"
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

func TestHTMLRenderer(t *testing.T) {
	_, err := processors.NewHTMLRenderer(nil)
	assert.Err(t, err)

	params := map[string]string{}
	_, err = processors.NewHTMLRenderer(params)
	assert.Err(t, err)

	params["templates"] = "test/html/*.tpl"
	_, err = processors.NewHTMLRenderer(params)
	assert.Err(t, err)

	params["store.in"] = "html"
	_, err = processors.NewHTMLRenderer(params)
	assert.Err(t, err)

	params["templateroot"] = "main"
	_, err = processors.NewHTMLRenderer(params)
	assert.NoErr(t, err)

	params["templates"] = "invalid"
	_, err = processors.NewHTMLRenderer(params)
	assert.Err(t, err)
}

func TestHTMLRendererCreate(t *testing.T) {
	params := map[string]string{
		"type":         "HTMLRenderer",
		"templates":    "test/html/*.tpl",
		"store.in":     "html",
		"templateroot": "main",
	}
	proc, err := processors.Create(params)
	assert.FailOnErr(t, err)
	assert.Equal(t, proc.Name(), "HTMLRenderer")
}

func TestHTMLRendererProcess(t *testing.T) {
	params := map[string]string{
		"templates":    "test/html/*.tpl",
		"store.in":     "html",
		"templateroot": "main",
	}
	html, err := processors.NewHTMLRenderer(params)
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
		"store.in":     "html",
		"templateroot": "main",
	}
	html, err := processors.NewHTMLRenderer(params)
	assert.FailOnErr(t, err)
	assert.Equal(t, html.Name(), "HTMLRenderer")
}
