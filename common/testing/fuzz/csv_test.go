package fuzz

import (
	"bufio"
	"github.com/marcusva/docproc/common/testing/assert"
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
		assert.FailOnErr(t, SetLines(min, max))
		csv, err := CSV([]string{"int"}, ';', true)
		assert.FailOnErr(t, err)
		assert.Equal(t, (csv.Lines >= min && csv.Lines <= max), true)
	}
	assert.FailOnErr(t, SetLines(0, 0))
	csv, err := CSV([]string{"int"}, ';', true)
	assert.FailOnErr(t, err)
	assert.Equal(t, (csv.Lines >= 0 && csv.Lines <= 1), true)

	assert.FailOnErr(t, SetLines(5, 5))
	csv, err = CSV([]string{"int"}, ';', true)
	assert.FailOnErr(t, err)
	assert.Equal(t, csv.Lines == 5, true)

	assert.NoErr(t, SetLines(-10, 9))
	csv, err = CSV([]string{"int"}, ';', true)
	assert.FailOnErr(t, err)
	assert.Equal(t, (csv.Lines >= 0 && csv.Lines <= 9), true)

	assert.Err(t, SetLines(-1, -200))
	assert.Err(t, SetLines(10, 9))
}

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
