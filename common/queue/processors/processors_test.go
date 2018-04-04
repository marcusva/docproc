package processors

import (
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

func TestProcessorsCreateInvalid(t *testing.T) {
	_, err := Create(nil)
	assert.Err(t, err)

	params := map[string]string{}
	_, err = Create(params)
	assert.Err(t, err)

	params["type"] = "unknown"
	_, err = Create(params)
	assert.Err(t, err)
}

func TestProcessorsTypes(t *testing.T) {
	known := []string{
		"CommandProc",
		"ContentValidator",
		"FileWriter",
		"HTMLRenderer",
		"HTTPSender",
		"PerformanceChecker",
		"TemplateTransformer",
		"ValueEnricher",
	}
	types := Types()
	for _, k := range known {
		assert.ContainsS(t, types, k)
	}
}
