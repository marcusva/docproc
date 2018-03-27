package processors

import (
	"encoding/json"
	"fmt"
	"github.com/marcusva/docproc/common/data"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/path"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/rules"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	fwName = "FileWriter"
)

func init() {
	Register(fwName, NewFileWriter)
}

// FileWriter writes a specific content entry of a queue.Message into a file.
// The message must contain the content as well as a file name to use.
type FileWriter struct {
	// readFrom denotes the entry of the message content to write.
	// data := message.Content[readFrom]
	readFrom string

	// filename denotes the entry of the filename within the message content.
	// fname := message.Content[filename]
	filename string

	// directry is the directory to write the file(s) into.
	directory string

	// rules contains the rules to be executed on processing a message.
	// The FileWriter only writes files for messages, which satisfy the rules.
	rules []rules.Rule
}

// Name returns the name to be used in configuration files.
func (fw *FileWriter) Name() string {
	return fwName
}

// Process processes the message and writes certain content of it to the
// configured directory. The message needs to have two key-value pairs, one
// containing the content to write to the file and another one containing the
// filename to use.
func (fw *FileWriter) Process(msg *queue.Message) error {
	buf, ok := msg.Content[fw.readFrom]
	if !ok {
		return fmt.Errorf("message '%s' misses identifier '%s'", msg.Metadata[queue.MetaID], fw.readFrom)
	}
	filename, ok := msg.Content[fw.filename].(string)
	if !ok {
		return fmt.Errorf("message '%s' misses filename '%s'", msg.Metadata[queue.MetaID], fw.filename)
	}

	for _, rule := range fw.rules {
		if ok, err := rule.Test(msg.Content); err != nil {
			return err
		} else if !ok {
			// TODO: this way of doing things feels wrong - better use a
			// message filter beforehand
			return fmt.Errorf("message '%s' does not satisfy the rules", msg.Metadata[queue.MetaID])
		}
	}
	bytebuf, err := data.Bytes(buf)
	if err != nil {
		return err
	}
	fpath := filepath.Join(fw.directory, filename)
	log.Debugf("writing html to '%s'", fpath)
	return ioutil.WriteFile(fpath, bytebuf, os.FileMode(0644))
}

// NewFileWriter creates a FileWriter.
// The parameter map params must contain the following entries:
//
// * "read.from": the key of the message content entry to write.
// * "filename": the key of the filename entry of the message content.
//
func NewFileWriter(params map[string]string) (queue.Processor, error) {
	identifier, ok := params["read.from"]
	if !ok {
		return nil, fmt.Errorf("parameter 'read.from' missing")
	}
	filename, ok := params["filename"]
	if !ok {
		return nil, fmt.Errorf("parameter 'filename' missing")
	}
	directory, ok := params["path"]
	if !ok {
		return nil, fmt.Errorf("parameter 'path' missing")
	}
	if ok, err := path.DirExists(directory); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("path '%s' is not a directory", directory)
	}
	if ok, err := path.Writable(directory); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("path '%s' is not a writable", directory)
	}

	rulefile, ok := params["rules"]
	if !ok {
		return nil, fmt.Errorf("parameter 'rules' missing")
	}
	data, err := ioutil.ReadFile(rulefile)
	if err != nil {
		return nil, err
	}
	var rules []rules.Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}
	return &FileWriter{
		readFrom:  identifier,
		filename:  filename,
		rules:     rules,
		directory: directory,
	}, nil
}
