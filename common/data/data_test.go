package data

import (
	"bytes"
	"encoding/gob"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

func TestBytes(t *testing.T) {
	str := "this is a test"
	buf, err := Bytes(str)
	assert.NoErr(t, err)
	assert.Equal(t, buf, []byte("this is a test"))

	bbuf := []byte{0x10, 0x10, 0x00, 0x99}
	buf, err = Bytes(bbuf)
	assert.NoErr(t, err)
	assert.Equal(t, buf, bbuf)

	data := map[string]int{
		"key":  1234,
		"test": 9876,
	}
	buf, err = Bytes(data)
	assert.NoErr(t, err)

	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	var decdata map[string]int
	err = dec.Decode(&decdata)
	assert.FailOnErr(t, err)
	assert.Equal(t, data, decdata)
}
