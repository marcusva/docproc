package input

import (
	"fmt"
	"github.com/marcusva/docproc/common/queue"
	"sync"
)

var (
	handlers = make(map[string]WebHandlerBuilder)
	hmu      = sync.Mutex{}
)

// WebHandlerBuilder defines a factory method for creating a WebHandler
type WebHandlerBuilder func(wq queue.WriteQueue, params map[string]string) (WebHandler, error)

// Register associates the passed in name with a specific WebHandler.
func Register(name string, builder WebHandlerBuilder) {
	hmu.Lock()
	handlers[name] = builder
	hmu.Unlock()
}

// Create creates a WebHandler using the provided configuration information
// supplied via params.
func Create(wq queue.WriteQueue, params map[string]string) (WebHandler, error) {
	name, ok := params["type"]
	if !ok {
		return nil, fmt.Errorf("parameter 'type' missing")
	}
	hmu.Lock()
	builder, ok := handlers[name]
	hmu.Unlock()
	if !ok {
		return nil, fmt.Errorf("invalid processor '%s'", name)
	}
	return builder(wq, params)
}

// WebHandlers returns the available handlers.
func WebHandlers() []string {
	hmu.Lock()
	defer hmu.Unlock()
	keys := make([]string, len(handlers))
	i := 0
	for k := range handlers {
		keys[i] = k
		i++
	}
	return keys
}
