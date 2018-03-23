package queue

import (
	"errors"
	"github.com/marcusva/docproc/common/log"
	"sync"
	"time"
)

var (
	memblocks = make(map[string][][]byte)
	memmu     = sync.Mutex{}
)

func init() {
	Register("memory", newMemRQ, newMemWQ)
}

// MemWQ represents a writable queue.
type memWQ struct {
	topic string
	open  bool
}

// newMemWQ creates a new in-memory queue to write messages to.
func newMemWQ(params map[string]string) (WriteQueue, error) {
	topic, ok := params["topic"]
	if !ok {
		return nil, errors.New("'topic' parameter found")
	}
	return &memWQ{
		topic: topic,
		open:  false,
	}, nil
}

// Open opens the in-memory queue.
func (wq *memWQ) Open() error {
	memmu.Lock()
	defer memmu.Unlock()
	wq.open = true
	if memblocks[wq.topic] == nil {
		memblocks[wq.topic] = make([][]byte, 0)
	}
	return nil
}

// IsOpen checks, if the in-memory queue is opened.
func (wq *memWQ) IsOpen() bool {
	return wq.open
}

// Close closes the in-memory queue.
func (wq *memWQ) Close() error {
	memmu.Lock()
	wq.open = false
	memmu.Unlock()
	return nil
}

// Publish writes a message to the queue.
func (wq *memWQ) Publish(msg *Message) error {
	memmu.Lock()
	defer memmu.Unlock()
	if wq.open == false {
		return errors.New("queue not open")
	}

	data, err := msg.ToJSON()
	if err != nil {
		return err
	}
	memblocks[wq.topic] = append(memblocks[wq.topic], data)
	return nil
}

// Topic returns the topic being used to write messages to.
func (wq *memWQ) Topic() string {
	return wq.topic
}

// memRQ provides an in-memory readable queue
type memRQ struct {
	topic    string
	consumer Consumer
	running  bool
}

// newMemRQ creates a new in-memory queue to read from
func newMemRQ(params map[string]string) (ReadQueue, error) {
	topic, ok := params["topic"]
	if !ok {
		return nil, errors.New("'topic' parameter not found")
	}
	return &memRQ{
		topic:    topic,
		consumer: nil,
		running:  false,
	}, nil
}

func (rq *memRQ) consume() {
	memmu.Lock()
	defer memmu.Unlock()
	if len(memblocks[rq.topic]) > 0 {
		msg, err := MsgFromJSON(memblocks[rq.topic][0])
		if err != nil {
			log.Errorf("could not convert message: %v", err)
			return
		}
		if err := rq.consumer.Consume(msg); err != nil {
			log.Errorf("could not consume message: %v", err)
			return
		}
		memblocks[rq.topic] = memblocks[rq.topic][1:]
	}
}

func (rq *memRQ) watchTopic() {
	rq.running = true
	for rq.running {
		select {
		case <-time.After(2 * time.Second):
			rq.consume()
		}
	}
}

// Open opens the in-memory queue for reading messages.
func (rq *memRQ) Open(consumer Consumer) error {
	if rq.consumer != nil {
		return errors.New("queue already opened")
	}

	memmu.Lock()
	if memblocks[rq.topic] == nil {
		memblocks[rq.topic] = make([][]byte, 0)
	}
	memmu.Unlock()

	rq.consumer = consumer
	go rq.watchTopic()
	return nil
}

// Close closes the queue.
func (rq *memRQ) Close() error {
	rq.consumer = nil
	rq.running = false
	return nil
}

// Topic returns the topic being used to read messages from.
func (rq *memRQ) Topic() string {
	return rq.topic
}
