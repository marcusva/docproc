package input

import (
	"context"
	"github.com/marcusva/docproc/common/log"
	"net/http"
	"time"
)

// WebHandler provides a request handler for HTTP requests.
type WebHandler interface {
	Transform(w http.ResponseWriter, r *http.Request)
}

// WebService is a simple web service implementation that can use one or
// more Processor implementations.
type WebService struct {
	http.Server
}

// Start starts the WebService
func (ws *WebService) Start() {
	log.Infof("starting WebService on '%s'", ws.Server.Addr)
	go func() {
		if err := ws.Server.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()
}

// Stop stops the WebService
func (ws *WebService) Stop() error {
	log.Infof("stopping WebService")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return ws.Server.Shutdown(ctx)
}

// Bind binds a WebHandler to the passed endpoint
func (ws *WebService) Bind(endpoint string, handler WebHandler) {
	log.Infof("adding endpoint handler for '%s'", endpoint)
	mux := ws.Server.Handler.(*http.ServeMux)
	mux.HandleFunc(endpoint, handler.Transform)
}

// NewWebService creates a new WebService
func NewWebService(address string) *WebService {
	return &WebService{
		Server: http.Server{
			Addr:         address,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			Handler:      http.NewServeMux(),
			ErrorLog:     log.Logger(),
		},
	}
}
