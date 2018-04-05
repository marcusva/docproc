package queue_test

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"sync/atomic"
	"testing"
	"time"
)

type memProc struct {
	Count int32
}

func (mp *memProc) Name() string { return "memProc" }
func (mp *memProc) Consume(msg *queue.Message) error {
	atomic.AddInt32(&mp.Count, 1)
	return nil
}

func TestMemWQPublish(t *testing.T) {
	msg := queue.NewMessage(map[string]interface{}{
		"ID": 1,
	})
	params := map[string]string{"topic": "test"}
	wq, err := queue.NewMemWQ(params)
	assert.FailOnErr(t, err)
	assert.Err(t, wq.Publish(msg))
	assert.NoErr(t, wq.Open())
	assert.NoErr(t, wq.Publish(msg))
}

func TestMemRQConsume(t *testing.T) {
	messages := []*queue.Message{}
	for i := 0; i < 100; i++ {
		msg := queue.NewMessage(map[string]interface{}{"ID": i})
		messages = append(messages, msg)

	}
	params := map[string]string{"topic": "test"}
	wq, err := queue.NewMemWQ(params)
	assert.FailOnErr(t, err)
	assert.Err(t, wq.Publish(messages[0]))
	assert.NoErr(t, wq.Open())

	for _, msg := range messages {
		assert.NoErr(t, wq.Publish(msg))
	}

	consumer := &memProc{Count: 0}
	rq, err := queue.NewMemRQ(params)
	assert.FailOnErr(t, err)
	assert.NoErr(t, rq.Open(consumer))

	time.Sleep(10 * time.Millisecond)
	cnt := atomic.LoadInt32(&consumer.Count)
	assert.FailIf(t, cnt < int32(len(messages)),
		"processing took too long: %d of %d", cnt, len(messages))
}
