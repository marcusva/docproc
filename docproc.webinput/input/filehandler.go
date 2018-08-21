package input

import (
	"fmt"
	"github.com/marcusva/docproc/common/httputil"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/queue"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	fileName  = "FileHandler"
	chunkSize = 1024 * 1024
)

func init() {
	Register(fileName, NewFileHandler)
}

// FileHandler is a simple web handler for storing files sent via a
// multipart-form request in 'file'.
type FileHandler struct {
	// directory is the directory to store the retrieved file in
	directory string

	// filePrefix is the prefix to use for uploaded file names
	filePrefix string

	// fileSuffix is the suffx to use for uploaded file names
	fileSuffix string

	// maxSize is the maximum size of the file in bytes.
	maxSize int64
}

// Transform processes a http.Request's MultipartForm for a 'file' entry and
// stores the retrieved data in a randomly named file in the FileHandler's
// directory.
func (fh *FileHandler) Transform(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(fh.maxSize); err != nil {
		log.Errorf("error on parsing the request: %v", err)
		if err == multipart.ErrMessageTooLarge {
			httputil.BadRequest(w, "file exceeds the allowed file size")
		} else {
			httputil.BadRequest(w, "error on parsing request")
		}
	}

	fbuf, _, err := r.FormFile("file")
	if err != nil {
		log.Errorf("could not retrieve multipart form: %v", err)
		httputil.BadRequest(w, "could not retrieve multipart form 'file'")
		return
	}
	defer fbuf.Close()

	fname := fmt.Sprintf("%s%d%s", fh.filePrefix, time.Now().UnixNano(), fh.fileSuffix)
	fname = filepath.Join(fh.directory, fname)
	fp, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		log.Errorf("could not create file: %v", err)
		httputil.InternalServerError(w)
		return
	}
	defer fp.Close()

	if _, err := io.Copy(fp, fbuf); err != nil {
		log.Errorf("could not write file: %v", err)
		httputil.InternalServerError(w)
	}

	if err := fp.Sync(); err != nil {
		log.Errorf("error on flushing the file cache to disk: %v", err)
		httputil.InternalServerError(w)
	}
}

// NewFileHandler creates a new FileHandler
func NewFileHandler(wq queue.WriteQueue, params map[string]string) (WebHandler, error) {
	directory, ok := params["folder.out"]
	if !ok {
		return nil, fmt.Errorf("parameter 'folder.out' missing")
	}
	filePrefix, ok := params["file.prefix"]
	if !ok {
		return nil, fmt.Errorf("parameter 'file.prefix' missing")
	}
	fileSuffix, ok := params["file.suffix"]
	if !ok {
		return nil, fmt.Errorf("parameter 'file.suffix' missing")
	}
	size, ok := params["maxsize"]
	if !ok {
		return nil, fmt.Errorf("parameter 'maxsize' missing")
	}
	maxSize, err := strconv.ParseInt(size, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("parameter 'maxsize' is not an integer")
	}
	maxSize *= 1024
	return &FileHandler{
		directory:  directory,
		filePrefix: filePrefix,
		fileSuffix: fileSuffix,
		maxSize:    maxSize,
	}, nil
}
