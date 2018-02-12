package input

import (
	"errors"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"io/ioutil"
	"os"
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
			log.Errorf("could not open queue: %v", err)
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
	if err := renameFile(filename, fnameproc); err != nil {
		return err
	}

	data, err := ioutil.ReadFile(fnameproc)
	if err != nil {
		log.Errorf("could not read file %s: %v", fnameproc, err)
		renameFile(fnameproc, filename+suffixFailed)
		return err
	}

	var msgs []*queue.Message
	if msgs, err = handler.Transform(data); err != nil {
		log.Errorf("could not transform input data: %v", err)
		renameFile(fnameproc, filename+suffixFailed)
		return err
	}

	// Pass the data into the queue
	for _, msg := range msgs {
		if err := handler.Consume(msg); err != nil {
			log.Errorf("could not publish message to queue: %v", err)
			renameFile(fnameproc, filename+suffixFailed)
			return err
		}
	}
	return renameFile(fnameproc, suffixDone)
}

func renameFile(curname, target string) error {
	if err := os.Rename(curname, target); err != nil {
		log.Errorf("could not rename file %s: %v", curname, err)
		return err
	}
	return nil
}
