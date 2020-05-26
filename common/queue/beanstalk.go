// +build beanstalk

package queue

import (
	"errors"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/marcusva/docproc/common/log"
)

func init() {
	Register("beanstalk", newBeanstalkRQ, newBeanstalkWQ)
}

// beanstalkWQ provides writable access to a beanstalk queue.
type beanstalkWQ struct {
	tube *beanstalk.Tube
	host string
}

// newBeanstalkWQ creates a new beanstalk producer
func newBeanstalkWQ(params map[string]string) (WriteQueue, error) {
	host, ok := params["host"]
	if !ok {
		return nil, errors.New("'host' parameter missing")
	}
	topic, ok := params["topic"]
	if !ok {
		return nil, errors.New("'topic' parameter missing")
	}
	return &beanstalkWQ{
		tube: &beanstalk.Tube{
			Conn: nil,
			Name: topic,
		},
		host: host,
	}, nil
}

// Topic returns the topic being used to write messages to.
func (wq *beanstalkWQ) Topic() string {
	return wq.tube.Name
}

// Open opens the queue for writing.
func (wq *beanstalkWQ) Open() error {
	if wq.tube.Conn != nil {
		return errors.New("queue already opened")
	}
	conn, err := beanstalk.Dial("tcp", wq.host)
	if err != nil {
		return err
	}
	wq.tube.Conn = conn
	return nil
}

func (wq *beanstalkWQ) IsOpen() bool {
	if wq.tube.Conn == nil {
		return false
	}
	if _, err := wq.tube.Stats(); err != nil {
		log.Warningf("queue not reachable anymore, resetting it; error: %v", err)
		wq.Close()
		return false
	}
	return true
}

// Close closes the queue.
func (wq *beanstalkWQ) Close() error {
	if wq.tube.Conn == nil {
		return nil
	}
	err := wq.tube.Conn.Close()
	if err != nil {
		log.Warningf("error on closing the queue: %v", err)
	}
	wq.tube.Conn = nil
	return err
}

// Publish sends a message to the queue using the configured topic.
func (wq *beanstalkWQ) Publish(msg *Message) error {
	if wq.tube.Conn == nil {
		return errors.New("queue not open")
	}
	buf, err := msg.ToJSON()
	if err != nil {
		return err
	}
	_, err = wq.tube.Put(buf, 1, 0, 1*time.Hour)
	return err
}

// beanstalkRQ provides readable access to a beanstalk queue.
type beanstalkRQ struct {
	tubeset  *beanstalk.TubeSet
	host     string
	topic    string
	consumer Consumer
}

// newBeanstalkRQ creates a new NSQ consumer
func newBeanstalkRQ(params map[string]string) (ReadQueue, error) {
	host, ok := params["host"]
	if !ok {
		return nil, errors.New("'host' parameter missing")
	}
	topic, ok := params["topic"]
	if !ok {
		return nil, errors.New("'topic' parameter missing")
	}
	return &beanstalkRQ{
		tubeset:  beanstalk.NewTubeSet(nil, topic),
		host:     host,
		topic:    topic,
		consumer: nil,
	}, nil
}

func (rq *beanstalkRQ) watch() {
	for rq.tubeset.Conn != nil {
		id, data, err := rq.tubeset.Reserve(1 * time.Second)
		if err != nil {
			log.Infof("error on querying the queue: %v\n", err)
			continue
		}
		msg, err := MsgFromJSON(data)
		if err != nil {
			log.Errorf("error on converting the queue message: %v\n", err)
			continue
		}
		if err := rq.consumer.Consume(msg); err != nil {
			log.Errorf("error on processing: %v\n", err)
			continue
		} else {
			if err := rq.tubeset.Conn.Delete(id); err != nil {
				log.Errorf("error on deleting the job from the queue: %v\n", err)
			}
		}
	}
}

// Open opens the beanstalk queue for reading and receicing messages.
func (rq *beanstalkRQ) Open(consumer Consumer) error {
	if rq.tubeset.Conn != nil {
		return errors.New("queue already opened")
	}
	conn, err := beanstalk.Dial("tcp", rq.host)
	if err != nil {
		return err
	}
	rq.tubeset.Conn = conn
	go rq.watch()
	return nil
}

// Close closes the queue.
func (rq *beanstalkRQ) Close() error {
	if rq.tubeset.Conn == nil {
		return nil
	}
	err := rq.tubeset.Conn.Close()
	if err != nil {
		log.Warningf("error on closing the queue: %v", err)
	}
	rq.tubeset.Conn = nil
	return err
}

// Topic returns the topic being used to read messages from.
func (rq *beanstalkRQ) Topic() string {
	return rq.topic
}
