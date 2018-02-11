package queue

import (
	"fmt"
	"github.com/marcusva/docproc/common/log"
)

// Writer receives data from a queue, runs them through one or more
// Processor instances and places the processed result on the writing queue.
// If an error queue is set on the Writer, messages, which fail on processing
// will be stored within that queue.
type Writer struct {
	Queue     WriteQueue
	ErrQueue  WriteQueue
	ProcChain []Processor
}

// NewWriter creates a new Writer
func NewWriter(queue WriteQueue) *Writer {
	return &Writer{
		Queue:     queue,
		ErrQueue:  nil,
		ProcChain: make([]Processor, 0),
	}
}

// Open opens all bound queues.
func (qw *Writer) Open() error {
	if qw.Queue != nil && !qw.Queue.IsOpen() {
		if err := qw.Queue.Open(); err != nil {
			return err
		}
	}
	if qw.ErrQueue != nil && !qw.ErrQueue.IsOpen() {
		if err := qw.ErrQueue.Open(); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the bound queues of the Writer.
func (qw *Writer) Close() error {
	var err error
	if qw.Queue != nil && qw.Queue.IsOpen() {
		err = qw.Queue.Close()
	}
	if qw.ErrQueue != nil && qw.ErrQueue.IsOpen() {
		if err2 := qw.ErrQueue.Close(); err2 != nil {
			if err == nil {
				err = err2
			}
		}
	}
	return err
}

// AddProcessor adds a Processor to the Writer's processing chain.
// It will not check, if there is already a Processor of the same type.
func (qw *Writer) AddProcessor(pp Processor) {
	qw.ProcChain = append(qw.ProcChain, pp)
}

// Consume runs the passed in message through the associated ProcChain and
// places the processed result on the bound WriteQueue.
//
// If an error occurs, processing will stop and nothing will be placed on the
// bound WriteQueue. The message will be published on the ErrQueue of the Writer,
// if set.
func (qw *Writer) Consume(msg *Message) error {
	log.Infof("Received message '%v'", msg.Metadata[MetaID])
	for _, proc := range qw.ProcChain {
		log.Debugf("Executing processor %s...", proc.Name())
		if err := proc.Process(msg); err != nil {
			if qw.ErrQueue != nil {
				// TODO: handle error
				if err2 := qw.ErrQueue.Publish(msg); err2 != nil {
					log.Errorf("could not pass the message to the error queue: %v", err2)
				}
			}
			return fmt.Errorf("processing failed in '%s': %v", proc.Name(), err)
		}
	}
	err := qw.Queue.Publish(msg)
	if err != nil {
		log.Errorf("could not pass the message to the output queue: %v", err)
	}
	return err
}
