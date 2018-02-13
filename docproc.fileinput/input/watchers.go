package input

import (
	"fmt"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/service"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// CheckInterval represents the default interval for checking the file
	// system for new files.
	CheckInterval = time.Duration(10 * time.Second)
)

var (
	watchers = make(map[string]FileTransformerBuilder)
	wmu      = sync.Mutex{}
)

// MetadataProc is a simple post processor for the input message transformers.
// It adds 'format' and 'source' metadata information to the queue.Message
// objects created by the Transformer implementations.
type MetadataProc struct {
	Source string
	Format string
}

// Name gets the name of the MetadataProc.
func (imp *MetadataProc) Name() string {
	return "MetadataProc"
}

// Process processes a queue.Message and sets its metadata information
// 'filetype' and 'source' to the configured values of the MetadataProc.
func (imp *MetadataProc) Process(msg *queue.Message) error {
	msg.Metadata[queue.MetaFormat] = imp.Format
	msg.Metadata[queue.MetaSource] = imp.Source
	return nil
}

// FileTransformerBuilder defines a factory method for creating a FileTransformer
type FileTransformerBuilder func(params map[string]string) (FileTransformer, error)

// Register associates the passed in name with a specific Watcher.
func Register(name string, builder FileTransformerBuilder) {
	wmu.Lock()
	watchers[name] = builder
	wmu.Unlock()
}

// Create creates a FileWatcher with an associated FileTransformer.
// The parameter map must contain the following entries:
//
// * "transformer": a name of an known FileTransformer (see FileTransformers())
// * "format": the file or data format to be used as metadata information
// * "folder.in": the directory to watch for new files to be processed
// * "pattern": the file pattern to use when looking for new files in "folder.in"
// * "interval": the interval in seconds to use for checking
//
func Create(wq queue.WriteQueue, params map[string]string) (*service.FileWatcher, error) {
	tfname, ok := params["transformer"]
	if !ok {
		return nil, fmt.Errorf("parameter 'transformer' missing")
	}
	format, ok := params["format"]
	if !ok {
		return nil, fmt.Errorf("parameter 'format' missing")
	}
	directory, ok := params["folder.in"]
	if !ok {
		return nil, fmt.Errorf("parameter 'folder.in' missing")
	}
	pattern, ok := params["pattern"]
	if !ok {
		return nil, fmt.Errorf("parameter 'pattern' missing")
	}
	interval := CheckInterval
	iv, ok := params["interval"]
	if ok {
		vv, err := strconv.Atoi(iv)
		if err != nil {
			return nil, fmt.Errorf("parameter 'interval' is not an integer")
		}
		if vv <= 0 {
			return nil, fmt.Errorf("parameter 'interval' must be greater than 0")
		}
		interval = time.Duration(vv) * time.Second
	}

	for _, s := range []string{suffixProcess, suffixDone, suffixFailed} {
		if strings.Contains(pattern, s) {
			return nil, fmt.Errorf("pattern must not contain '%s'", s)
		}
	}

	wmu.Lock()
	builder, ok := watchers[tfname]
	wmu.Unlock()
	if !ok {
		return nil, fmt.Errorf("invalid transfomer '%s'", tfname)
	}
	transformer, err := builder(params)
	if err != nil {
		return nil, err
	}
	handler := NewFileHandler(wq, transformer)
	handler.AddProcessor(&MetadataProc{Source: "file", Format: format})
	return service.NewFileWatcher(directory, pattern, interval, handler)
}

// FileTransfomers returns the available FileTransformer builder types.
func FileTransfomers() []string {
	wmu.Lock()
	defer wmu.Unlock()
	keys := make([]string, len(watchers))
	i := 0
	for k := range watchers {
		keys[i] = k
		i++
	}
	return keys
}
