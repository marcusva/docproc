package service

import (
	"fmt"
	"github.com/marcusva/docproc/common/log"
	"github.com/marcusva/docproc/common/path"
	"path/filepath"
	"time"
)

// FileProcessor processes a single file matching the pattern of the
// FileWatcher
type FileProcessor interface {
	Process(filename string) error
}

// FileProcFunc is a convenience type to avoid having to create a type to
// implement the FileProcessor interface. It can be used like this:
//
// 	NewFileWatcher(..., service.FileProcFunc(func(filename string) error {
// 		// Process the file
// 	}))
type FileProcFunc func(filename string) error

// Process implements the FileProcessor interface.
func (fn FileProcFunc) Process(filename string) error {
	return fn(filename)
}

// FileWatcher will check a specific directory periodically for files matching
// a specific pattern. If one or more files are found, it executes a processing
// function on each individual file.
type FileWatcher struct {
	// Directory is the directory to be watched for new files matching a
	// certain pattern.
	Directory string

	// Pattern is the filename pattern to use for identifying new files.
	// The FileWatcher uses the filepath.Glob() function to find matching
	// files in the specific directory.
	Pattern string

	// Interval is the time interval to use for checking the directory.
	Interval time.Duration

	// Processor receives the found file for further processing.
	Processor FileProcessor

	// stop is the stop signal channel to break out of the Watch() loop.
	stop chan (bool)
}

// NewFileWatcher creates a new FileWatcher. If the passed directory cannot
// be resolved to an absolute path, an error is returned.
func NewFileWatcher(directory, pattern string, interval time.Duration,
	processor FileProcessor) (*FileWatcher, error) {
	dir, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}
	exists, err := path.DirExists(dir)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("directory '%s' does not exist", dir)
	}

	return &FileWatcher{
		Directory: dir,
		Pattern:   pattern,
		Interval:  interval,
		Processor: processor,
		stop:      make(chan bool),
	}, nil
}

// Check checks the directory of the FileWatcher for files matching the
// specified pattern. Each file matching the pattern is processed via the
// FileWatcher's Processor.
func (w *FileWatcher) Check() {
	ppath := filepath.Join(w.Directory, w.Pattern)
	files, err := filepath.Glob(ppath)
	if err != nil {
		log.Errorf("Checking the directory failed: %v", err)
		return
	}
	for _, fpath := range files {
		log.Debugf("Processing file '%s'...", fpath)
		if err := w.Processor.Process(fpath); err != nil {
			log.Errorf("An error occurred on processing %s: %v", fpath, err)
		}
	}
}

// Watch causes the FileWatcher to check its directory periodically via the
// Check function. It will run in an endless loop until Stop is called.
func (w *FileWatcher) Watch() {
	log.Infof("Starting FileWatcher...")
	for {
		select {
		case <-time.After(w.Interval):
			log.Debugf("checking directory for new files...")
			w.Check()
		case <-w.stop:
			log.Infof("FileWatcher stopped.")
			return
		}
	}

}

// Stop stops the FileWatcher's Watch function.
func (w *FileWatcher) Stop() {
	log.Infof("Stopping FileWatcher...")
	w.stop <- true
}
