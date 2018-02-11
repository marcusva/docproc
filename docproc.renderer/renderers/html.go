package renderers

import (
	"bytes"
	"fmt"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"html/template"
)

func init() {
	Register("HTMLRenderer", NewHTMLRenderer)
}

// HTMLRenderer is a template-based renderer.
type HTMLRenderer struct {
	OutputID     string
	TemplateRoot string
	Templates    *template.Template
}

func (html *HTMLRenderer) Name() string {
	return "HTMLRenderer"
}

func (html *HTMLRenderer) Process(msg *queue.Message) error {
	buf := bytes.NewBufferString("")
	err := html.Templates.ExecuteTemplate(buf, html.TemplateRoot, msg.Content)
	if err != nil {
		log.Errorf("error on executing the templates: %v", err)
		return err
	}
	msg.Content[html.OutputID] = buf.String()
	msg.Content[ContentType] = "text/html"
	return nil
}

// NewHTMLRenderer creates a simple HTML renderer, which uses go's
// html template package
func NewHTMLRenderer(params map[string]string) (queue.Processor, error) {
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
	return &HTMLRenderer{
		OutputID:     output,
		TemplateRoot: tplroot,
		Templates:    tmpl,
	}, nil
}
