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
	Register("RDITransfomer", func(p map[string]string) (FileTransformer, error) {
		return &RDITransformer{}, nil
	})
}

const (
	rdiHeader  = 'H'
	rdiControl = 'C'
	rdiSort    = 'S'
	rdiData    = 'D'
)

// RDISection represents a simple key-value store for RDI data blocks
type RDISection struct {
	Name    string
	Content map[string]string
}

// RDIDocument represents a set of sections belongig to a RDI document
type RDIDocument struct {
	Sections []*RDISection
}

// docAsMap converts a RDIDocument to a simple map representation.
func docAsMap(doc *RDIDocument) map[string]interface{} {
	ret := make(map[string]interface{})
	sections := make([]interface{}, len(doc.Sections))
	for idx, sec := range doc.Sections {
		sections[idx] = secAsMap(sec)
	}
	ret["sections"] = sections
	return ret
}

// secAsMap converts a RDISection to a simple map representation.
func secAsMap(sec *RDISection) map[string]interface{} {
	ret := make(map[string]interface{})
	ret["name"] = sec.Name
	ret["content"] = sec.Content
	return ret
}

// RDITransformer represents a simple SAP RDI to queue.Message transformer.
type RDITransformer struct {
}

// Transform transforms the passed in RDI stream into a set of queue.Message
// objects. Splitting the stream into multiple documents - and thus multiple
// messages - will be done based on the RDI header marker.
//
// Only RDI data entries will be processed and end up in the content section
// of the generated Message objects. RID control or sort entries will be
// ignored.
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
//
func (tf *RDITransformer) Transform(data []byte) ([]*queue.Message, error) {
	gzreader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gzreader.Close()

	var documents []*RDIDocument
	var curdoc *RDIDocument
	var cursection *RDISection
	insection := false
	offset := 0

	scanner := bufio.NewScanner(gzreader)
	for scanner.Scan() {
		offset++
		line := strings.TrimSpace(scanner.Text())
		switch line[0] {
		case rdiControl, rdiSort:
			// Skip control data and sort records
			insection = false
			continue
		case rdiHeader:
			// New document
			insection = false
			curdoc = &RDIDocument{Sections: []*RDISection{}}
			documents = append(documents, curdoc)
		case rdiData:
			// Contents
			// Section starts at D + Windowname + 2 = 11
			if !insection {
				cursection = &RDISection{
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
