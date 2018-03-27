package fuzz

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

var (
	maxLinesCSV  = 2500
	maxLenString = 10
)

// FuzzedCSV is an io.Rader that contains randomly generated CSV data.
type FuzzedCSV struct {
	io.Reader
	// Columns contains the column count of the CSV.
	Columns int

	// Lines holds the line count of the generated CSV, excluding the optional
	// header of the CSV.
	Lines int
}

// SetMaxLines sets the maximum number of CSV lines to generate. If maxlines
// is smaller than 1, 1 will be set.
func SetMaxLines(maxlines int) {
	if maxlines <= 0 {
		maxlines = 1
	}
	maxLinesCSV = maxlines
}

// SetMaxLenString sets the maximum length of a single string columns.
// If maxlen is smaller than 1, 1 will be set.
func SetMaxLenString(maxlen int) {
	if maxlen <= 0 {
		maxlen = 1
	}
	maxLenString = maxlen
}

func createRecord(types []string) []string {
	record := make([]string, len(types))
	for idx, t := range types {
		switch t {
		case "int":
			record[idx] = strconv.FormatInt(rand.Int63(), 10)
		case "float":
			record[idx] = strconv.FormatFloat(rand.Float64(), 'e', rand.Intn(24), 64)
		case "bool":
			record[idx] = strconv.FormatBool(rand.Int63n(2) > 0)
		case "string":
			len := rand.Intn(maxLenString) + 1
			buf := make([]byte, len)
			for i := 0; i < len; i++ {
				buf[i] = byte(rand.Int63n(94) + 32) // any ASCII character in the range 0x20 - 0x7E
			}
			record[idx] = string(buf)
		default:
			record[idx] = ""
		}
	}
	return record
}

func validColumnTypes(types []string) error {
	for _, t := range types {
		switch t {
		case "int":
		case "float":
		case "bool":
		case "string":
			continue
		default:
			return fmt.Errorf("invalid column type '%s'", t)
		}
	}
	return nil
}

// CSV returns an in-memory io.Reader containing random CSV data.
func CSV(types []string, delim rune, headers bool) (*FuzzedCSV, error) {
	if err := validColumnTypes(types); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = delim
	switch runtime.GOOS {
	case "windows":
		writer.UseCRLF = true
	default:
		writer.UseCRLF = false
	}

	rand.Seed(time.Now().UnixNano())

	if headers {
		headline := make([]string, len(types))
		for idx, t := range types {
			headline[idx] = fmt.Sprintf("Header [%s]", t)
		}
		if err := writer.Write(headline); err != nil {
			return nil, err
		}
	}
	maxlines := rand.Intn(maxLinesCSV)
	for i := 0; i < maxlines; i++ {
		if err := writer.Write(createRecord(types)); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return &FuzzedCSV{
		Reader:  &buf,
		Lines:   maxlines,
		Columns: len(types),
	}, nil
}
