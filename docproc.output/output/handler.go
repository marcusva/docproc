package output

import (
	"fmt"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
)

type Handler struct {
	ProcChain []queue.Processor
}

func NewHandler() *Handler {
	return &Handler{
		ProcChain: make([]queue.Processor, 0),
	}
}

func (hnd *Handler) AddProcessor(pp queue.Processor) {
	hnd.ProcChain = append(hnd.ProcChain, pp)
}

func (hnd *Handler) Consume(msg *queue.Message) error {
	log.Infof("Received message '%v'", msg.Metadata[queue.MetaID])
	for _, proc := range hnd.ProcChain {
		if err := proc.Process(msg); err != nil {
			return fmt.Errorf("processing failed in '%s': %v", proc.Name(), err)
		}
	}
	return nil
}
