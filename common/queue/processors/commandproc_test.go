package processors_test

import (
	"encoding/base64"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/queue/processors"
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
	_, err := processors.NewCommandProc(nil)
	assert.Err(t, err)

	params := map[string]string{}
	_, err = processors.NewCommandProc(params)
	assert.Err(t, err)

	params["read.from"] = "CITY"
	_, err = processors.NewCommandProc(params)
	assert.Err(t, err)

	params["store.in"] = "cmdfield"
	_, err = processors.NewCommandProc(params)
	assert.Err(t, err)

	params["exec"] = "cmd"
	_, err = processors.NewCommandProc(params)
	assert.NoErr(t, err)

	params["store.base64"] = "banana"
	_, err = processors.NewCommandProc(params)
	assert.Err(t, err)

	params["store.base64"] = "t"
	_, err = processors.NewCommandProc(params)
	assert.NoErr(t, err)
}

func TestCommandProcCreate(t *testing.T) {
	params := map[string]string{
		"type":      "CommandProc",
		"read.from": "CITY",
		"store.in":  "cmdfield",
		"exec":      "cmd",
	}
	proc, err := processors.Create(params)
	assert.FailOnErr(t, err)
	assert.Equal(t, proc.Name(), "CommandProc")
}

func TestCommandProcName(t *testing.T) {
	params := map[string]string{
		"read.from": "CITY",
		"store.in":  "cmdfield",
		"exec":      "cmd",
	}
	cmd, err := processors.NewCommandProc(params)
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
	cmd, err := processors.NewCommandProc(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(cmdmsg))
	assert.FailOnErr(t, err)
	assert.NoErr(t, cmd.Process(msg))
	assert.Equal(t, msg.Content["cmdfield"], "New York")

	params["read.from"] = "unknown"
	cmd, err = processors.NewCommandProc(params)
	assert.FailOnErr(t, err)
	assert.Err(t, cmd.Process(msg))
}

func TestCommandProcBase64(t *testing.T) {
	params := map[string]string{
		"read.from":    "CITY",
		"store.in":     "cmdfield",
		"store.base64": "t",
	}
	switch runtime.GOOS {
	case "windows":
		params["exec"] = "cmd /C type"
	default:
		params["exec"] = "cat"
	}
	cmd, err := processors.NewCommandProc(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(cmdmsg))
	assert.FailOnErr(t, err)
	assert.NoErr(t, cmd.Process(msg))
	assert.NotEqual(t, msg.Content["cmdfield"], "New York")
	assert.Equal(t, msg.Content["cmdfield"], base64.StdEncoding.EncodeToString([]byte("New York")))

}

func TestCommandProcProcessBroken(t *testing.T) {
	params := map[string]string{
		"read.from": "CITY",
		"store.in":  "cmdfield",
		"exec":      "invalidcmd",
	}
	cmd, err := processors.NewCommandProc(params)
	assert.FailOnErr(t, err)

	msg, err := queue.MsgFromJSON([]byte(cmdmsg))
	assert.FailOnErr(t, err)
	assert.Err(t, cmd.Process(msg))
}
