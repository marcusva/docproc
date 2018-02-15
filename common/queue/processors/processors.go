package processors

import (
	"fmt"
	"github.com/marcusva/docproc/common/queue"
	"sync"
)

var (
	processors = make(map[string]Builder)
	pmu        = sync.Mutex{}
)

// Builder defines a factory method for creating a queue.Processor
type Builder func(params map[string]string) (queue.Processor, error)

// Register associates the passed in name with a specific queue.Processor.
func Register(name string, builder Builder) {
	pmu.Lock()
	processors[name] = builder
	pmu.Unlock()
}

// Create creates a Processor using the provided configuration information
// supplied via params.
func Create(params map[string]string) (queue.Processor, error) {
	name, ok := params["type"]
	if !ok {
		return nil, fmt.Errorf("parameter 'type' missing")
	}
	pmu.Lock()
	builder, ok := processors[name]
	pmu.Unlock()
	if !ok {
		return nil, fmt.Errorf("invalid processor '%s'", name)
	}
	return builder(params)
}

// Types returns the available processor types.
func Types() []string {
	pmu.Lock()
	defer pmu.Unlock()
	keys := make([]string, len(processors))
	i := 0
	for k := range processors {
		keys[i] = k
		i++
	}
	return keys
}
