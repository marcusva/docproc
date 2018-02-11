package renderers

import (
	"fmt"
	"github.com/marcusva/docproc/common/queue"
	"sync"
)

var (
	renderers = make(map[string]Builder)
	emu       = sync.Mutex{}
)

var (
	// ContentType represents the rendered result's content type.
	ContentType = "mime-type"
)

// Builder defines a factory method for creating an Renderer
type Builder func(params map[string]string) (queue.Processor, error)

// Register associates the passed in name with a specific Builder.
func Register(name string, builder Builder) {
	emu.Lock()
	renderers[name] = builder
	emu.Unlock()
}

// Create creates a Renderer based on on the 'renderer' parameter
// TODO: be more elaborate.
func Create(params map[string]string) (queue.Processor, error) {
	name, ok := params["renderer"]
	if !ok {
		return nil, fmt.Errorf("parameter 'renderer' missing")
	}
	emu.Lock()
	builder, ok := renderers[name]
	emu.Unlock()
	if !ok {
		return nil, fmt.Errorf("invalid renderer '%s'", name)
	}
	return builder(params)
}

// Renderers returns the available Renderer builder types.
func Renderers() []string {
	emu.Lock()
	defer emu.Unlock()
	keys := make([]string, len(renderers))
	i := 0
	for k := range renderers {
		keys[i] = k
		i++
	}
	return keys
}
