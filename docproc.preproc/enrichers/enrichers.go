package enrichers

import (
	"fmt"
	"github.com/marcusva/docproc/common/queue"
	"sync"
)

var (
	enrichers = make(map[string]Builder)
	emu       = sync.Mutex{}
)

// Builder defines a factory method for creating an Enricher
type Builder func(params map[string]string) (queue.Processor, error)

// Register associates the passed in name with a specific queue.Processor.
func Register(name string, builder Builder) {
	emu.Lock()
	enrichers[name] = builder
	emu.Unlock()
}

// Create creates an Enricher based on on the 'enricher' parameter
// TODO: be more elaborate.
func Create(params map[string]string) (queue.Processor, error) {
	name, ok := params["enricher"]
	if !ok {
		return nil, fmt.Errorf("parameter 'enricher' missing")
	}
	emu.Lock()
	builder, ok := enrichers[name]
	emu.Unlock()
	if !ok {
		return nil, fmt.Errorf("invalid enricher '%s'", name)
	}
	return builder(params)
}

// Enrichers returns the available Enricher builder types.
func Enrichers() []string {
	emu.Lock()
	defer emu.Unlock()
	keys := make([]string, len(enrichers))
	i := 0
	for k := range enrichers {
		keys[i] = k
		i++
	}
	return keys
}
