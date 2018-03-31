package processors

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

const (
	cvmessage = `{
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

func TestNewContentValidator(t *testing.T) {
	_, err := NewContentValidator(nil)
	assert.Err(t, err)

	params := map[string]string{}
	_, err = NewContentValidator(params)
	assert.Err(t, err)

	params = map[string]string{"rules": "test/brokenrules.json"}
	_, err = NewContentValidator(params)
	assert.Err(t, err)

	params = map[string]string{"rules": "test/invalid"}
	_, err = NewContentValidator(params)
	assert.Err(t, err)

	params = map[string]string{"rules": "test/xml-template.tpl"}
	_, err = NewContentValidator(params)
	assert.Err(t, err)

	params = map[string]string{"rules": "test/cvrules.json"}
	_, err = NewContentValidator(params)
	assert.NoErr(t, err)
}

func TestContentValidatorCreate(t *testing.T) {
	params := map[string]string{
		"type":  "ContentValidator",
		"rules": "test/testrules.json",
	}
	proc, err := Create(params)
	assert.FailOnErr(t, err)
	assert.Equal(t, proc.Name(), "ContentValidator")
}

func TestContentValidatorName(t *testing.T) {
	cv, err := NewContentValidator(map[string]string{"rules": "test/testrules.json"})
	assert.FailOnErr(t, err)
	assert.Equal(t, cv.Name(), "ContentValidator")
}

func TestContentValidatorProcess(t *testing.T) {
	cv, err := NewContentValidator(map[string]string{"rules": "test/testrules.json"})
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(cvmessage))
	assert.FailOnErr(t, err)
	assert.Err(t, cv.Process(msg))

	cv, err = NewContentValidator(map[string]string{"rules": "test/cvrules.json"})
	assert.FailOnErr(t, err)

	msg, err = queue.MsgFromJSON([]byte(cvmessage))
	assert.FailOnErr(t, err)
	assert.NoErr(t, cv.Process(msg))
}
