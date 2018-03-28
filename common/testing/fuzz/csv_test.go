package fuzz

import (
	"bufio"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

func TestCSV(t *testing.T) {
	csv, err := CSV([]string{"string", "int", "string", "string"}, ';', true)
	assert.FailOnErr(t, err)

	lines := 0
	scanner := bufio.NewScanner(csv)
	for scanner.Scan() {
		lines++
	}
	// lines includes the header, which's omitted in csv.Lines
	assert.Equal(t, lines, csv.Lines+1)
}

func BenchmarkCSV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CSV([]string{"string", "int", "string", "string"}, ';', true)
	}
}
