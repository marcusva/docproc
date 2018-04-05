package queue_test

import (
	"errors"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

type testProc struct{}

func (tp *testProc) Name() string { return "testProc" }
func (tp *testProc) Process(msg *queue.Message) error {
	if _, ok := msg.Content["works"]; !ok {
		return errors.New("failed")
	}
	msg.Content["testProc"] = true
	return nil
}

func TestNewSimpleConsumer(t *testing.T) {
	sc := queue.NewSimpleConsumer()
	assert.Equal(t, len(sc.Processors), 0)
}

func TestSimpleConsumerAdd(t *testing.T) {
	proc := &testProc{}
	sc := queue.NewSimpleConsumer()
	sc.Add(proc)
	assert.Equal(t, len(sc.Processors), 1)
	assert.Equal(t, sc.Processors[0], proc)
}

func TestSimpleConsumerConsume(t *testing.T) {
	proc := &testProc{}
	sc := queue.NewSimpleConsumer()
	sc.Add(proc)
	assert.Err(t, sc.Consume(queue.NewMessage(nil)))
	assert.Err(t, sc.Consume(queue.NewMessage(map[string]interface{}{"test": "test"})))

	msg := queue.NewMessage(map[string]interface{}{"works": true})
	assert.NoErr(t, sc.Consume(msg))
	assert.Equal(t, msg.Content["testProc"], true)
}
