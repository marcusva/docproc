package queue

import (
	"github.com/marcusva/docproc/common/log"
	"sync"
)

// Writer receives data from a queue, runs them through one or more
// Processor instances and places the processed result on the writing queue.
// If an error queue is set on the Writer, messages, which fail on processing
// will be stored within that queue.
type Writer struct {
	processors []Processor
	queue      WriteQueue
	errQueue   WriteQueue
	mu         sync.RWMutex
}

// NewWriter creates a new Writer
func NewWriter(queue, errqueue WriteQueue) *Writer {
	return &Writer{
		processors: make([]Processor, 0),
		queue:      queue,
		errQueue:   errqueue,
		mu:         sync.RWMutex{},
	}
}

// Open opens all bound queues.
func (qw *Writer) Open() error {
	if qw.queue != nil && !qw.queue.IsOpen() {
		if err := qw.queue.Open(); err != nil {
			return err
		}
	}
	if qw.errQueue != nil && !qw.errQueue.IsOpen() {
		if err := qw.errQueue.Open(); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the bound queues of the Writer.
func (qw *Writer) Close() error {
	var err error
	qw.mu.Lock()
	defer qw.mu.Unlock()
	if qw.queue != nil && qw.queue.IsOpen() {
		err = qw.queue.Close()
	}
	if qw.errQueue != nil && qw.errQueue.IsOpen() {
		err2 := qw.errQueue.Close()
		if err2 != nil && err == nil {
			return err2
		}
	}
	return err
}

// Add adds a Processor to the Writer's processing chain.
// It will not check, if there is already a Processor of the same type.
func (qw *Writer) Add(pp Processor) {
	qw.mu.Lock()
	defer qw.mu.Unlock()
	qw.processors = append(qw.processors, pp)
}

// Consume consumes a Message and executes the associated processors on it.
//
// If an error occurs, processing will stop and nothing will be placed on the
// bound WriteQueue. The message will be published on the ErrQueue of the
// Writer, if set.
func (qw *Writer) Consume(msg *Message) error {
	log.Infof("Received message '%v'", msg.Metadata[MetaID])
	qw.mu.RLock()
	defer qw.mu.RUnlock()
	for _, pp := range qw.processors {
		if err := pp.Process(msg); err != nil {
			if qw.errQueue == nil {
				return err
			}
			if err2 := qw.errQueue.Publish(msg); err2 != nil {
				log.Errorf("could not pass the message to the error queue: %v", err2)
			}
			return err
		}
	}
	if qw.queue == nil {
		return nil
	}
	err := qw.queue.Publish(msg)
	if err != nil {
		log.Errorf("could not pass the message to the output queue: %v", err)
	}
	return err
}
