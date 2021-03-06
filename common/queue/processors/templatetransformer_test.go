package processors_test

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/queue/processors"
	"github.com/marcusva/docproc/common/testing/assert"
	"strings"
	"testing"
)

const (
	tplmessage = `{
    "metadata": {
        "batch": 1517607828,
        "created": "2018-02-02T22:43:48.0220047+01:00"
    },
    "content": {
        "CITY": "New York",
        "CUSTNO": "100112",
        "DOCTYPE": "INVOICE",
        "DATE": "2017-04-07",
        "FIRSTNAME": "John",
        "GROSS": "12386.86",
        "LASTNAME": "Doe",
        "NET": "10394.00",
        "STREET": "Example Lane 384",
        "ZIP": "10006"
    }
}`
)

func TestNewTemplateTransformer(t *testing.T) {
	params := map[string]string{}
	_, err := processors.NewTemplateTransformer(params)
	assert.Err(t, err)

	params["templates"] = "test//*.tpl"
	_, err = processors.NewTemplateTransformer(params)
	assert.Err(t, err)

	params["templateroot"] = "main"
	_, err = processors.NewTemplateTransformer(params)
	assert.Err(t, err)

	params["store.in"] = "_output_"
	_, err = processors.NewTemplateTransformer(params)
	assert.FailOnErr(t, err)

	params["templates"] = "invalid//*.tpl"
	_, err = processors.NewTemplateTransformer(params)
	assert.Err(t, err)
}

func TestTemplateTransformerCreate(t *testing.T) {
	params := map[string]string{
		"type":         "TemplateTransformer",
		"templates":    "test//*.tpl",
		"store.in":     "_xml_",
		"templateroot": "main",
	}
	proc, err := processors.Create(params)
	assert.FailOnErr(t, err)
	assert.Equal(t, proc.Name(), "TemplateTransformer")
}

func TestTemplateTransformerProcess(t *testing.T) {
	params := map[string]string{
		"templates":    "test//*.tpl",
		"store.in":     "_xml_",
		"templateroot": "main",
	}
	lnt, err := processors.NewTemplateTransformer(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(tplmessage))
	assert.FailOnErr(t, err)

	err = lnt.Process(msg)
	assert.FailOnErr(t, err)

	data, ok := msg.Content["_xml_"]
	assert.FailIfNot(t, ok, "no _xml_ section found")
	assert.FailIfNot(t, strings.Contains(data.(string), "<invoice>"), "<invoice> not found")
}
