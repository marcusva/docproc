package input

import (
	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	wq, _ := queue.CreateWQ("memory", map[string]string{"topic": "q"})
	_, err := Create(nil, nil)
	assert.Err(t, err)

	_, err = Create(wq, map[string]string{
		"type": "InvalidHandler",
	})
	assert.Err(t, err)

	_, err = Create(wq, map[string]string{
		"type":    "RawHandler",
		"maxsize": "1",
	})
	assert.NoErr(t, err)

	_, err = Create(wq, map[string]string{
		"type":        "FileHandler",
		"maxsize":     "1",
		"file.prefix": "out-",
		"file.suffix": ".csv",
		"folder.out":  "/app/data",
	})
	assert.NoErr(t, err)
}

func TestWebHandlers(t *testing.T) {
	handlers := WebHandlers()
	assert.NotNil(t, handlers)
	assert.ContainsS(t, handlers, "RawHandler")
	assert.ContainsS(t, handlers, "FileHandler")
}
