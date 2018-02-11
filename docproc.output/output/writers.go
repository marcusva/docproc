package output

import (
	"fmt"
	"github.com/marcusva/docproc/common/queue"
	"sync"
)

var (
	writers = make(map[string]Builder)
	emu     = sync.Mutex{}
)

// Builder defines a factory method for creating an OutputWriter
type Builder func(params map[string]string) (queue.Processor, error)

// Register associates the passed in name with a specific OutputBuilder.
func Register(name string, builder Builder) {
	emu.Lock()
	writers[name] = builder
	emu.Unlock()
}

// Create creates an OutputWriter based on on the 'writer' parameter
// TODO: be more elaborate.
func Create(params map[string]string) (queue.Processor, error) {
	name, ok := params["writer"]
	if !ok {
		return nil, fmt.Errorf("parameter 'writer' missing")
	}
	emu.Lock()
	builder, ok := writers[name]
	emu.Unlock()
	if !ok {
		return nil, fmt.Errorf("invalid writer '%s'", name)
	}
	return builder(params)
}

// Writers returns the available output writer builder types.
func Writers() []string {
	emu.Lock()
	defer emu.Unlock()
	keys := make([]string, len(writers))
	i := 0
	for k := range writers {
		keys[i] = k
		i++
	}
	return keys
}
