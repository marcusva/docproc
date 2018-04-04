// +build nsq

package queue

import (
	"errors"
	"github.com/marcusva/docproc/common/log"
	"github.com/nsqio/go-nsq"
	//"runtime"
)

func init() {
	Register("nsq", newNsqRQ, newNsqWQ)
}

func mapLogLevel(lvl log.Level) nsq.LogLevel {
	switch lvl {
	case log.LevelDebug:
		return nsq.LogLevelDebug
	case log.LevelInfo, log.LevelNotice:
		return nsq.LogLevelInfo
	case log.LevelWarning:
		return nsq.LogLevelWarning
	case log.LevelError, log.LevelAlert, log.LevelCritical, log.LevelEmergency:
		return nsq.LogLevelError
	default:
		return nsq.LogLevelInfo
	}
}

// nsqWQ provides writable access to a NSQ queue.
type nsqWQ struct {
	nsqproducer *nsq.Producer
	topic       string
	host        string
}

// newNsqWQ creates a new NSQ producer
func newNsqWQ(params map[string]string) (WriteQueue, error) {
	host, ok := params["host"]
	if !ok {
		return nil, errors.New("'host' parameter missing")
	}
	topic, ok := params["topic"]
	if !ok {
		return nil, errors.New("'topic' parameter missing")
	}

	return &nsqWQ{
		nsqproducer: nil,
		topic:       topic,
		host:        host,
	}, nil
}

// Topic returns the topic being used to write messages to.
func (wq *nsqWQ) Topic() string {
	return wq.topic
}

// IsOpen checks, if the NSQ queue is available and open for writing.
func (wq *nsqWQ) IsOpen() bool {
	if wq.nsqproducer == nil {
		return false
	}
	if err := wq.nsqproducer.Ping(); err != nil {
		log.Errorf("queue not reachable anymore, resetting it; error: %v", err)
		wq.nsqproducer = nil
		return false
	}
	return true
}

// Open opens the queue for writing.
func (wq *nsqWQ) Open() error {
	if wq.nsqproducer != nil {
		return errors.New("queue already opened")
	}

	nsqprod, err := nsq.NewProducer(wq.host, nsq.NewConfig())
	if err != nil {
		return err
	}
	nsqprod.SetLogger(log.Logger(), mapLogLevel(log.CurrentLevel()))

	err = nsqprod.Ping()
	if err == nil {
		wq.nsqproducer = nsqprod
	}
	return err
}

// Close closes the queue.
func (wq *nsqWQ) Close() error {
	if wq.nsqproducer == nil {
		return nil
	}
	wq.nsqproducer.Stop()
	wq.nsqproducer = nil
	return nil
}

// Publish sends a message to the queue using the configured topic.
func (wq *nsqWQ) Publish(msg *Message) error {
	if wq.nsqproducer == nil {
		return errors.New("queue not open")
	}
	data, err := msg.ToJSON()
	if err != nil {
		return err
	}
	err = wq.nsqproducer.Publish(wq.topic, data)
	if err == nsq.ErrNotConnected {
		log.Errorf("queue not reachable, error: %v", err)
		wq.nsqproducer = nil
	}
	return err
}

// nsqRQ provides readable access to a NSQ queue.
type nsqRQ struct {
	nsqconsumer *nsq.Consumer
	topic       string
	host        string
	consumer    Consumer
}

// newNsqRQ creates a new NSQ consumer
func newNsqRQ(params map[string]string) (ReadQueue, error) {
	host, ok := params["host"]
	if !ok {
		return nil, errors.New("'host' parameter missing")
	}
	topic, ok := params["topic"]
	if !ok {
		return nil, errors.New("'topic' parameter missing")
	}
	return &nsqRQ{
		nsqconsumer: nil,
		consumer:    nil,
		topic:       topic,
		host:        host,
	}, nil
}

// Open opens the NSQ queue for reading and receicing messages.
func (rq *nsqRQ) Open(consumer Consumer) error {
	if rq.nsqconsumer != nil {
		return errors.New("queue already opened")
	}
	nsqconsumer, err := nsq.NewConsumer(rq.topic, "docproc", nsq.NewConfig())
	if err != nil {
		return err
	}
	nsqconsumer.SetLogger(log.Logger(), mapLogLevel(log.CurrentLevel()))

	rq.consumer = consumer
	rq.nsqconsumer = nsqconsumer

	handler := nsq.HandlerFunc(func(m *nsq.Message) error {
		msg, err := MsgFromJSON(m.Body)
		if err != nil {
			log.Errorf("could not convert message: %v", err)
			return err
		}
		if err = rq.consumer.Consume(msg); err != nil {
			log.Errorf("error on processing: %v", err)
			return err
		}
		return nil
	})
	//log.Debugf("Setting up %d concurrent consumers", runtime.GOMAXPROCS(0))
	//rq.nsqconsumer.AddConcurrentHandlers(handler, runtime.GOMAXPROCS(0))
	rq.nsqconsumer.AddHandler(handler)
	return rq.nsqconsumer.ConnectToNSQLookupd(rq.host)
}

// Close closes the NSQ queue.
func (rq *nsqRQ) Close() error {
	if rq.nsqconsumer == nil {
		return nil
	}
	rq.nsqconsumer.Stop()
	err := rq.nsqconsumer.DisconnectFromNSQLookupd(rq.host)
	rq.nsqconsumer = nil
	rq.consumer = nil
	return err
}

// Topic returns the topic being used to read messages from.
func (rq *nsqRQ) Topic() string {
	return rq.topic
}
