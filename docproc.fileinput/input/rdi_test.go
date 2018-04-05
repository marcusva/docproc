package input_test

import (
	// "encoding/json"
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"github.com/marcusva/docproc/docproc.fileinput/input"
	"os"
	"testing"
)

func TestRDITransformer(t *testing.T) {
	t.Skip()

	fp, err := os.Open("test/testrdi.gz")
	assert.FailOnErr(t, err)
	defer fp.Close()

	tf := &input.RDITransformer{}

	messages, err := tf.Transform(fp)
	assert.FailOnErr(t, err)

	assert.Equal(t, len(messages), 7)

	ts := messages[0].Metadata[queue.MetaBatch]
	for _, m := range messages {
		assert.Equal(t, ts, m.Metadata[queue.MetaBatch])
		assert.Equal(t, len(m.Content["sections"].([]interface{})) > 0, true)
		// if true {
		// 	data, _ := json.MarshalIndent(m, "", "  ")
		// 	t.Errorf("%s", data)
		// }
	}
}
