package input

import (
	"encoding/csv"
	"fmt"
	"github.com/marcusva/docproc/common/queue"
	"io"
	"time"
)

func init() {
	Register("CSVTransformer", NewCSVTransformer)
}

// CSVTransformer is a simple CSV to queue.Message transformer
type CSVTransformer struct {
	Delim rune
}

// NewCSVTransformer creates a new CSVTransformer. The CSV delimiter field
// can be configured via a 'delim' entry in the params map. If no 'delim' entry
// is provided, a comma ',' will be used as delimiter.
func NewCSVTransformer(params map[string]string) (FileTransformer, error) {
	delim, ok := params["delim"]
	if !ok {
		delim = ","
	}
	if len(delim) > 1 {
		return nil, fmt.Errorf("Invalid delimiter '%s', only one character allowed", delim)
	}
	return &CSVTransformer{
		Delim: rune(delim[0]),
	}, nil
}

// Transform creates queue.Message objects from the passed input reader. For
// each row of the CSV input data, a queue.Message will be created.
func (tf *CSVTransformer) Transform(r io.Reader) ([]*queue.Message, error) {
	reader := csv.NewReader(r)
	reader.Comma = tf.Delim

	// The first row represents the field names
	keys, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("no CSV data found")
	}
	columns := len(keys)

	msgs := []*queue.Message{}
	ts := time.Now().Unix()

	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		content := make(map[string]interface{})
		for i := 0; i < columns; i++ {
			content[keys[i]] = rec[i]
		}
		msg := queue.NewMessage(content)
		msg.Metadata[queue.MetaBatch] = ts
		msgs = append(msgs, msg)
	}
	return msgs, nil
}
