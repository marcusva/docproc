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
	// assert.ContainsS(t, rqtypes, "nats")
	// assert.ContainsS(t, rqtypes, "nsq")
	// assert.ContainsS(t, rqtypes, "beanstalk")
}
