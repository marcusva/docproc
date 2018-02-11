package input

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/marcusva/docproc/common/queue"
	"strconv"
	"strings"
	"time"
)

func init() {
	Register("RdiTransfomer", func(p map[string]string) (FileTransformer, error) {
		return &RdiTransformer{}, nil
	})
}

const (
	rdiHeader  = 'H'
	rdiControl = 'C'
	rdiSort    = 'S'
	rdiData    = 'D'
)

// RdiSection represents a simple key-value store for RDI data blocks
type RdiSection struct {
	Name    string
	Content map[string]string
}

// RdiDocument represents a set of sections belongig to a RDI document
type RdiDocument struct {
	Sections []*RdiSection
}

func docAsMap(doc *RdiDocument) map[string]interface{} {
	ret := make(map[string]interface{})
	sections := make([]interface{}, len(doc.Sections))
	for idx, sec := range doc.Sections {
		sections[idx] = secAsMap(sec)
	}
	ret["sections"] = sections
	return ret
}

func secAsMap(sec *RdiSection) map[string]interface{} {
	ret := make(map[string]interface{})
	ret["name"] = sec.Name
	ret["content"] = sec.Content
	return ret
}

// RdiTransformer represents a simple SAP RDI to queue.Message transformer.
type RdiTransformer struct {
}

// Transform transforms the passed in RDI stream into a JSON array of one or
// more RdiDocument objects.
// Only RDI data entries will be processed. Splitting the
// stream into multiple documents will be done based on the RDI header marker.
//
// An SAP RDI Data record consists of
// - A single D as data record indicator
// - 8 characters for the window name
// - 2 characters for the new window and new element start (X or space)
// - 30 characters for the element name
// - 130 characters for the symbol name
// - 1 character for the succession flag (X or space)
// - 3 characters for the value length
// - up to 255 characters for the value contents
func (tf *RdiTransformer) Transform(data []byte) ([]*queue.Message, error) {
	gzreader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gzreader.Close()

	var documents []*RdiDocument
	var curdoc *RdiDocument
	var cursection *RdiSection
	insection := false
	offset := 0

	scanner := bufio.NewScanner(gzreader)
	for scanner.Scan() {
		offset++
		line := strings.TrimSpace(scanner.Text())
		switch {
		case line[0] == rdiControl, line[0] == rdiSort:
			// Skip control data and sort records
			insection = false
			continue
		case line[0] == rdiHeader:
			// New document
			insection = false
			curdoc = &RdiDocument{Sections: []*RdiSection{}}
			documents = append(documents, curdoc)
		case line[0] == rdiData:
			// Contents
			// Section starts at D + Windowname + 2 = 11
			if !insection {
				cursection = &RdiSection{
					Name:    strings.TrimSpace(line[11:41]),
					Content: make(map[string]string),
				}
				curdoc.Sections = append(curdoc.Sections, cursection)
				insection = true
			}
			elem := strings.TrimSpace(line[41:171])
			// elem might be empty - we will skip it then
			if elem == "" {
				continue
			}
			// elem might be <SECTION>-<ELEM>, split it on demand
			splitted := strings.Split(elem, "-")
			if len(splitted) == 2 {
				elem = splitted[1]
			}
			llen, err := strconv.Atoi(line[172:175])
			if err != nil {
				return nil, err
			}
			if len(line) < (175 + llen) {
				return nil, fmt.Errorf("line %d: offset mismatch for content size %d", offset, llen)
			}
			// We will ignore the length and just copy everything (w/o the spaces)
			cursection.Content[elem] = strings.TrimSpace(line[175:])
		default:
			return nil, fmt.Errorf("line %d: unknown RDI control type '%c''", offset, line[0])
		}
	}
	ts := time.Now().Unix()
	msgs := make([]*queue.Message, len(documents))
	for idx, doc := range documents {
		msgs[idx] = queue.NewMessage(docAsMap(doc))
		msgs[idx].Metadata[queue.MetaBatch] = ts
	}
	return msgs, nil
}
