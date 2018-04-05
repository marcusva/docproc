package fuzz_test

import (
	"bufio"
	"github.com/marcusva/docproc/common/testing/assert"
	"github.com/marcusva/docproc/common/testing/fuzz"
	"testing"
)

func TestSetLines(t *testing.T) {
	pairs := [][]int{
		{0, 1},
		{0, 10},
		{5, 20},
		{44, 207},
	}
	for _, p := range pairs {
		min, max := p[0], p[1]
		assert.FailOnErr(t, fuzz.SetLines(min, max))
		csv, err := fuzz.CSV([]string{"int"}, ';', true)
		assert.FailOnErr(t, err)
		assert.Equal(t, (csv.Lines >= min && csv.Lines <= max), true)
	}
	assert.FailOnErr(t, fuzz.SetLines(0, 0))
	csv, err := fuzz.CSV([]string{"int"}, ';', true)
	assert.FailOnErr(t, err)
	assert.Equal(t, (csv.Lines >= 0 && csv.Lines <= 1), true)

	assert.FailOnErr(t, fuzz.SetLines(5, 5))
	csv, err = fuzz.CSV([]string{"int"}, ';', true)
	assert.FailOnErr(t, err)
	assert.Equal(t, csv.Lines == 5, true)

	assert.NoErr(t, fuzz.SetLines(-10, 9))
	csv, err = fuzz.CSV([]string{"int"}, ';', true)
	assert.FailOnErr(t, err)
	assert.Equal(t, (csv.Lines >= 0 && csv.Lines <= 9), true)

	assert.Err(t, fuzz.SetLines(-1, -200))
	assert.Err(t, fuzz.SetLines(10, 9))
}

func TestCSV(t *testing.T) {
	csv, err := fuzz.CSV([]string{"string", "int", "string", "string"}, ';', true)
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
		fuzz.CSV([]string{"string", "int", "string", "string"}, ';', true)
	}
}
