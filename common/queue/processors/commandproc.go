package processors

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/marcusva/docproc/common/data"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	cmdName = "CommandProc"
)

func init() {
	Register(cmdName, NewCommandProc)
}

// CommandProc is a simple command processor for queue.Message objects.
type CommandProc struct {
	readFrom string
	storeIn  string
	asBase64 bool
	execArgs []string
}

// Name returns the name to be used in configuration files.
func (cmd *CommandProc) Name() string {
	return cmdName
}

// Process processes the queue.Message by passing its content as file
// to the configured command.
func (cmd *CommandProc) Process(msg *queue.Message) error {
	buf, ok := msg.Content[cmd.readFrom]
	if !ok {
		return fmt.Errorf("message '%s' misses identifier '%s'", msg.Metadata[queue.MetaID], cmd.readFrom)
	}

	file, err := ioutil.TempFile(os.TempDir(), cmd.Name())
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	bytebuf, err := data.Bytes(buf)
	if err != nil {
		return err
	}
	if _, err = file.Write(bytebuf); err != nil {
		return err
	}
	file.Sync()
	file.Close()

	app := cmd.execArgs[0]
	varargs := make([]string, len(cmd.execArgs)-1)
	copy(varargs, cmd.execArgs[1:])
	varargs = append(varargs, file.Name())

	command := exec.Command(app, varargs...)
	stderr, err := command.StderrPipe()
	if err != nil {
		log.Infof("could not connect stderr for command")
		stderr = nil
	}
	log.Debugf("Executing '%s'", strings.Join(command.Args, " "))
	output, err := command.Output()
	if err != nil {
		if stderr != nil {
			rd := bufio.NewReader(stderr)
			errbuf, newerr := rd.ReadString(0)
			if newerr == nil {
				err = fmt.Errorf(errbuf)
			}
		}
		log.Errorf("command '%s', arguments '%s' failed: %s", app, varargs, err)
		return err
	}
	if cmd.asBase64 {
		msg.Content[cmd.storeIn] = base64.StdEncoding.EncodeToString(output)
	} else {
		msg.Content[cmd.storeIn] = string(output)
	}
	return nil
}

// NewCommandProc creates a new command processor
func NewCommandProc(params map[string]string) (queue.Processor, error) {
	inputid, ok := params["read.from"]
	if !ok {
		return nil, fmt.Errorf("parameter 'read.from' missing")
	}
	outputid, ok := params["store.in"]
	if !ok {
		return nil, fmt.Errorf("parameter 'store.in' missing")
	}
	var err error
	base64 := false
	b64, ok := params["store.base64"]
	if ok {
		base64, err = strconv.ParseBool(b64)
		if err != nil {
			return nil, err
		}
	}
	args, ok := params["exec"]
	if !ok {
		return nil, fmt.Errorf("parameter 'exec' missing")
	}
	return &CommandProc{
		readFrom: inputid,
		storeIn:  outputid,
		asBase64: base64,
		execArgs: strings.Split(args, " "),
	}, nil
}
