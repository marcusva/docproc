// +build nats

package queue

import (
	"errors"
	"github.com/marcusva/docproc/common/log"
	"github.com/nats-io/go-nats"
)

func init() {
	Register("nats", newNatsRQ, newNatsWQ)
}

//
type natsQueue struct {
	conn  *nats.Conn
	host  string
	topic string
}

type natsRQ struct {
	natsQueue
	sub *nats.Subscription
}

type natsWQ struct {
	natsQueue
}

func newNatsRQ(params map[string]string) (ReadQueue, error) {
	host, ok := params["host"]
	if !ok {
		return nil, errors.New("'host' parameter missing")
	}
	topic, ok := params["topic"]
	if !ok {
		return nil, errors.New("'topic' parameter missing")
	}
	return &natsRQ{
		natsQueue: natsQueue{
			conn:  nil,
			host:  host,
			topic: topic,
		},
		sub: nil,
	}, nil
}

func newNatsWQ(params map[string]string) (WriteQueue, error) {
	host, ok := params["host"]
	if !ok {
		return nil, errors.New("'host' parameter missing")
	}
	topic, ok := params["topic"]
	if !ok {
		return nil, errors.New("'topic' parameter missing")
	}
	return &natsWQ{natsQueue{
		conn:  nil,
		host:  host,
		topic: topic,
	}}, nil
}

func (q *natsQueue) close() error {
	if q.conn == nil {
		return nil
	}
	q.conn.Close()
	q.conn = nil
	return nil
}

func (q *natsQueue) IsOpen() bool {
	if q.conn == nil || q.conn.IsClosed() {
		return false
	}
	return q.conn.IsConnected()
}

func (q *natsQueue) open() error {
	if q.conn != nil {
		if !q.conn.IsClosed() {
			return errors.New("queue already opened")
		}
	}
	conn, err := nats.Connect(q.host)
	if err != nil {
		return err
	}
	q.conn = conn
	return nil
}

func (q *natsQueue) Topic() string {
	return q.topic
}

func (q *natsRQ) Close() error {
	if q.sub != nil {
		q.sub.Unsubscribe()
		q.sub = nil
	}
	return q.close()
}

func (q *natsRQ) Open(c Consumer) error {
	if err := q.open(); err != nil {
		return err
	}
	sub, err := q.conn.Subscribe(q.topic, func(msg *nats.Msg) {
		m, err2 := MsgFromJSON(msg.Data)
		if err2 != nil {
			log.Errorf("could not convert message: %v", err2)
		}
		if err2 = c.Consume(m); err2 != nil {
			log.Errorf("error on processing: %v", err2)
		}
	})
	if err != nil {
		return err
	}
	q.sub = sub
	return nil
}

func (q *natsWQ) Close() error {
	return q.close()
}

func (q *natsWQ) Open() error {
	if err := q.open(); err != nil {
		return err
	}
	return nil
}

func (q *natsWQ) Publish(msg *Message) error {
	if !q.IsOpen() {
		if err := q.Open(); err != nil {
			return err
		}
	}
	buf, err := msg.ToJSON()
	if err != nil {
		return err
	}
	err = q.conn.Publish(q.topic, buf)
	if err != nil {
		return err
	}
	return nil
}
