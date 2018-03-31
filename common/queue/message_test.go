package queue

import (
	"github.com/marcusva/docproc/common/testing/assert"
	"reflect"
	"testing"
)

func TestNewMessage(t *testing.T) {
	msg := NewMessage(nil)
	assert.NotNil(t, msg)
	assert.NotNil(t, msg.Content)
	assert.NotNil(t, msg.Metadata[MetaID])
	assert.NotNil(t, msg.Metadata[MetaCreated])
	assert.Nil(t, msg.Metadata[MetaBatch])
	assert.Nil(t, msg.Metadata[MetaFormat])
	assert.Nil(t, msg.Metadata[MetaSource])

	data := map[string]interface{}{
		"entry":     1234,
		"something": "test",
	}
	msg = NewMessage(data)
	assert.NotNil(t, msg)
	assert.Equal(t, msg.Content["entry"], 1234)
	data["entry"] = 555
	assert.Equal(t, msg.Content["entry"], 555)
}

func TestMessageClear(t *testing.T) {
	msg := NewMessage(nil)
	assert.Equal(t, len(msg.Content), 0)
	msg.Clear()
	assert.Equal(t, len(msg.Content), 0)

	msg = NewMessage(map[string]interface{}{
		"entry":     1234,
		"something": "test",
	})
	assert.Equal(t, len(msg.Content), 2)
	msg.Clear()
	assert.Equal(t, len(msg.Content), 0)
}

func TestMessageJSON(t *testing.T) {
	msg := NewMessage(nil)
	id := msg.Metadata[MetaID]
	cr := msg.Metadata[MetaCreated]

	buf, err := msg.ToJSON()
	assert.NoErr(t, err)

	msg.Metadata[MetaID] = "123"
	msg.Metadata[MetaCreated] = "123"

	err = msg.FromJSON(buf)
	assert.NoErr(t, err)
	assert.Equal(t, msg.Metadata[MetaID], id)
	assert.Equal(t, msg.Metadata[MetaCreated], cr)
}

func TestMSGFromJSON(t *testing.T) {
	msg := NewMessage(map[string]interface{}{
		"entry":     1234,
		"something": "test",
	})
	buf, err := msg.ToJSON()
	assert.NoErr(t, err)

	msg2, err := MsgFromJSON(buf)
	assert.NoErr(t, err)
	assert.FailIf(t, reflect.DeepEqual(msg, msg2))

	_, err = MsgFromJSON([]byte("tmp"))
	assert.Err(t, err)
}
