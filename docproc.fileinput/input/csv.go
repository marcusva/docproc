package input

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/marcusva/docproc/common/queue"
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

// Transform creates queue.Message objects from the passed input data. For each
// row of the CSV input data, a queue.Message will be created.
func (tf *CSVTransformer) Transform(data []byte) ([]*queue.Message, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = tf.Delim
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("no CSV data found")
	}

	// The first row represents the field names
	keys := records[0]
	records = records[1:]
	columns := len(keys)

	msgs := make([]*queue.Message, 1)
	ts := time.Now().Unix()
	for _, rec := range records {
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
