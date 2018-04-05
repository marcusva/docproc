package data_test

import (
	"bytes"
	"encoding/gob"
	"github.com/marcusva/docproc/common/data"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

func TestBytes(t *testing.T) {
	str := "this is a test"
	buf, err := data.Bytes(str)
	assert.NoErr(t, err)
	assert.Equal(t, buf, []byte("this is a test"))

	bbuf := []byte{0x10, 0x10, 0x00, 0x99}
	buf, err = data.Bytes(bbuf)
	assert.NoErr(t, err)
	assert.Equal(t, buf, bbuf)

	mapped := map[string]int{
		"key":  1234,
		"test": 9876,
	}
	buf, err = data.Bytes(mapped)
	assert.NoErr(t, err)

	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	var decdata map[string]int
	err = dec.Decode(&decdata)
	assert.FailOnErr(t, err)
	assert.Equal(t, mapped, decdata)
}
