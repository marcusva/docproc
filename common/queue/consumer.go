package queue

import (
	"github.com/marcusva/docproc/common/log"
	"sync"
)

// Consumer is a simple reader function for messages
type Consumer interface {
	Consume(*Message) error
}

// ProcConsumer combines the Consumer for reading messages with the possibility
// to use Processor implementations on Consume()
type ProcConsumer interface {
	Consumer
	Add(Processor)
}

type SimpleConsumer struct {
	// Processors is the list of Processor instances to be executed.
	Processors []Processor
	// mu represents a mutex for thread-safety on running Add() or Consume()
	mu sync.Mutex
}

// NewSimpleConsumer creates a new ProcConsumer
func NewSimpleConsumer() *SimpleConsumer {
	return &SimpleConsumer{
		Processors: make([]Processor, 0),
		mu:         sync.Mutex{},
	}
}

// Add adds a Processor to the SimpleConsumer's processing chain.
// It will not check, if there is already a Processor of the same type.
func (pc *SimpleConsumer) Add(pp Processor) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.Processors = append(pc.Processors, pp)
}

// Consume consumes a Message and executes the associated processors on it.
// If a processor fails, the execution will stop immediately and an error is
// returned.
func (pc *SimpleConsumer) Consume(msg *Message) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	for _, pp := range pc.Processors {
		if err := pp.Process(msg); err != nil {
			log.Errorf("error on processing in '%s': %v", pp.Name(), err)
			return err
		}
	}
	return nil
}
