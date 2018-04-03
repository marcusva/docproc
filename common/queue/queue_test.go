package queue

import (
	"github.com/marcusva/docproc/common/testing/assert"
	"testing"
)

func TestPackage(t *testing.T) {
	assert.NotNil(t, readqueues, "_readqueues is nil, although a package initialization was done")
	assert.NotNil(t, writequeues, "_writequeues is nil, although a package initialization was done")
}

func TestQueueTypes(t *testing.T) {
	rqtypes := ReadTypes()
	assert.ContainsS(t, rqtypes, "memory")
	assert.ContainsS(t, rqtypes, "nats")
	assert.ContainsS(t, rqtypes, "nsq")
	assert.ContainsS(t, rqtypes, "beanstalk")

	wqtypes := WriteTypes()
	assert.ContainsS(t, wqtypes, "memory")
	assert.ContainsS(t, wqtypes, "nats")
	assert.ContainsS(t, wqtypes, "nsq")
	assert.ContainsS(t, wqtypes, "beanstalk")
}

func TestCreateQueue(t *testing.T) {

	for _, rq := range ReadTypes() {
		_, err := CreateRQ(rq, nil)
		assert.Err(t, err)
	}
	for _, wq := range WriteTypes() {
		_, err := CreateWQ(wq, nil)
		assert.Err(t, err)
	}

	pmap := map[string]map[string]string{
		"memory": {
			"topic": "input",
		},
		"nats": {
			"topic": "input",
			"host":  "127.0.0.1:1234",
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
	for _, rq := range ReadTypes() {
		q, err := CreateRQ(rq, pmap[rq])
		assert.NoErr(t, err)
		assert.NotNil(t, q)
	}
	for _, wq := range WriteTypes() {
		q, err := CreateWQ(wq, pmap[wq])
		assert.NoErr(t, err)
		assert.NotNil(t, q)
	}

	_, err := CreateRQ("unknown", map[string]string{})
	assert.Err(t, err)
	assert.Equal(t, err.Error(), "Invalid readable queue type 'unknown'")
	_, err = CreateWQ("unknown", map[string]string{})
	assert.Err(t, err)
	assert.Equal(t, err.Error(), "Invalid writeable queue type 'unknown'")
}
	