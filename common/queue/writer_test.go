package queue_test

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/queue/processors"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

const (
	message = `{
	"metadata": {
		"batch": 1517607828,
		"created": "2018-02-02T22:43:48.0220047+01:00"
	},
	"content": {
		"CITY": "New York",
		"CUSTNO": "100112",
		"DATE": "2017-04-07",
		"FIRSTNAME": "John",
		"GROSS": "12386.86",
		"LASTNAME": "Doe",
		"NET": "10394.00",
		"STREET": "Example Lane 384",
		"ZIP": "10006"
	}
}
`
)

func TestWriter(t *testing.T) {
	w := queue.NewWriter(nil, nil)
	assert.NotNil(t, w)

	assert.NoErr(t, w.Open())
	assert.NoErr(t, w.Close())

	proc, err := processors.NewValueEnricher(map[string]string{"rules": "processors/test/testrules.json"})
	assert.FailOnErr(t, err)

	w.Add(proc)

	msg, err := queue.MsgFromJSON([]byte(message))
	assert.FailOnErr(t, err)

	assert.NoErr(t, w.Consume(msg))
	assert.Equal(t, msg.Content["DOCTYPE"], "INVOICE")
}

type nopProc struct{}

func (p *nopProc) Name() string { return "nopProc" }
func (p *nopProc) Process(msg *queue.Message) error {
	return nil
}

func BenchmarkWriterEmpty(b *testing.B) {
	w := queue.NewWriter(nil, nil)
	w.Add(&nopProc{})

	msg, _ := queue.MsgFromJSON([]byte(message))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Consume(msg)
	}
}

func BenchmarkWriterQueues(b *testing.B) {
	wq, _ := queue.CreateWQ("memory", map[string]string{"topic": "out"})
	errq, _ := queue.CreateWQ("memory", map[string]string{"topic": "error"})

	w := queue.NewWriter(wq, errq)
	w.Add(&nopProc{})
	msg, _ := queue.MsgFromJSON([]byte(message))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Consume(msg)
	}
}
