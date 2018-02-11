package queue

import (
	"fmt"
	"sync"
)

// Consumer is a simple reader function for messages
type Consumer interface {
	Consume(*Message) error
}

// WriteQueue provides write access to a queue
type WriteQueue interface {
	Open() error
	IsOpen() bool
	Close() error
	Topic() string
	Publish(*Message) error
}

// ReadQueue provides read access to a queue. Implementations should ensure that
// messages will stay in the queue (or get republished), if the bound consumer
// fails to process them successfully.
type ReadQueue interface {
	Open(Consumer) error
	Close() error
	Topic() string
}

// WriteQueueBuilder defines a factory method for creating a WriteQueue
type WriteQueueBuilder func(params map[string]string) (WriteQueue, error)

// ReadQueueBuilder defines a factory method for creating a ReadQueue
type ReadQueueBuilder func(params map[string]string) (ReadQueue, error)

var (
	readqueues  = make(map[string]ReadQueueBuilder)
	writequeues = make(map[string]WriteQueueBuilder)
	qmu         = sync.Mutex{}
)

// Register associates the passed in name with specific builders for ReadQueue
// and WriteQueue instances.
func Register(name string, rq ReadQueueBuilder, wq WriteQueueBuilder) {
	RegisterRQ(name, rq)
	RegisterWQ(name, wq)
}

// RegisterRQ associates the passed in name with a specific builder for
// ReadQueue instance.
func RegisterRQ(name string, rq ReadQueueBuilder) {
	qmu.Lock()
	readqueues[name] = rq
	qmu.Unlock()
}

// RegisterWQ associates the passed in name with a specific builder for
// WriteQueue instance.
func RegisterWQ(name string, wq WriteQueueBuilder) {
	qmu.Lock()
	writequeues[name] = wq
	qmu.Unlock()
}

// CreateRQ creates a ReadQueue
func CreateRQ(queue string, params map[string]string) (ReadQueue, error) {
	qmu.Lock()
	builder, ok := readqueues[queue]
	qmu.Unlock()
	if !ok {
		return nil, fmt.Errorf("Invalid readable queue type '%s'", queue)
	}
	return builder(params)
}

// CreateWQ creates a WriteQueue
func CreateWQ(queue string, params map[string]string) (WriteQueue, error) {
	qmu.Lock()
	builder, ok := writequeues[queue]
	qmu.Unlock()
	if !ok {
		return nil, fmt.Errorf("Invalid writeable queue type '%s'", queue)
	}
	return builder(params)
}

// ReadTypes returns the available ReadQueue builder types.
func ReadTypes() []string {
	qmu.Lock()
	defer qmu.Unlock()
	keys := make([]string, len(readqueues))
	i := 0
	for k := range readqueues {
		keys[i] = k
		i++
	}
	return keys
}

// WriteTypes returns the available WriteQueue builder types.
func WriteTypes() []string {
	qmu.Lock()
	defer qmu.Unlock()
	keys := make([]string, len(writequeues))
	i := 0
	for k := range writequeues {
		keys[i] = k
		i++
	}
	return keys
}
