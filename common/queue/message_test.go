package queue_test

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := queue.NewMessage(nil)
	assert.NotNil(t, msg)
	assert.NotNil(t, msg.Content)
	assert.NotNil(t, msg.Metadata[queue.MetaID])
	assert.NotNil(t, msg.Metadata[queue.MetaCreated])
	assert.Nil(t, msg.Metadata[queue.MetaBatch])
	assert.Nil(t, msg.Metadata[queue.MetaFormat])
	assert.Nil(t, msg.Metadata[queue.MetaSource])

	data := map[string]interface{}{
		"entry":     1234,
		"something": "test",
	}
	msg = queue.NewMessage(data)
	assert.NotNil(t, msg)
	assert.Equal(t, msg.Content["entry"], 1234)
	data["entry"] = 555
	assert.Equal(t, msg.Content["entry"], 555)
}

func TestMessageClear(t *testing.T) {
	msg := queue.NewMessage(nil)
	assert.Equal(t, len(msg.Content), 0)
	msg.Clear()
	assert.Equal(t, len(msg.Content), 0)

	msg = queue.NewMessage(map[string]interface{}{
		"entry":     1234,
		"something": "test",
	})
	assert.Equal(t, len(msg.Content), 2)
	msg.Clear()
	assert.Equal(t, len(msg.Content), 0)
}

func TestMessageJSON(t *testing.T) {
	msg := queue.NewMessage(nil)
	id := msg.Metadata[queue.MetaID]
	cr := msg.Metadata[queue.MetaCreated]

	buf, err := msg.ToJSON()
	assert.NoErr(t, err)

	msg.Metadata[queue.MetaID] = "123"
	msg.Metadata[queue.MetaCreated] = "123"

	err = msg.FromJSON(buf)
	assert.NoErr(t, err)
	assert.Equal(t, msg.Metadata[queue.MetaID], id)
	assert.Equal(t, msg.Metadata[queue.MetaCreated], cr)
}

func TestMSGFromJSON(t *testing.T) {
	msg := queue.NewMessage(map[string]interface{}{
		"entry":     1234,
		"something": "test",
	})
	buf, err := msg.ToJSON()
	assert.NoErr(t, err)

	msg2, err := queue.MsgFromJSON(buf)
	assert.NoErr(t, err)
	assert.FailIf(t, reflect.DeepEqual(msg, msg2))

	_, err = queue.MsgFromJSON([]byte("tmp"))
	assert.Err(t, err)
}
