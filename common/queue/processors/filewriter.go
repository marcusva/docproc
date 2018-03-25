package processors

import (
	"encoding/json"
	"fmt"
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

// FileWriter writes a portion of a queue.Message to a file.
type FileWriter struct {
	identifier string
	filename   string
	path       string
	rules      []rules.Rule
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
	buf, ok := msg.Content[fw.identifier]
	if !ok {
		return fmt.Errorf("message '%s' misses identifier '%s'", msg.Metadata[queue.MetaID], fw.identifier)
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
	var bytebuf []byte
	switch buf.(type) {
	case []byte:
		bytebuf = buf.([]byte)
	case string:
		bytebuf = []byte(buf.(string))
	default:
		log.Infof("content '%s' is not a string or byte buffer, using standard conversion", fw.identifier)
		bytebuf = []byte(fmt.Sprintf("%v", buf))
	}
	fpath := filepath.Join(fw.path, filename)
	log.Debugf("writing html to '%s'", fpath)
	return ioutil.WriteFile(fpath, bytebuf, os.FileMode(0644))
}

// NewFileWriter creates a FileWriter
func NewFileWriter(params map[string]string) (queue.Processor, error) {
	identifier, ok := params["identifier"]
	if !ok {
		return nil, fmt.Errorf("parameter 'identifier' missing")
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
		identifier: identifier,
		filename:   filename,
		rules:      rules,
		path:       directory,
	}, nil
}
