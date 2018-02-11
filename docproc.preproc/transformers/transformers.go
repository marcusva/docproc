package transformers

import (
	"fmt"
	"github.com/marcusva/docproc/common/queue"
	"sync"
)

var (
	transformers = make(map[string]Builder)
	emu          = sync.Mutex{}
)

// Builder defines a factory method for creating a Transformer
type Builder func(params map[string]string) (queue.Processor, error)

// Register associates the passed in name with a specific queue.Processor.
func Register(name string, builder Builder) {
	emu.Lock()
	transformers[name] = builder
	emu.Unlock()
}

// Create creates a Transformer based on on the 'transformer' parameter
// TODO: be more elaborate.
func Create(params map[string]string) (queue.Processor, error) {
	name, ok := params["transformer"]
	if !ok {
		return nil, fmt.Errorf("parameter 'transformer' missing")
	}
	emu.Lock()
	builder, ok := transformers[name]
	emu.Unlock()
	if !ok {
		return nil, fmt.Errorf("invalid transformer '%s'", name)
	}
	return builder(params)
}

// Transformers returns the available transformer builder types.
func Transformers() []string {
	emu.Lock()
	defer emu.Unlock()
	keys := make([]string, len(transformers))
	i := 0
	for k := range transformers {
		keys[i] = k
		i++
	}
	return keys
}
