// Package log is a simple wrapper around the log package, which adds RFC5424
// severity thresholds.
package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

const (
	// LevelEmergency represents the most severe level EMERGENCY
	LevelEmergency Level = iota
	// LevelAlert represents the ALERT level
	LevelAlert
	// LevelCritical represents the CRITICAL level
	LevelCritical
	// LevelError represents the ERROR level
	LevelError
	// LevelWarning represents the WARNING level
	LevelWarning
	// LevelNotice represents the NOTICE level
	LevelNotice
	// LevelInfo represents the INFO level
	LevelInfo
	// LevelDebug represents the DEBUG level
	LevelDebug
)

// Level represents The log level threshold value type
type Level int8

var (
	logger     *log.Logger
	fpLogfile  *os.File
	showCaller bool
	threshold  Level
	mux        = sync.Mutex{}
)

func init() {
	// package initialization
	if logger == nil {
		Init(ioutil.Discard, LevelError, false)
	}
}

// Logger gets the logger being used.
func Logger() *log.Logger {
	mux.Lock()
	defer mux.Unlock()
	return logger
}

// CurrentLevel returns the current log level for the logger.
func CurrentLevel() Level {
	mux.Lock()
	defer mux.Unlock()
	return threshold
}

// GetLogLevel gets a matching LogLevel value from the passed string. Returns
// an error,  if the string does not match a valid log level.
func GetLogLevel(level string) (Level, error) {
	switch level {
	case "0", "Emergency":
		return LevelEmergency, nil
	case "1", "Alert":
		return LevelAlert, nil
	case "2", "Critical":
		return LevelCritical, nil
	case "3", "Error":
		return LevelError, nil
	case "4", "Warning":
		return LevelWarning, nil
	case "5", "Notice":
		return LevelNotice, nil
	case "6", "Informational", "Info":
		return LevelInfo, nil
	case "7", "Debug":
		return LevelDebug, nil
	default:
		return -1, fmt.Errorf("unknown log level '%s'", level)
	}
}

// InitFile initializes the logging functionality using the passed file.
// This will close the currently open logfile, if the logger has been
// initialized with InitFile before.
func InitFile(logfile string, level Level, caller bool) error {
	logfp, err := os.OpenFile(logfile, os.O_RDWR, os.FileMode(0600))
	if err != nil {
		return err
	}
	Init(logfp, level, caller)
	mux.Lock()
	defer mux.Unlock()
	fpLogfile = logfp
	return nil
}

// Init (re)initializes the logging functionality. out will be the stream to
// write to, level is the threshold to use as maximum value to log.
// This will close the currently open logfile, if the logger has been
// initialized with InitFile before.
func Init(out io.Writer, level Level, caller bool) {
	mux.Lock()
	defer mux.Unlock()
	if fpLogfile != nil {
		fpLogfile.Close()
		fpLogfile = nil
	}
	threshold = level
	showCaller = caller
	logger = log.New(out, "", log.LstdFlags)
}

// NoisyInit initializes the logging functionality with a debug level on the
// standard output.
// This will close the currently open logfile, if the logger has been
// initialized with InitFile before.
func NoisyInit() {
	Init(os.Stdout, LevelDebug, true)
}

func _printval(prefix string, v ...interface{}) {
	if showCaller {
		// Arg to Caller(): 0 = this func, 1 = previous (_log.XXX), 2: caller
		if _, file, line, ok := runtime.Caller(2); ok {
			fdata := fmt.Sprintf("[%s:%d]", filepath.Base(file), line)
			logger.Printf("%-9s %s %v\n", prefix, fdata, v)
			return
		}
	}
	logger.Printf("%-9s %v\n", prefix, v)
}

func _printstr(prefix string, output string) {
	if showCaller {
		// Arg to Caller(): 0 = this func, 1 = previous (_log.XXX), 2: caller
		if _, file, line, ok := runtime.Caller(2); ok {
			fdata := fmt.Sprintf("[%s:%d]", filepath.Base(file), line)
			logger.Printf("%-9s %s %s\n", prefix, fdata, output)
			return
		}
	}
	logger.Printf("%-9s %s\n", prefix, output)
}

// Debug writes a debug message to the log.
func Debug(args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelDebug {
		_printval("DEBUG", args)
	}
}

// Debugf writes a debug message to the log.
func Debugf(format string, args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelDebug {
		_printstr("DEBUG", fmt.Sprintf(format, args...))
	}
}

// Info writes an informational message to the log.
func Info(args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelInfo {
		_printval("INFO", args)
	}
}

// Infof writes an informational message to the log.
func Infof(format string, args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelInfo {
		_printstr("INFO", fmt.Sprintf(format, args...))
	}
}

// Notice writes a notice message to the log.
func Notice(args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelNotice {
		_printval("NOTICE", args)
	}
}

// Noticef writes a notice message to the log.
func Noticef(format string, args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelNotice {
		_printstr("NOTICE", fmt.Sprintf(format, args...))
	}
}

// Warning writes a warning message to the log.
func Warning(args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelWarning {
		_printval("WARNING", args)
	}
}

// Warningf writes a warning message to the log.
func Warningf(format string, args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelWarning {
		_printstr("WARNING", fmt.Sprintf(format, args...))
	}
}

// Error writes an error message to the log.
func Error(args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelError {
		_printval("ERROR", args)
	}
}

// Errorf writes an error message to the log.
func Errorf(format string, args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelError {
		_printstr("ERROR", fmt.Sprintf(format, args...))
	}
}

// Critical writes a critical message to the log.
func Critical(args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelCritical {
		_printval("CRITICAL", args)
	}
}

// Criticalf writes a critical message to the log.
func Criticalf(format string, args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelCritical {
		_printstr("CRITICAL", fmt.Sprintf(format, args...))
	}
}

// Alert writes an alert message to the log.
func Alert(args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelAlert {
		_printval("ALERT", args)
	}
}

// Alertf rites an alert message to the log.
func Alertf(format string, args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelAlert {
		_printstr("ALERT", fmt.Sprintf(format, args...))
	}
}

// Emergency writes an emergency message to the log.
func Emergency(args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelEmergency {
		_printval("EMERGENCY", args)
	}
}

// Emergencyf writes an emergency message to the log.
func Emergencyf(format string, args ...interface{}) {
	mux.Lock()
	defer mux.Unlock()
	if threshold >= LevelEmergency {
		_printstr("EMERGENCY", fmt.Sprintf(format, args...))
	}
}
