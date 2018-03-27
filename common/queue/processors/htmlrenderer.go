package processors

import (
	"bytes"
	"fmt"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"html/template"
)

const (
	htmlName = "HTMLRenderer"
)

var (
	// ContentType represents the rendered result's content type.
	ContentType = "mime-type"
)

func init() {
	Register(htmlName, NewHTMLRenderer)
}

// HTMLRenderer is a template-based renderer for HTML content.
type HTMLRenderer struct {
	storeIn      string
	templateRoot string
	templates    *template.Template
}

// Name returns the name to be used in configuration files.
func (html *HTMLRenderer) Name() string {
	return htmlName
}

// Process processes the passed in message using the configured templates.
// On success, the result will be stored as a key-value pair using the
// HTMLRenderer's set storeIn value. Additionally, a "mime-type" : "text/html"
// key-value pair will be set on the message.
func (html *HTMLRenderer) Process(msg *queue.Message) error {
	buf := bytes.NewBufferString("")
	err := html.templates.ExecuteTemplate(buf, html.templateRoot, msg.Content)
	if err != nil {
		log.Errorf("error on executing the templates: %v", err)
		return err
	}
	msg.Content[html.storeIn] = buf.String()
	msg.Content[ContentType] = "text/html"
	return nil
}

// NewHTMLRenderer creates a simple HTML renderer, which uses go's
// html template package.
func NewHTMLRenderer(params map[string]string) (queue.Processor, error) {
	tplfiles, ok := params["templates"]
	if !ok {
		return nil, fmt.Errorf("parameter 'templates' missing")
	}
	output, ok := params["store.in"]
	if !ok {
		return nil, fmt.Errorf("parameter 'store.in' missing")
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
		storeIn:      output,
		templateRoot: tplroot,
		templates:    tmpl,
	}, nil
}
