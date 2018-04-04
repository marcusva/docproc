package processors

import (
	"fmt"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"os"
	"time"
)

const (
	perfName = "PerformanceChecker"
)

func init() {
	Register(perfName, NewPerformanceChecker)
}

// PerformanceChecker measures the processing time of queue.Messages.
type PerformanceChecker struct {
	filename string
}

// Name gets the name of the PerformanceChecker.
func (perf *PerformanceChecker) Name() string {
	return perfName
}

// Process processes a queue.Message and writes the processing time since its
// initial creation and current consumption into the logs.
func (perf *PerformanceChecker) Process(msg *queue.Message) error {
	curtime := time.Now()

	id, ok := msg.Metadata[queue.MetaID]
	if !ok {
		log.Info("message metadata does not contain an ID, using '<unknown>'")
		id = "<unknown>"
	}
	ts, ok := msg.Metadata[queue.MetaCreated]
	if !ok {
		log.Info("message metadata does not contain a timestamp")
	}

	var durCreation time.Duration
	switch v := ts.(type) {
	case time.Time:
		durCreation = curtime.Sub(v)
	case string:
		msgtime, err := time.Parse(time.RFC3339Nano, v)
		if err != nil {
			log.Errorf("could not retrieve the message time for '%s'", id)
		}
		durCreation = curtime.Sub(msgtime)
	default:
		// ts might be nil or something different
		log.Errorf("could not retrieve the message time for '%s'", id)
		durCreation = time.Duration(-1)
	}
	durCurrent := curtime.Sub(msg.Timestamp)

	log.Infof("proc time of '%s': cur: %v, creation: %v", id, durCurrent, durCreation)
	if perf.filename != "" {
		fp, err := os.OpenFile(perf.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(0644))
		if err != nil {
			return err
		}
		defer fp.Close()

		line := fmt.Sprintf("%s;%v;%v;\n", id, durCurrent, durCreation)
		if _, err := fp.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}

// NewPerformanceChecker creates a PerformanceChecker.
func NewPerformanceChecker(params map[string]string) (queue.Processor, error) {
	filename, ok := params["write.to"]
	if !ok {
		log.Info("'write.to not set, using standard logger'")
		filename = ""
	}
	return &PerformanceChecker{
		filename: filename,
	}, nil
}
