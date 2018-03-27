package data

import (
	"bytes"
	"encoding/gob"
	"github.com/marcusva/docproc/common/log"
)

// Bytes converts the given value into a buffer of bytes.
func Bytes(val interface{}) ([]byte, error) {
	switch v := val.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	}
	log.Infof("content is not a string or byte buffer, using standard conversion")

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
