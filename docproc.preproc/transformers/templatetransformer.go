package transformers

import (
	"bytes"
	"fmt"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"text/template"
)

func init() {
	Register("TemplateTransformer", NewTemplateTransformer)
}

// TemplateTransformer provides a simple mechanism to add additional data to
// a queue.Message via Go's text/template transformation support.
// The transformation result will be stored in the message's Content section and
// can be reached via the transformer's Identifier.
type TemplateTransformer struct {
	OutputID     string
	TemplateRoot string
	Templates    *template.Template
}

// Name returns the name of the TemplateTransformer to be used in configuration
// files.
func (tf *TemplateTransformer) Name() string {
	return "TemplateTransformer"
}

// Process takes a queue.Message and executes the configured templates using
// the message's Content.
func (tf *TemplateTransformer) Process(msg *queue.Message) error {
	buf := bytes.NewBufferString("")
	err := tf.Templates.ExecuteTemplate(buf, tf.TemplateRoot, msg.Content)
	if err != nil {
		log.Errorf("Executing the template '%s' failed for content %v", tf.TemplateRoot, msg.Content)
		return err
	}
	msg.Content[tf.OutputID] = buf.String()
	return nil
}

// NewTemplateTransformer creates a new TemplateTransformer
// TODO: describe param map
func NewTemplateTransformer(params map[string]string) (queue.Processor, error) {
	tplfiles, ok := params["templates"]
	if !ok {
		return nil, fmt.Errorf("parameter 'templates' missing")
	}
	output, ok := params["output"]
	if !ok {
		return nil, fmt.Errorf("parameter 'output' missing")
	}
	tplroot, ok := params["templateroot"]
	if !ok {
		return nil, fmt.Errorf("parameter 'templateroot' missing")
	}
	tmpl, err := template.ParseGlob(tplfiles)
	if err != nil {
		return nil, err
	}

	return &TemplateTransformer{
		OutputID:     output,
		TemplateRoot: tplroot,
		Templates:    tmpl,
	}, nil
}
