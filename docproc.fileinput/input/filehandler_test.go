package input

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

type nopTransformer int

func (nop nopTransformer) Transform(data []byte) ([]*queue.Message, error) {
	return []*queue.Message{
		queue.NewMessage(map[string]interface{}{"buf": data}),
	}, nil
}

var (
	buf = []byte("some sample data")
)

func createWQ() queue.WriteQueue {
	wq, _ := queue.CreateWQ("memory", map[string]string{"topic": "test"})
	return wq
}

func TestTransform(t *testing.T) {
	fh := NewFileHandler(createWQ(), nopTransformer(0))

	messages, err := fh.Transform(buf)
	assert.NoErr(t, err)
	assert.FailIf(t, len(messages) != 1)
	assert.Equal(t, string(messages[0].Content["buf"].([]byte)), string(buf))
}

func TestProcess(t *testing.T) {
	// FIXME: implement this
	t.Skip("not yet implemented")
}
