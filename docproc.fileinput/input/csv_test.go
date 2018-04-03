package input

import (
	"github.com/marcusva/docproc/common/testing/assert"
	"github.com/marcusva/docproc/common/testing/fuzz"
	"io"
	"os"
	"testing"
)

func TestCSVTransformer(t *testing.T) {
	tf, err := NewCSVTransformer(nil)
	assert.NoErr(t, err)
	assert.Equal(t, tf.(*CSVTransformer).Delim, ',')

	tf, err = NewCSVTransformer(map[string]string{"delim": ";"})
	assert.NoErr(t, err)
	assert.Equal(t, tf.(*CSVTransformer).Delim, ';')

	_, err = NewCSVTransformer(map[string]string{"delim": "###"})
	assert.Err(t, err)
}

func TestCSVTransform(t *testing.T) {
	fp, err := os.Open("test/testrecords.csv")
	assert.FailOnErr(t, err)
	defer fp.Close()

	tf, err := NewCSVTransformer(nil)
	assert.NoErr(t, err)
	tf.(*CSVTransformer).Delim = ';'

	msgs, err := tf.Transform(fp)
	assert.NoErr(t, err)
	assert.Equal(t, len(msgs), 4)

	col, ok := msgs[0].Content["CUSTNO"].(string)
	assert.Equal(t, ok, true)
	assert.Equal(t, col, "100112")

	col2, ok := msgs[0].Content["NET"].(string)
	assert.Equal(t, ok, true)
	assert.Equal(t, col2, "10394.00")
}

func TestCSVTransformFuzzed(t *testing.T) {
	tf, err := NewCSVTransformer(nil)
	assert.NoErr(t, err)
	tf.(*CSVTransformer).Delim = ';'

	for i := 0; i < 500; i++ {
		csv, err := fuzz.CSV([]string{"string", "int", "string", "string", "float", "int"}, ';', true)
		assert.FailOnErr(t, err)
		msgs, err := tf.Transform(csv)
		assert.FailOnErr(t, err)
		assert.Equal(t, len(msgs), csv.Lines)
	}
}

func BenchmarkCSVTransformLarge(b *testing.B) {
	tf, _ := NewCSVTransformer(nil)
	tf.(*CSVTransformer).Delim = ';'

	fuzz.SetLines(100000, 100000)
	csv, _ := fuzz.CSV([]string{"string", "int", "string", "string", "float", "int"}, ';', true)
	fuzz.SetLines(fuzz.MinLinesCSV, fuzz.MaxLinesCSV)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tf.Transform(csv)
		csv.Seek(0, io.SeekStart)
	}
}

func BenchmarkCSVTransform(b *testing.B) {
	tf, _ := NewCSVTransformer(nil)
	tf.(*CSVTransformer).Delim = ';'
	for i := 0; i < b.N; i++ {
		csv, _ := fuzz.CSV([]string{"string", "int", "string", "string", "float", "int"}, ';', true)
		tf.Transform(csv)
	}
}
