package input

import (
	"github.com/marcusva/docproc/common/testing/assert"
	"io/ioutil"
	"testing"
)

func TestCSVTransformer(t *testing.T) {
	tf, err := NewCSVTransformer(nil)
	assert.NoErr(t, err)
	assert.Equal(t, tf.(*CSVTransformer).Delim, ',')

	tf, err = NewCSVTransformer(map[string]string{"delim": ";"})
	assert.NoErr(t, err)
	assert.Equal(t, tf.(*CSVTransformer).Delim, ';')

	tf, err = NewCSVTransformer(map[string]string{"delim": "###"})
	assert.Err(t, err)
}

func TestCsvTransform(t *testing.T) {
	buf, err := ioutil.ReadFile("test/testrecords.csv")
	assert.FailOnErr(t, err)

	tf, err := NewCSVTransformer(nil)
	assert.NoErr(t, err)
	tf.(*CSVTransformer).Delim = ';'

	msgs, err := tf.Transform(buf)
	assert.NoErr(t, err)
	assert.Equal(t, len(msgs), 4)

	col, ok := msgs[0].Content["CUSTNO"].(string)
	assert.Equal(t, ok, true)
	assert.Equal(t, col, "100112")

	col2, ok := msgs[0].Content["NET"].(string)
	assert.Equal(t, ok, true)
	assert.Equal(t, col2, "10394.00")

}
