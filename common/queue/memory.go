package queue

import (
	"errors"
	"github.com/marcusva/docproc/common/log"
	"runtime"
	"sync"
)

const (
	// maxTopicSize is the maximum of allowed messages per topic
	maxTopicSize = 500
)

var (
	topics = make(map[string]*topic)
	memmu  = sync.Mutex{}
)

func init() {
	Register("memory", newMemRQ, newMemWQ)
}

// topic is a simple name/channel storage
type topic struct {
	// Name is the name of the topic
	Name string

	// Incoming is a r/w channel containing the messages of the topic
	Incoming chan []byte
}

// newTopic creates a new topic
func newTopic(name string) *topic {
	return &topic{
		Name:     name,
		Incoming: make(chan []byte, maxTopicSize),
	}
}

// memWQ is an in-memory writable queue.
type memWQ struct {
	open  bool
	topic *topic
}

// Topic gets the topic the memWQ is associated with.
func (wq *memWQ) Topic() string {
	return wq.topic.Name
}

// Close closes the memWQ, but leaves the underlying topic unmodified.
func (wq *memWQ) Close() error {
	wq.open = false
	return nil
}

// IsOpen checks, if the memWQ can publish messages to its associated topic.
func (wq *memWQ) IsOpen() bool {
	return wq.open
}

// Open opens the memWQ for publishing.
func (wq *memWQ) Open() error {
	wq.open = true
	return nil
}

// Publish publishs a Message to the memWQ's topic.
func (wq *memWQ) Publish(msg *Message) error {
	if !wq.open {
		return errors.New("queue is not open")
	}
	data, err := msg.ToJSON()
	if err != nil {
		return err
	}
	select {
	case wq.topic.Incoming <- data:
		return nil
	default:
		return errors.New("queue is full")
	}
}

// newMemWQ creates a new in-memory writable queue.
func newMemWQ(params map[string]string) (WriteQueue, error) {
	name, ok := params["topic"]
	if !ok {
		return nil, errors.New("parameter 'topic' is missing")
	}
	memmu.Lock()
	tp, ok := topics[name]
	if !ok {
		tp = newTopic(name)
		topics[name] = tp
	}
	memmu.Unlock()
	return &memWQ{
		open:  false,
		topic: tp,
	}, nil
}

// memRQ is an in-memory queue that consumes messages for a specific topic.
type memRQ struct {
	open  bool
	topic *topic
	stop  chan int
}

// Topic gets the topic the memRQ is associated with.
func (rq *memRQ) Topic() string {
	return rq.topic.Name
}

// Close closes the memRQ, causing its consumer to stop receiving messages.
func (rq *memRQ) Close() error {
	if !rq.open {
		return nil
	}
	rq.stop <- 1
	close(rq.stop)
	rq.open = false
	return nil
}

// Open opens the memRQ, causing the given Consumer to receive messages for the
// memRQ's topic. The Consumer may be run concurrently, depending on the value
// of runtime.GOMAXPROCS.
func (rq *memRQ) Open(c Consumer) error {
	if rq.open {
		return errors.New("queue already open")
	}
	rq.stop = make(chan int)
	handler := func() {
		for {
			select {
			case _, ok := <-rq.stop:
				if !ok {
					return
				}
			default:
			}

			buf, ok := <-rq.topic.Incoming
			if !ok {
				return
			}
			msg, err := MsgFromJSON(buf)
			if err != nil {
				log.Errorf("could not convert message: %v", err)
			}
			if err = c.Consume(msg); err != nil {
				log.Errorf("error on processing: %v", err)
			}
		}
	}

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go handler()
	}
	rq.open = true
	return nil
}

// newMemRQ creates a new in-memory readable queue.
func newMemRQ(params map[string]string) (ReadQueue, error) {
	name, ok := params["topic"]
	if !ok {
		return nil, errors.New("parameter 'topic' is missing")
	}
	memmu.Lock()
	tp, ok := topics[name]
	if !ok {
		tp = newTopic(name)
		topics[name] = tp
	}
	memmu.Unlock()
	return &memRQ{
		open:  false,
		topic: tp,
		stop:  nil,
	}, nil
}
