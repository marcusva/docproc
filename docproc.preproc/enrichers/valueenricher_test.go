package enrichers

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

const (
	message = `{
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

func TestNewValueRuleEnricher(t *testing.T) {
	params := map[string]string{"norules": "1234"}
	_, err := NewValueEnricher(params)
	assert.FailIf(t, err == nil, "NewValueEricher() must fail, if no 'rule' arg is provided")

	params = map[string]string{"rules": "test/norules.json"}
	_, err = NewValueEnricher(params)
	assert.FailIf(t, err == nil, "NewValueEricher() must fail, if there is no rules file")

	params = map[string]string{"rules": "test/brokenrules.json"}
	_, err = NewValueEnricher(params)
	assert.FailIf(t, err == nil, "NewValueEricher() must fail, if the rules are broken")

	params = map[string]string{"rules": "test/testrules.json"}
	_, err = NewValueEnricher(params)
	assert.FailOnErr(t, err)
}

func TestValueEnricherProcess(t *testing.T) {
	p := map[string]string{"rules": "test/testrules.json"}
	ve, err := NewValueEnricher(p)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(message))
	assert.FailOnErr(t, err)
	assert.FailOnErr(t, ve.Process(msg))

	rval, ok := msg.Content["DOCTYPE"]
	assert.FailIfNot(t, ok, "no 'DOCTYPE' found")
	assert.Equal(t, rval, "INVOICE")

	filename, ok := msg.Content["filename"]
	assert.FailIfNot(t, ok, "no 'filename' found")
	assert.Equal(t, filename, "fn-100112.html")

	multivar, ok := msg.Content["multi-var"]
	assert.FailIfNot(t, ok, "no 'multi-var' found")
	assert.Equal(t, multivar, "100112.10394.00-INVOICE")
}

func TestQueueProcessing(t *testing.T) {
	wq, err := queue.CreateWQ("memory", map[string]string{"topic": "q"})
	assert.FailOnErr(t, err)

	p := map[string]string{"rules": "test/testrules.json"}
	ve, err := NewValueEnricher(p)
	assert.FailOnErr(t, err)

	writer := queue.NewWriter(wq)
	writer.AddProcessor(ve)

	msg, err := queue.MsgFromJSON([]byte(message))
	assert.FailOnErr(t, err)

	assert.FailOnErr(t, writer.Consume(msg))

	rval, ok := msg.Content["DOCTYPE"]
	assert.FailIfNot(t, ok, "no 'DOCTYPE' found")
	assert.Equal(t, rval, "INVOICE")
}

func TestInvalidRules(t *testing.T) {
	p := map[string]string{"rules": "test/invalidrules.json"}
	ve, err := NewValueEnricher(p)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(message))
	assert.FailOnErr(t, err)
	assert.Err(t, ve.Process(msg))
}
