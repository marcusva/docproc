package processors

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"runtime"
	"testing"
)

const (
	cmdmsg = `{
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

func TestNewCommandProc(t *testing.T) {
	_, err := NewCommandProc(nil)
	assert.Err(t, err)

	params := map[string]string{}
	_, err = NewCommandProc(params)
	assert.Err(t, err)

	params["read.from"] = "CITY"
	_, err = NewCommandProc(params)
	assert.Err(t, err)

	params["store.in"] = "cmdfield"
	_, err = NewCommandProc(params)
	assert.Err(t, err)

	params["exec"] = "cmd"
	_, err = NewCommandProc(params)
	assert.NoErr(t, err)
}

func TestCommandProcName(t *testing.T) {
	params := map[string]string{
		"read.from": "CITY",
		"store.in":  "cmdfield",
		"exec":      "cmd",
	}
	cmd, err := NewCommandProc(params)
	assert.FailOnErr(t, err)
	assert.Equal(t, cmd.Name(), "CommandProc")
}

func TestCommandProcProcess(t *testing.T) {
	params := map[string]string{
		"read.from": "CITY",
		"store.in":  "cmdfield",
	}
	switch runtime.GOOS {
	case "windows":
		params["exec"] = "cmd /C type"
	default:
		params["exec"] = "cat"
	}
	cmd, err := NewCommandProc(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(cmdmsg))
	assert.FailOnErr(t, err)
	assert.NoErr(t, cmd.Process(msg))
	assert.Equal(t, msg.Content["cmdfield"], "New York")
}

func TestCommandProcProcessBroken(t *testing.T) {
	params := map[string]string{
		"read.from": "CITY",
		"store.in":  "cmdfield",
		"exec":      "invalidcmd",
	}
	cmd, err := NewCommandProc(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(cmdmsg))
	assert.FailOnErr(t, err)
	assert.Err(t, cmd.Process(msg))
}
