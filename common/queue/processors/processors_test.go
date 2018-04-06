package processors_test

import (
	"github.com/marcusva/docproc/common/queue/processors"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

func TestProcessorsCreateInvalid(t *testing.T) {
	_, err := processors.Create(nil)
	assert.Err(t, err)

	params := map[string]string{}
	_, err = processors.Create(params)
	assert.Err(t, err)

	params["type"] = "unknown"
	_, err = processors.Create(params)
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
	types := processors.Types()
	for _, k := range known {
		assert.ContainsS(t, types, k)
	}
}
