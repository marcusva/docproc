package input

import (
	"errors"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"io/ioutil"
	"os"
	"time"
)

const (
	suffixProcess = ".IN_PROCESS"
	suffixDone    = ".DONE"
	suffixFailed  = ".FAILED"
)

// FileTransformer converts an input file into one or more Message objects
type FileTransformer interface {
	Transform(data []byte) ([]*queue.Message, error)
}

// FileHandler transforms files provided by a FileProcessor and publishes the
// resulting queue.Message objects to its bound queue.Writer.
type FileHandler struct {
	*queue.Writer
	FileTransformer
}

// NewFileHandler creates a FileHandler.
func NewFileHandler(wq queue.WriteQueue, tf FileTransformer) *FileHandler {
	return &FileHandler{
		Writer:          queue.NewWriter(wq),
		FileTransformer: tf,
	}
}

func (handler *FileHandler) checkQueue() error {
	if !handler.Queue.IsOpen() {
		log.Warningf("bound queue is currently not available, trying to open it...")
		if err := handler.Queue.Open(); err != nil {
			log.Errorf("could not open queue, waiting for 10 seconds...")
			// TODO: better move this to Process()?
			time.Sleep(10 * time.Second)
			return err
		}
		return errors.New("bound queue is not open")
	}
	return nil
}

// Process reads all data from the passed file and places new messages with
// the file contents into the processing queue. On processing, the file
// will be renamed twice to indicate the processing and when it's done.
func (handler *FileHandler) Process(filename string) error {
	if err := handler.checkQueue(); err != nil {
		return err
	}
	fnameproc := filename + suffixProcess
	if err := os.Rename(filename, fnameproc); err != nil {
		log.Errorf("Could not rename file %s: %v", filename, err)
		return err
	}

	data, err := ioutil.ReadFile(fnameproc)
	if err != nil {
		log.Errorf("Could not read file %s: %v", fnameproc, err)
		return err
	}

	suffix := suffixDone
	var msgs []*queue.Message
	if msgs, err = handler.Transform(data); err != nil {
		log.Errorf("Could not transform input data: %v", err)
		suffix = suffixFailed
	} else {
		// Pass the data into the queue
		for _, msg := range msgs {
			if err := handler.Consume(msg); err != nil {
				log.Errorf("Could not publish message to queue: %v", err)
				suffix = suffixFailed
			}
		}
	}
	fnamefinal := filename + suffix
	if err := os.Rename(fnameproc, fnamefinal); err != nil {
		log.Errorf("Could not rename file %s: %v", filename, err)
		return err
	}
	return nil
}
