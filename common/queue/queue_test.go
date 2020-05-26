package queue_test

import (
	"testing"

	"github.com/marcusva/docproc/common/queue"
	"github.com/marcusva/docproc/common/testing/assert"
)

func TestPackage(t *testing.T) {
	assert.FailIf(t, len(queue.ReadTypes()) == 0,
		"ReadTypes() is nil, although a package initialization was done")
	assert.FailIf(t, len(queue.WriteTypes()) == 0,
		"WriteTypes() is nil, although a package initialization was done")
}

func TestQueueTypes(t *testing.T) {
	rqtypes := queue.ReadTypes()
	assert.ContainsS(t, rqtypes, "memory")
	assert.ContainsS(t, rqtypes, "nsq")
	assert.ContainsS(t, rqtypes, "beanstalk")

	wqtypes := queue.WriteTypes()
	assert.ContainsS(t, wqtypes, "memory")
	assert.ContainsS(t, wqtypes, "nsq")
	assert.ContainsS(t, wqtypes, "beanstalk")
}

func TestCreateQueue(t *testing.T) {

	for _, rq := range queue.ReadTypes() {
		_, err := queue.CreateRQ(rq, nil)
		assert.Err(t, err)
	}
	for _, wq := range queue.WriteTypes() {
		_, err := queue.CreateWQ(wq, nil)
		assert.Err(t, err)
	}

	pmap := map[string]map[string]string{
		"memory": {
			"topic": "input",
		},
		"nsq": {
			"topic": "input",
			"host":  "127.0.0.1:1234",
		},
		"beanstalk": {
			"topic": "input",
			"host":  "127.0.0.1:1234",
		},
	}
	for _, rq := range queue.ReadTypes() {
		q, err := queue.CreateRQ(rq, pmap[rq])
		assert.NoErr(t, err)
		assert.NotNil(t, q)
	}
	for _, wq := range queue.WriteTypes() {
		q, err := queue.CreateWQ(wq, pmap[wq])
		assert.NoErr(t, err)
		assert.NotNil(t, q)
	}

	_, err := queue.CreateRQ("unknown", map[string]string{})
	assert.Err(t, err)
	assert.Equal(t, err.Error(), "Invalid readable queue type 'unknown'")
	_, err = queue.CreateWQ("unknown", map[string]string{})
	assert.Err(t, err)
	assert.Equal(t, err.Error(), "Invalid writeable queue type 'unknown'")
}
