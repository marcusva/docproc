package queue

import (
	"errors"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

type testProc struct{}

func (tp *testProc) Name() string { return "testProc" }
func (tp *testProc) Process(msg *Message) error {
	if _, ok := msg.Content["works"]; !ok {
		return errors.New("failed")
	}
	msg.Content["testProc"] = true
	return nil
}

func TestNewSimpleConsumer(t *testing.T) {
	sc := NewSimpleConsumer()
	assert.Equal(t, len(sc.Processors), 0)
}

func TestSimpleConsumerAdd(t *testing.T) {
	proc := &testProc{}
	sc := NewSimpleConsumer()
	sc.Add(proc)
	assert.Equal(t, len(sc.Processors), 1)
	assert.Equal(t, sc.Processors[0], proc)
}

func TestSimpleConsumerConsume(t *testing.T) {
	proc := &testProc{}
	sc := NewSimpleConsumer()
	sc.Add(proc)
	assert.Err(t, sc.Consume(NewMessage(nil)))
	assert.Err(t, sc.Consume(NewMessage(map[string]interface{}{"test": "test"})))

	msg := NewMessage(map[string]interface{}{"works": true})
	assert.NoErr(t, sc.Consume(msg))
	assert.Equal(t, msg.Content["testProc"], true)
}
